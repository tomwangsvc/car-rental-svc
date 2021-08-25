package context

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKeyType int

// Context constants
const (
	ContextKeyCloudSchedulerPush                  contextKeyType = iota
	ContextKeyCloudTaskCreatedDate                contextKeyType = iota
	ContextKeyCloudTaskId                         contextKeyType = iota
	ContextKeyCloudTasksPush                      contextKeyType = iota
	ContextKeyCorrelationId                       contextKeyType = iota
	ContextKeyIntegrationTest                     contextKeyType = iota
	ContextKeyIntegrationTestPubsubAutoAckDisable contextKeyType = iota
	ContextKeyLcCaller                            contextKeyType = iota
	ContextKeyPubsubMessageId                     contextKeyType = iota
	ContextKeyPubsubMessagePublishTime            contextKeyType = iota
	ContextKeyPubsubPush                          contextKeyType = iota
	ContextKeyHttpRequestHeaderIfNoneMatch        contextKeyType = iota
	ContextKeyTest                                contextKeyType = iota
	ContextKeyTestMetadata                        contextKeyType = iota
	ContextKeyXLcLocationHeaders                  contextKeyType = iota

	CleanUpCorrelationId        = "CLEAN_UP"
	CloudSchedulerCorrelationId = "CLOUD_SCHEDULER"
	CloudTasksCorrelationId     = "CLOUD_TASKS"
	MissingCorrelationId        = "MISSING"
	PubsubCorrelationId         = "PUB_SUB"
	StartUpCorrelationId        = "START_UP"
)

// NewContext creates a new context from an old context
// -> This new context is disconnected from the original, meaning it can continue after the original context is cancelled
// -> This is useful for goroutines that need to live after an http request is handled
func NewContext(ctx context.Context) context.Context {
	newContext := context.Background()
	newContext = WithCloudSchedulerPush(newContext, CloudSchedulerPush(ctx))
	newContext = WithCloudTaskCreatedDate(newContext, CloudTaskCreatedDate(ctx))
	newContext = WithCloudTaskId(newContext, CloudTaskId(ctx))
	newContext = WithCloudTasksPush(newContext, CloudTasksPush(ctx))
	newContext = WithCorrelationIdAppend(WithCorrelationId(newContext, uuid.New().String()), CorrelationId(ctx))
	newContext = WithIntegrationTest(newContext, IntegrationTest(ctx))
	newContext = WithLcCaller(newContext, LcCaller(ctx))
	newContext = WithIntegrationTestPubsubAutoAckDisable(newContext, IntegrationTestPubsubAutoAckDisable(ctx))
	newContext = WithPubsubMessageId(newContext, PubsubMessageId(ctx))
	newContext = WithPubsubMessagePublishTime(newContext, PubsubMessagePublishTime(ctx))
	newContext = WithPubsubPush(newContext, PubsubPush(ctx))
	newContext = WithTest(newContext, Test(ctx))
	newContext = WithTestMetadata(newContext, TestMetadata(ctx))
	newContext = WithXLcLocationHeaders(newContext, XLcLocationHeaders(ctx))
	newContext = iamNewContext(newContext, ctx)
	return newContext
}

// NewStartUpContext with start up correlation ID
func NewStartUpContext() context.Context {
	return WithCorrelationId(context.Background(), StartUpCorrelationId)
}

// NewCleanUpContext with clean up correlation ID
func NewCleanUpContext() context.Context {
	return WithCorrelationId(context.Background(), CleanUpCorrelationId)
}

// CloudSchedulerPush returns the CloudSchedulerPush bool of the ctx
func CloudSchedulerPush(ctx context.Context) bool {
	v, ok := ctx.Value(ContextKeyCloudSchedulerPush).(bool)
	if !ok {
		return false
	}
	return v
}

// WithCloudSchedulerPush creates a new context with CloudSchedulerPush bool
func WithCloudSchedulerPush(ctx context.Context, cloudSchedulerPush bool) context.Context {
	return context.WithValue(ctx, ContextKeyCloudSchedulerPush, cloudSchedulerPush)
}

// CloudTaskCreatedDate returns the CloudTaskCreatedDate string of the ctx
func CloudTaskCreatedDate(ctx context.Context) time.Time {
	v, ok := ctx.Value(ContextKeyCloudTaskCreatedDate).(time.Time)
	if !ok {
		return time.Time{}
	}
	return v
}

// WithCloudTaskCreatedDate creates a new context with CloudTaskCreatedDate time.Time
func WithCloudTaskCreatedDate(ctx context.Context, cloudTaskCreatedDate time.Time) context.Context {
	return context.WithValue(ctx, ContextKeyCloudTaskCreatedDate, cloudTaskCreatedDate)
}

// CloudTaskId returns the CloudTaskId string of the ctx
func CloudTaskId(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyCloudTaskId).(string)
	if !ok {
		return ""
	}
	return v
}

// WithCloudTaskId creates a new context with CloudTaskId string
func WithCloudTaskId(ctx context.Context, cloudTaskId string) context.Context {
	return context.WithValue(ctx, ContextKeyCloudTaskId, cloudTaskId)
}

// CloudTasksPush returns the CloudTasksPush bool of the ctx
func CloudTasksPush(ctx context.Context) bool {
	v, ok := ctx.Value(ContextKeyCloudTasksPush).(bool)
	if !ok {
		return false
	}
	return v
}

// WithCloudTasksPush creates a new context with CloudTasksPush bool
func WithCloudTasksPush(ctx context.Context, cloudTasksPush bool) context.Context {
	return context.WithValue(ctx, ContextKeyCloudTasksPush, cloudTasksPush)
}

// CorrelationId returns the correlation ID string of the ctx
func CorrelationId(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyCorrelationId).(string)
	if !ok {
		return ""
	}
	return v
}

// WithCorrelationId creates a new context with correlationId
func WithCorrelationId(ctx context.Context, correlationId string) context.Context {
	return context.WithValue(ctx, ContextKeyCorrelationId, correlationId)
}

// WithCorrelationIdAppend creates a new context with correlationId appended to existing
func WithCorrelationIdAppend(ctx context.Context, correlationId string) context.Context {
	c := CorrelationId(ctx)
	if c == "" {
		c = fmt.Sprintf("%s,%s", MissingCorrelationId, correlationId)
	} else if correlationId != "" {
		c = fmt.Sprintf("%s,%s", c, correlationId)
	}
	return WithCorrelationId(ctx, c)
}

// WithCorrelationIdPrepend creates a new context with correlationId prepended to existing
func WithCorrelationIdPrepend(ctx context.Context, correlationId string) context.Context {
	c := CorrelationId(ctx)
	if c == "" {
		c = fmt.Sprintf("%s,%s", correlationId, MissingCorrelationId)
	} else if correlationId != "" {
		c = fmt.Sprintf("%s,%s", correlationId, c)
	}
	return WithCorrelationId(ctx, c)
}

// IntegrationTest returns the Integration test bool of the ctx
func IntegrationTest(ctx context.Context) bool {
	v, ok := ctx.Value(ContextKeyIntegrationTest).(bool)
	if !ok {
		return false
	}
	return v
}

// WithIntegrationTest creates a new context with IntegrationTest bool
func WithIntegrationTest(ctx context.Context, integrationTest bool) context.Context {
	return context.WithValue(ctx, ContextKeyIntegrationTest, integrationTest)
}

// IntegrationTestPubsubAutoAckDisable returns the IntegrationTestPubsubAutoAckDisable bool of the ctx
func IntegrationTestPubsubAutoAckDisable(ctx context.Context) bool {
	v, ok := ctx.Value(ContextKeyIntegrationTestPubsubAutoAckDisable).(bool)
	if !ok {
		return false
	}
	return v
}

// WithIntegrationTestPubsubAutoAckDisable creates a new context with IntegrationTestPubsubAutoAckDisable string
func WithIntegrationTestPubsubAutoAckDisable(ctx context.Context, isIntegrationTestPubsubAutoAckDisable bool) context.Context {
	return context.WithValue(ctx, ContextKeyIntegrationTestPubsubAutoAckDisable, isIntegrationTestPubsubAutoAckDisable)
}

// LcCaller returns the TomWang referer bool of the ctx
func LcCaller(ctx context.Context) bool {
	v, ok := ctx.Value(ContextKeyLcCaller).(bool)
	if !ok {
		return false
	}
	return v
}

// WithLcCaller creates a new context with WithLcCaller bool
func WithLcCaller(ctx context.Context, isLcCaller bool) context.Context {
	return context.WithValue(ctx, ContextKeyLcCaller, isLcCaller)
}

// PubsubMessageId returns the PubsubMessageId string of the ctx
func PubsubMessageId(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyPubsubMessageId).(string)
	if !ok {
		return ""
	}
	return v
}

// WithPubsubMessageId creates a new context with PubsubMessageId string
func WithPubsubMessageId(ctx context.Context, pubsubMessageId string) context.Context {
	return context.WithValue(ctx, ContextKeyPubsubMessageId, pubsubMessageId)
}

// PubsubMessagePublishTime returns the PubsubMessagePublishTime string of the ctx
func PubsubMessagePublishTime(ctx context.Context) time.Time {
	v, ok := ctx.Value(ContextKeyPubsubMessagePublishTime).(time.Time)
	if !ok {
		return time.Time{}
	}
	return v
}

// WithPubsubMessagePublishTime creates a new context with PubsubMessagePublishTime time.Time
func WithPubsubMessagePublishTime(ctx context.Context, cloudTaskCreatedDate time.Time) context.Context {
	return context.WithValue(ctx, ContextKeyPubsubMessagePublishTime, cloudTaskCreatedDate)
}

// PubsubPush returns the PubsubPush bool of the ctx
func PubsubPush(ctx context.Context) bool {
	v, ok := ctx.Value(ContextKeyPubsubPush).(bool)
	if !ok {
		return false
	}
	return v
}

// WithPubsubPush creates a new context with PubsubPush bool
func WithPubsubPush(ctx context.Context, pubsubPush bool) context.Context {
	return context.WithValue(ctx, ContextKeyPubsubPush, pubsubPush)
}

func HttpRequestHeaderIfNoneMatch(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyHttpRequestHeaderIfNoneMatch).(string)
	if !ok {
		return ""
	}
	return v
}

func WithHttpRequestHeaderIfNoneMatch(ctx context.Context, httpRequestHeaderIfNoneMatch string) context.Context {
	return context.WithValue(ctx, ContextKeyHttpRequestHeaderIfNoneMatch, httpRequestHeaderIfNoneMatch)
}

// Test returns the test bool of the ctx
func Test(ctx context.Context) bool {
	v, ok := ctx.Value(ContextKeyTest).(bool)
	if !ok {
		return false
	}
	return v
}

// WithTest creates a new context with test bool
func WithTest(ctx context.Context, test bool) context.Context {
	if !test && IntegrationTest(ctx) {
		ctx = WithIntegrationTest(ctx, false)
	}
	return context.WithValue(ctx, ContextKeyTest, test)
}

// TestMetadata returns the test metadata string of the ctx
func TestMetadata(ctx context.Context) string {
	v, ok := ctx.Value(ContextKeyTestMetadata).(string)
	if !ok {
		return ""
	}
	return v
}

// WithTestMetadata creates a new context with testMetadata
func WithTestMetadata(ctx context.Context, testMetadata string) context.Context {
	return context.WithValue(ctx, ContextKeyTestMetadata, testMetadata)
}

// XLcLocationHeaders returns the location headers of the ctx
func XLcLocationHeaders(ctx context.Context) http.Header {
	v, ok := ctx.Value(ContextKeyXLcLocationHeaders).(http.Header)
	if !ok {
		return nil
	}
	return v
}

// WithXLcLocationHeaders creates a new context with locations headers associated with the request
func WithXLcLocationHeaders(ctx context.Context, locationHeaders http.Header) context.Context {
	return context.WithValue(ctx, ContextKeyXLcLocationHeaders, locationHeaders)
}
