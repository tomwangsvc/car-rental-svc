package pubsub

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_http "github.com/tomwangsvc/lib-svc/http"
	lib_json "github.com/tomwangsvc/lib-svc/json"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_token_gcp "github.com/tomwangsvc/lib-svc/token/gcp"
	lib_token_svc "github.com/tomwangsvc/lib-svc/token/svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Constants for pubsub
const (
	EventListenerDefaultTimeout = time.Second * 30

	pushMessageExpiry = time.Hour * 24 * 7

	// Task behavior from Pubsub can be achieved by using Synchronous mode true, a flow control is no longer needed
	// 1. https://github.com/GoogleCloudPlatform/google-cloud-go/issues/1193
	// 2. https://github.com/GoogleCloudPlatform/google-cloud-go/issues/1088
	//      -> https://github.com/GoogleCloudPlatform/google-cloud-go/issues/919#issuecomment-372412775

	// Example of where this is useful
	// Consider orders involving items we have not seen before:
	//   These orders require manual configuration before they can be ack'ed
	//   If we nack these orders they are immediately re-delivered in milliseconds which is not helpful
	//   Instead we want a delay before attempting to reprocess them i.e. a 'delay queue'

	SubscriptionDefaultLongRetryDelay                        = time.Second * 60
	SubscriptionDefaultShortRetryDelay                       = time.Second * 10
	subscriptionDefaultReceiveSettingsMaxOutstandingMessages = -1
	subscriptionDefaultReceiveSettingsSynchronous            = true // Required to be true for task behavior e.g. a delay queue
)

var (
	jsonRedactionsByPathForMessage = map[string]lib_json.Redaction{
		"message.data": redactionForMessage,
	}
)

func redactionForMessage(data string) string {
	d, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		lib_log.Warn(context.Background(), fmt.Sprintf("\"value\":\"NOT BASE64 DECODABLE\",\"error\":\"%+v\"", err), lib_log.FmtString("data", data))
		return fmt.Sprintf("{\"data_decoded\":\"NOT BASE64 DECODABLE\",\"bytes_length\":%d}", len(data))
	}
	s := string(d)

	if len(s) > 1024 {
		return fmt.Sprintf("{\"data_decoded\":\"%s\",\"bytes_length\":%d}", s[:1024], len(s))
	}

	return s
}

func DelayQueueSubscriptionConfig(topic *pubsub.Topic, retryDelay time.Duration) pubsub.SubscriptionConfig {
	return pubsub.SubscriptionConfig{Topic: topic, AckDeadline: retryDelay}
}

func DelayQueueSubscriptionReceiveSettings(receiveSettings pubsub.ReceiveSettings, subscriptionConfig pubsub.SubscriptionConfig) pubsub.ReceiveSettings {
	receiveSettings.MaxExtension = subscriptionConfig.AckDeadline // Used to enforce the subscription AckDeadline
	receiveSettings.MaxOutstandingMessages = subscriptionDefaultReceiveSettingsMaxOutstandingMessages
	receiveSettings.Synchronous = subscriptionDefaultReceiveSettingsSynchronous
	return receiveSettings
}

type Client interface {
	Close() error

	Authorize(next http.Handler) http.Handler

	CreateTopic(ctx context.Context, topicId string) (*pubsub.Topic, error)
	CreateTopics(ctx context.Context, topicIds []string) (topicByTopicId map[string]*pubsub.Topic, err error)
	Topic(ctx context.Context, topicId string) (*pubsub.Topic, error)
	Topics(ctx context.Context, topicIds []string) (topicByTopicId map[string]*pubsub.Topic, err error)

	PublishMessage(ctx context.Context, topic *pubsub.Topic, data []byte) (messageId string, err error)
	PublishMessages(ctx context.Context, topic *pubsub.Topic, datas [][]byte) (messageIds []string, err error)
	Unmarshal(ctx context.Context, topic, subscription string, msg *pubsub.Message, dst interface{}) (newCtx context.Context, createdInIntegrationTest, createdInIntegrationTestAutoAckDisable bool, err error)
	UnmarshalWithNoAttributes(ctx context.Context, topic, subscription string, msg *pubsub.Message, dst interface{}) error
}

type Config struct {
	Env lib_env.Env
}

func NewClient(ctx context.Context, config Config, tokenSvcClient lib_token_svc.Client, tokenGcpClient lib_token_gcp.Client) (Client, error) {
	lib_log.Info(ctx, "Initializing", lib_log.FmtAny("config", config))

	pubsubClient, err := pubsub.NewClient(ctx, config.Env.GcpProjectId)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed initializing pubsub client")
	}

	lib_log.Info(ctx, "Initialized")
	return client{
		config:            config,
		cloudPubsubClient: cloudPubsubClient{pubsubClient},
		tokenSvcClient:    tokenSvcClient,
		tokenGcpClient:    tokenGcpClient,
	}, nil
}

type client struct {
	config            Config
	cloudPubsubClient cloudPubsubClientInterface
	tokenSvcClient    lib_token_svc.Client
	tokenGcpClient    lib_token_gcp.Client
}

func (c client) Close() error {
	if err := c.cloudPubsubClient.Close(); err != nil {
		return lib_errors.Wrap(err, "Failed closing pubsub client")
	}
	return nil
}

//revive:disable:cyclomatic
func (c client) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		lib_log.Info(ctx, "Authorizing")

		ctx = lib_context.WithPubsubPush(ctx, true)

		token, err := lib_http.TokenFromAuthorizationHeader(r.Header)
		if err != nil {
			lib_log.Warn(ctx, "Unauthorized due to failure reading token from authorization header for pubsub push, will log request body and then respond with No Content so that message is ACKed", lib_log.FmtError(err))
			if err := logRequestBody(r); err != nil {
				lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed logging request body"))
				return
			}
			lib_http.RenderNoContent(ctx, w)
			return
		}

		if lib_context.IntegrationTest(ctx) {
			if _, err := c.tokenSvcClient.VerifyTokenAndExtractClaims(ctx, token); err != nil {
				lib_log.Warn(ctx, "Unauthorized due to failure verifying gcp or api token for pubsub push, will log request body and then respond with No Content so that message is ACKed", lib_log.FmtError(err))
				if err := logRequestBody(r); err != nil {
					lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed logging request body"))
					return
				}
				lib_http.RenderNoContent(ctx, w)
				return
			}

		} else if err := c.tokenGcpClient.VerifyForPubsubPush(ctx, token); err != nil {
			if lib_errors.IsCustomWithCode(err, http.StatusBadGateway) || lib_errors.IsCustomWithCode(err, http.StatusGatewayTimeout) || lib_errors.IsCustomWithCode(err, http.StatusServiceUnavailable) {
				lib_http.RenderError(ctx, w, err)
				return
			}
			lib_log.Warn(ctx, "Unauthorized due to failure verifying gcp token for pubsub push, will log request body and then respond with No Content so that message is ACKed", lib_log.FmtError(err))
			if err := logRequestBody(r); err != nil {
				lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed logging request body"))
				return
			}
			lib_http.RenderNoContent(ctx, w)
			return
		}

		var pubsubMessageId string
		var pubsubMessagePublishTime time.Time
		if lib_context.IntegrationTest(ctx) && !lib_context.IntegrationTestPubsubAutoAckDisable(ctx) {
			pubsubMessageId = r.Header.Get(lib_http.HeaderKeyXLcPubsubMessageId)
			pubsubMessagePublishTime, err = time.Parse(time.RFC3339Nano, r.Header.Get(lib_http.HeaderKeyXLcPubsubMessagePublishTime))
			if err != nil {
				lib_log.Warn(ctx, "Unauthorized due to failure parsing pubsub publish time as a time value, will log request body and then respond with No Content so that message is ACKed", lib_log.FmtError(err))
				if err := logRequestBody(r); err != nil {
					lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed logging request body"))
					return
				}
				lib_http.RenderNoContent(ctx, w)
				return
			}

		} else {
			body, err := lib_http.ReadRequestBodyWithLogRedactions(r, true, jsonRedactionsByPathForMessage)
			if err != nil {
				lib_http.RenderError(ctx, w, lib_errors.Wrap(err, "Failed reading request body with log redactions"))
				return
			}

			var pushMessage *pushMessage
			var data []byte
			ctx, pushMessage, data, err = c.unmarshalFromPush(ctx, body)
			if err != nil {
				lib_log.Warn(ctx, "Unauthorized due to failure unmarshalling from push, will respond with No Content so that message is ACKed", lib_log.FmtError(err))
				lib_http.RenderNoContent(ctx, w)
				return
			}

			if pushMessage.Message.Attributes.CreatedInIntegrationTest != "" {
				createdInIntegrationTest, err := strconv.ParseBool(pushMessage.Message.Attributes.CreatedInIntegrationTest)
				if err != nil {
					lib_log.Warn(ctx, "Unauthorized due to failure parsing created in integration test attribute as a bool value, will respond with No Content so that message is ACKed", lib_log.FmtError(err))
					lib_http.RenderNoContent(ctx, w)
					return
				}

				if createdInIntegrationTest {
					if lib_context.IntegrationTestPubsubAutoAckDisable(ctx) {
						lib_log.Info(ctx, "Message created in integration test and disable auto ACK is enabled, the message will be processed")
					} else {
						lib_log.Info(ctx, "Message created in integration test and disable auto ACK is NOT enabled, will respond with No Content so that message is ACKed")
						lib_http.RenderNoContent(ctx, w)
						return
					}
				}
			}

			logPotentiallyLargeBytes(ctx, data, "Replacing request body with message data")
			r.Body = io.NopCloser(bytes.NewBuffer(data))
			pubsubMessageId = pushMessage.Message.MessageId
			pubsubMessagePublishTime = pushMessage.Message.PublishTime
		}
		ctx = lib_context.WithPubsubMessageId(ctx, pubsubMessageId)
		ctx = lib_context.WithPubsubMessagePublishTime(ctx, pubsubMessagePublishTime)

		if time.Now().AddDate(0, 0, -6).After(pubsubMessagePublishTime) {
			lib_log.Error(ctx, "Pubsub message was published more than 6 days ago, pubsub message maximum expiry is 7 days, this pubsub message might expire soon", lib_log.FmtTime("pubsubMessagePublishTime", pubsubMessagePublishTime))
		}

		lib_log.Info(ctx, "Authorized")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
	//revive:enable:cyclomatic
}

func logRequestBody(r *http.Request) error {
	ctx := r.Context()

	body, err := lib_http.ReadRequestBody(r, true)
	if err != nil {
		return lib_errors.Wrap(err, "Failed reading request body")
	}

	logPotentiallyLargeBytes(ctx, body, "Request body")

	return nil
}

func logPotentiallyLargeBytes(ctx context.Context, bytes []byte, message string) {
	bytesLength := len(bytes)
	maxLogPayloadSize := lib_log.MaxPayloadSize()
	if bytesLength < maxLogPayloadSize {
		lib_log.Info(ctx, message, lib_log.FmtBytes("bytes", bytes))

	} else {
		chunkSize := maxLogPayloadSize / 2
		chunks := (bytesLength + chunkSize - 1) / chunkSize
		for i := 0; i < chunks; i++ {
			chunkStart := chunkSize * i
			chunkEnd := chunkSize * (i + 1)
			if chunkEnd > bytesLength {
				chunkEnd = bytesLength
			}
			lib_log.Info(ctx, fmt.Sprintf("%s part %d/%d", message, i+1, chunks), lib_log.FmtBytes("bytes[chunkStart:chunkEnd]", bytes[chunkStart:chunkEnd]))
		}
	}
}

type pushMessage struct {
	Message      message `json:"message"`
	Subscription string  `json:"subscription"`
}

type message struct {
	Attributes  attributes `json:"attributes"`
	Data        string     `json:"data"`
	MessageId   string     `json:"message_id"`
	PublishTime time.Time  `json:"publish_time"`
}

type attributes struct {
	CorrelationId                          string  `json:"correlation_id"`
	CreatedInIntegrationTest               string  `json:"created_in_integration_test"`
	CreatedInIntegrationTestAutoAckDisable *string `json:"created_in_integration_test_auto_ack_disable,omitempty"`
	Test                                   string  `json:"test"`
}

func (c client) unmarshalFromPush(ctx context.Context, requestBody []byte) (newCtx context.Context, pushMessage *pushMessage, data []byte, err error) {
	pushMessage, data, err = unmarshalWithNoAttributesFromPush(ctx, requestBody)
	if err != nil {
		newCtx = ctx
		err = lib_errors.Wrap(err, "Failed unmarshalling with no attributes from push")
		return
	}

	ctx = lib_context.WithCorrelationIdAppend(ctx, pushMessage.Message.Attributes.CorrelationId)
	if pushMessage.Message.Attributes.Test != "" {
		var test bool
		test, err = strconv.ParseBool(pushMessage.Message.Attributes.Test)
		if err != nil {
			newCtx = ctx
			err = lib_errors.Wrap(err, "Failed parsing Test attribute of pubsub message as bool")
			return
		}
		ctx = lib_context.WithTest(ctx, test || lib_context.Test(ctx))
	}

	if pushMessage.Message.Attributes.CreatedInIntegrationTest != "" {
		var createdInIntegrationTest bool
		createdInIntegrationTest, err = strconv.ParseBool(pushMessage.Message.Attributes.CreatedInIntegrationTest)
		if err != nil {
			newCtx = ctx
			err = lib_errors.Wrap(err, "Failed parsing CreatedInIntegrationTest attribute of pubsub message as bool")
			return
		}
		integrationTest := createdInIntegrationTest || lib_context.IntegrationTest(ctx)
		ctx = lib_context.WithIntegrationTest(ctx, integrationTest)

		if integrationTest {
			if pushMessage.Message.Attributes.CreatedInIntegrationTestAutoAckDisable != nil {
				var createdInIntegrationTestAutoAckDisable bool
				createdInIntegrationTestAutoAckDisable, err = strconv.ParseBool(*pushMessage.Message.Attributes.CreatedInIntegrationTestAutoAckDisable)
				if err != nil {
					newCtx = ctx
					err = lib_errors.Wrap(err, "Failed parsing CreatedInIntegrationTestAutoAckDisable attribute of pubsub message as bool")
					return
				}
				ctx = lib_context.WithIntegrationTestPubsubAutoAckDisable(ctx, createdInIntegrationTestAutoAckDisable || lib_context.IntegrationTestPubsubAutoAckDisable(ctx))
			}
		}
	}

	newCtx = ctx
	lib_log.Info(newCtx, "Added", lib_log.FmtString("lib_context.CorrelationId(newCtx)", lib_context.CorrelationId(newCtx)), lib_log.FmtBool("lib_context.Test(newCtx)", lib_context.Test(newCtx)))
	return
}

func unmarshalWithNoAttributesFromPush(ctx context.Context, requestBody []byte) (incomingPushMessage *pushMessage, data []byte, err error) {
	lib_log.Info(ctx, "Parsing", lib_log.FmtBytes("requestBody", requestBody))

	var p pushMessage
	if err = json.Unmarshal(requestBody, &p); err != nil {
		err = lib_errors.Wrap(err, "Failed unmarshalling into PushMessage")
		return
	}
	incomingPushMessage = &p

	data, err = base64.StdEncoding.DecodeString(incomingPushMessage.Message.Data)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed decoding message data as base64 string")
		return
	}

	lib_log.Info(ctx, "Parsed", lib_log.FmtAny("incomingPushMessage", incomingPushMessage), lib_log.FmtBytes("data", data))
	return
}

func (c client) CreateTopic(ctx context.Context, topicId string) (*pubsub.Topic, error) {
	lib_log.Info(ctx, "Creating", lib_log.FmtString("topicId", topicId))

	topic, err := c.createTopic(ctx, topicId)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed creating topic")
	}

	lib_log.Info(ctx, "Created")
	return topic, nil
}

func (c client) CreateTopics(ctx context.Context, topicIds []string) (topicByTopicId map[string]*pubsub.Topic, err error) {
	lib_log.Info(ctx, "Retrieving", lib_log.FmtStrings("topicIds", topicIds))

	if len(topicIds) == 0 {
		return
	}
	type topicAndNameAndError struct {
		Err     error
		Topic   *pubsub.Topic
		TopicId string
	}
	ch := make(chan topicAndNameAndError)
	topicByTopicId = make(map[string]*pubsub.Topic)
	for _, topicId := range topicIds {
		go func(topicId string) {
			topic, err := c.createTopic(ctx, topicId)
			if err != nil {
				ch <- topicAndNameAndError{
					Err: lib_errors.Wrap(err, "Failed creating topic"),
				}
				return
			}
			ch <- topicAndNameAndError{
				Topic:   topic,
				TopicId: topicId,
			}
		}(topicId)
	}
	for range topicIds {
		topicAndNameAndError := <-ch
		if topicAndNameAndError.Err != nil {
			lib_log.Error(ctx, "Failed creating topic in goroutine", lib_log.FmtError(topicAndNameAndError.Err))
			err = lib_errors.Wrap(topicAndNameAndError.Err, "Failed creating topic in goroutine")
		}
		topicByTopicId[topicAndNameAndError.TopicId] = topicAndNameAndError.Topic
	}
	if err != nil {
		err = lib_errors.Wrap(err, "Failed creating one or more topics")
		return
	}

	lib_log.Info(ctx, "Retrieved", lib_log.FmtInt("len(topicByTopicId)", len(topicByTopicId)))
	return
}

func (c client) createTopic(ctx context.Context, topicId string) (*pubsub.Topic, error) {
	topic, err := c.cloudPubsubClient.GetTopic(ctx, topicId)
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			return nil, lib_errors.Wrap(err, "Failed getting topic")
		}
		lib_log.Info(ctx, "Topic not found, will attempt to create", lib_log.FmtError(err))
		topic, err = c.cloudPubsubClient.CreateTopic(ctx, topicId)
		if err != nil {
			return nil, lib_errors.Wrap(err, "Failed creating topic")
		}
	}
	return topic, nil
}

func (c client) Topic(ctx context.Context, topicId string) (*pubsub.Topic, error) {
	lib_log.Info(ctx, "Retrieving", lib_log.FmtString("topicId", topicId))

	topic, err := c.topic(ctx, topicId)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed retrieving topic")
	}

	lib_log.Info(ctx, "Retrieved")
	return topic, nil
}

func (c client) Topics(ctx context.Context, topicIds []string) (topicByTopicId map[string]*pubsub.Topic, err error) {
	lib_log.Info(ctx, "Retrieving", lib_log.FmtStrings("topicIds", topicIds))

	if len(topicIds) == 0 {
		return
	}
	type topicAndNameAndError struct {
		Err     error
		Topic   *pubsub.Topic
		TopicId string
	}
	ch := make(chan topicAndNameAndError)
	topicByTopicId = make(map[string]*pubsub.Topic)
	for _, topicId := range topicIds {
		go func(topicId string) {
			topic, err := c.topic(ctx, topicId)
			if err != nil {
				ch <- topicAndNameAndError{
					Err: lib_errors.Wrap(err, "Failed retrieving topic"),
				}
				return
			}
			ch <- topicAndNameAndError{
				Topic:   topic,
				TopicId: topicId,
			}
		}(topicId)
	}
	for range topicIds {
		topicAndNameAndError := <-ch
		if topicAndNameAndError.Err != nil {
			lib_log.Error(ctx, "Failed getting topic in goroutine", lib_log.FmtError(topicAndNameAndError.Err))
			err = lib_errors.Wrap(topicAndNameAndError.Err, "Failed getting topic in goroutine")
		}
		topicByTopicId[topicAndNameAndError.TopicId] = topicAndNameAndError.Topic
	}
	if err != nil {
		err = lib_errors.Wrap(err, "Failed getting one or more topics")
		return
	}

	lib_log.Info(ctx, "Retrieved", lib_log.FmtInt("len(topicByTopicId)", len(topicByTopicId)))
	return
}

func (c client) topic(ctx context.Context, topicId string) (*pubsub.Topic, error) {
	topic, err := c.cloudPubsubClient.GetTopic(ctx, topicId)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed getting topic")
	}
	return topic, nil
}

func (c client) PublishMessage(ctx context.Context, topic *pubsub.Topic, data []byte) (messageId string, err error) {
	lib_log.Info(ctx, "Publishing message", lib_log.FmtString("topic.String()", topic.String()), lib_log.FmtBytes("data", data))

	msg := pubsub.Message{Data: data, Attributes: messageAttributes(ctx)}
	lib_log.Info(ctx, "Constructed", lib_log.FmtString("msg.ID", msg.ID), lib_log.FmtTime("msg.PublishTime", msg.PublishTime), lib_log.FmtAny("msg.Attributes", msg.Attributes), lib_log.FmtString("string(msg.Data))", string(msg.Data)))

	result := topic.Publish(ctx, &msg)

	messageId, err = result.Get(ctx)
	if err != nil {
		lib_log.Error(ctx, "Error publishing message", lib_log.FmtString("messageId", messageId), lib_log.FmtError(err))
		err = lib_errors.Wrapf(err, "Error publishing message '%s' to topic '%s'", messageId, topic.String())
		return
	}

	lib_log.Info(ctx, "Published message", lib_log.FmtString("messageId", messageId))
	return
}

func (c client) PublishMessages(ctx context.Context, topic *pubsub.Topic, datas [][]byte) (messageIds []string, err error) {
	lib_log.Info(ctx, "Publishing messages", lib_log.FmtString("topic.String()", topic.String()), lib_log.FmtInt("len(datas)", len(datas)))

	var results []*pubsub.PublishResult
	attr := messageAttributes(ctx)
	for _, data := range datas {
		msg := pubsub.Message{Data: data, Attributes: attr}
		lib_log.Info(ctx, "Constructed", lib_log.FmtAny("msg.Attributes", msg.Attributes), lib_log.FmtString("string(msg.Data))", string(msg.Data)))

		results = append(results, topic.Publish(ctx, &msg))
	}

	var errs []error
	for _, v := range results {
		var messageId string
		messageId, err := v.Get(ctx)
		if err != nil {
			lib_log.Error(ctx, "Error publishing message", lib_log.FmtString("messageId", messageId), lib_log.FmtError(err))
			errs = append(errs, err)

		} else {
			lib_log.Info(ctx, "Published message", lib_log.FmtString("messageId", messageId))
			messageIds = append(messageIds, messageId)
		}
	}

	if len(errs) > 0 {
		var errMsgs []string
		for _, err := range errs {
			errMsgs = append(errMsgs, err.Error())
		}
		err = lib_errors.Errorf("Error publishing messages to topic '%s': %s", topic.String(), strings.Join(errMsgs, " | "))
		return
	}

	lib_log.Info(ctx, "Published messages", lib_log.FmtInt("len(messageIds)", len(messageIds)))
	return
}

const (
	messageAttributeCorrelationId                          = "correlation_id"
	messageAttributeCreatedInIntegrationTest               = "created_in_integration_test"
	messageAttributeCreatedInIntegrationTestAutoAckDisable = "created_in_integration_test_auto_ack_disable"
	messageAttributeTest                                   = "test"
)

func messageAttributes(ctx context.Context) map[string]string {
	return map[string]string{
		messageAttributeCorrelationId:                          fmt.Sprintf("%s,%s", lib_context.PubsubCorrelationId, lib_context.CorrelationId(ctx)),
		messageAttributeCreatedInIntegrationTest:               strconv.FormatBool(lib_context.IntegrationTest(ctx)),
		messageAttributeCreatedInIntegrationTestAutoAckDisable: strconv.FormatBool(lib_context.IntegrationTestPubsubAutoAckDisable(ctx)),
		messageAttributeTest:                                   strconv.FormatBool(lib_context.Test(ctx)),
	}
}

func (c client) Unmarshal(ctx context.Context, topic, subscription string, msg *pubsub.Message, dst interface{}) (newCtx context.Context, createdInIntegrationTest, createdInIntegrationTestAutoAckDisable bool, err error) {
	newCtx = ctx

	if err = c.UnmarshalWithNoAttributes(newCtx, topic, subscription, msg, dst); err != nil {
		err = lib_errors.Wrap(err, "Failed unmarshalling with no attributes")
		return
	}

	newCtx = lib_context.WithCorrelationIdAppend(newCtx, msg.Attributes[messageAttributeCorrelationId])
	test, err := strconv.ParseBool(msg.Attributes[messageAttributeTest])
	if err != nil {
		err = lib_errors.Wrapf(err, "Failed parsing attribute %q as bool", messageAttributeTest)
		return
	}
	newCtx = lib_context.WithTest(newCtx, test || lib_context.Test(newCtx))

	if v, ok := msg.Attributes[messageAttributeCreatedInIntegrationTest]; ok {
		createdInIntegrationTest, err = strconv.ParseBool(v)
		if err != nil {
			err = lib_errors.Wrapf(err, "Failed parsing attribute %q (%q) as bool", messageAttributeCreatedInIntegrationTest, v)
			return
		}
	}

	if v, ok := msg.Attributes[messageAttributeCreatedInIntegrationTestAutoAckDisable]; ok {
		createdInIntegrationTestAutoAckDisable, err = strconv.ParseBool(v)
		if err != nil {
			err = lib_errors.Wrapf(err, "Failed parsing attribute %q (%q) as bool", messageAttributeCreatedInIntegrationTestAutoAckDisable, v)
			return
		}
	}

	lib_log.Info(newCtx, "Added", lib_log.FmtString("lib_context.CorrelationId(newCtx)", lib_context.CorrelationId(newCtx)), lib_log.FmtBool("lib_context.Test(newCtx)", lib_context.Test(newCtx)), lib_log.FmtBool("createdInIntegrationTest", createdInIntegrationTest), lib_log.FmtBool("createdInIntegrationTestAutoAckDisable", createdInIntegrationTestAutoAckDisable))
	return
}

func (c client) UnmarshalWithNoAttributes(ctx context.Context, topic, subscription string, msg *pubsub.Message, dst interface{}) error {
	lib_log.Info(ctx, "Unmarshalling", lib_log.FmtString("topic", topic), lib_log.FmtString("subscription", subscription), lib_log.FmtString("msg.ID", msg.ID), lib_log.FmtTime("msg.PublishTime", msg.PublishTime), lib_log.FmtAny("msg.Attributes", msg.Attributes), lib_log.FmtString("string(msg.Data))", string(msg.Data)))

	if err := json.Unmarshal(msg.Data, dst); err != nil {
		return lib_errors.Wrap(err, "Failed decoding message")
	}

	lib_log.Info(ctx, "Unmarshalled", lib_log.FmtAny("dst", dst))
	return nil
}
