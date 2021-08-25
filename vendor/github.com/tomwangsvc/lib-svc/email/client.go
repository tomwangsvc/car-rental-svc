package email

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	lib_brand "github.com/tomwangsvc/lib-svc/brand"
	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_pubsub "github.com/tomwangsvc/lib-svc/pubsub"
	lib_regexp "github.com/tomwangsvc/lib-svc/regexp"
	lib_secrets "github.com/tomwangsvc/lib-svc/secrets"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	AttachmentFileSizeLimit   = int64(10485760)
	TotalAttachmentsSizeLimit = int64(26214400)
)

var (
	emailRegExpCompile                 *regexp.Regexp
	emailWithOptionalNameRegExpCompile *regexp.Regexp
)

func init() {
	emailRegExpCompile = regexp.MustCompile(lib_regexp.EmailRegExp)
	emailWithOptionalNameRegExpCompile = regexp.MustCompile(lib_regexp.EmailWithOptionalNameRegExp)
}

type Client interface {
	SendViaEmailSvc(ctx context.Context, subject, plainTextContent, htmlContent, emailCategory string, ccEmails, bccEmails, storageObjectIds []string, fromDomain DomainWithRequiredReference, toDomains []DomainWithRequiredReference, relatedDomains []Domain, system, seperateToEmails bool) error
}

type Config struct {
	Env lib_env.Env
}

func NewClient(ctx context.Context, config Config, cloudPubsubClient *pubsub.Client, pubsubClient lib_pubsub.Client, secretsClient lib_secrets.Client) (Client, error) {
	lib_log.Info(ctx, "Initializing", lib_log.FmtAny("config", config))

	sendGridApiKey, err := secretsClient.ValueFromBase64WithNewLinesStripped(secretDomain, secretTypeApiKey)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed retrieving sengrid api key from secrets")
	}

	topic, err := cloudPubsubClient.CreateTopic(ctx, TopicNameEmailSvcSendEmail)
	if err != nil {
		if grpc.Code(err) != codes.AlreadyExists {
			return nil, lib_errors.Wrap(err, "Failed creating subscription")
		}
		lib_log.Info(ctx, "Ignoring error from `cloudPubsubClient.CreateSubscription`. We will always get 'already exists', if there are other fatal errors then the first use will fail and errors can be handled there", lib_log.FmtError(err))
		topic = cloudPubsubClient.Topic(TopicNameEmailSvcSendEmail)
		if topic == nil {
			return nil, lib_errors.Wrapf(err, "Failed getting reference to topic %q after topic creation failed", TopicNameEmailSvcSendEmail)
		}
	}

	lib_log.Info(ctx, "Initialized")
	return client{
		config:                 config,
		topicEmailSvcSendEmail: topic,
		pubsubClient:           pubsubClient,
		sendGridApiKey:         sendGridApiKey,
	}, nil
}

const (
	emailAttachmentsBytesSizeLimit = 1024 * 1024 * 2 // 10MB (the maximum request size) but we limit to 2MB
	secretDomain                   = "sendgrid"
	secretTypeApiKey               = "api-key"
	TopicNameEmailSvcSendEmail     = "email-svc_send-email"
	TopicNameEmailSvcReceiveEmail  = "email-svc_receive-email-notification"
)

func RequiredSecrets() lib_secrets.Required {
	return lib_secrets.Required{
		secretDomain: []string{
			secretTypeApiKey,
		},
	}
}

type client struct {
	config                 Config
	topicEmailSvcSendEmail *pubsub.Topic
	pubsubClient           lib_pubsub.Client
	sendGridApiKey         string
}

type Attachment struct {
	Content  []byte
	FileName string
	MimeType string
}

const (
	tokenExpiry = 86400 * time.Second
)

type Domain struct {
	Id        *string `json:"id,omitempty"`
	Reference *string `json:"reference,omitempty"`
	Type      string  `json:"type"`
}

type MlDomain struct {
	Id          *string `json:"id,omitempty"`
	Probability float64 `json:"probability"`
	Reference   *string `json:"reference,omitempty"`
	Type        string  `json:"type"`
}

type DomainWithRequiredReference struct {
	Id        *string `json:"id,omitempty"`
	Reference string  `json:"reference"`
	Type      string  `json:"type"`
}

func (d Domain) Map() map[string]interface{} {
	m := map[string]interface{}{
		"type": d.Type,
	}
	if d.Id != nil {
		m["id"] = *d.Id
	}
	if d.Reference != nil {
		m["reference"] = *d.Reference
	}

	return m
}

//revive:disable:argument-limit
func (c client) SendViaEmailSvc(ctx context.Context, subject, plainTextContent, htmlContent, emailCategory string, ccEmails, bccEmails, storageObjectIds []string, fromDomain DomainWithRequiredReference, toDomains []DomainWithRequiredReference, relatedDomains []Domain, system bool, seperateToEmails bool) error {
	lib_log.Info(ctx, "Sending",
		lib_log.FmtString("subject", subject),
		lib_log.FmtInt("len(plainTextContent)", len(plainTextContent)),
		lib_log.FmtInt("len(htmlContent)", len(htmlContent)),
		lib_log.FmtString("category", emailCategory),
		lib_log.FmtStrings("ccEmails", ccEmails),
		lib_log.FmtStrings("bccEmails", bccEmails),
		lib_log.FmtAny("fromDomain", fromDomain),
		lib_log.FmtAny("toDomain", toDomains),
		lib_log.FmtInt("len(relatedDomains)", len(relatedDomains)),
		lib_log.FmtInt("len(storageObjectIds)", len(storageObjectIds)),
		lib_log.FmtBool("system", system),
		lib_log.FmtBool("seperateToEmails", seperateToEmails),
	)

	if seperateToEmails {
		var seperatedEmailBytes [][]byte
		var seperatedEmails []email

		for _, v := range toDomains {
			seperatedEmail := initializeEmail(subject, plainTextContent, htmlContent, emailCategory, ccEmails, bccEmails, storageObjectIds, fromDomain, []DomainWithRequiredReference{v}, relatedDomains, system)

			if seperatedEmail != nil {
				seperatedEmails = append(seperatedEmails, *seperatedEmail)
			}
		}
		for _, v := range seperatedEmails {
			emailBytes, err := json.Marshal(v)
			if err != nil {
				return lib_errors.Wrap(err, "Failed marshalling email into bytes")
			}
			seperatedEmailBytes = append(seperatedEmailBytes, emailBytes)
		}

		if _, err := c.pubsubClient.PublishMessages(ctx, c.topicEmailSvcSendEmail, seperatedEmailBytes); err != nil {
			return lib_errors.Wrapf(err, "Failed publishing email messages to topic %q", TopicNameEmailSvcSendEmail)
		}
	} else {
		email := initializeEmail(subject, plainTextContent, htmlContent, emailCategory, ccEmails, bccEmails, storageObjectIds, fromDomain, toDomains, relatedDomains, system)

		emailBytes, err := json.Marshal(email)
		if err != nil {
			return lib_errors.Wrap(err, "Failed marshalling email into bytes")
		}

		if _, err := c.pubsubClient.PublishMessage(ctx, c.topicEmailSvcSendEmail, emailBytes); err != nil {
			return lib_errors.Wrapf(err, "Failed publishing email message to topic %q", TopicNameEmailSvcSendEmail)
		}
	}
	lib_log.Info(ctx, "Sent")

	return nil
	//revive:enable:argument-limit
}

//revive:disable:argument-limit
func initializeEmail(subject, plainTextContent, htmlContent, emailCategory string, ccEmails, bccEmails, storageObjectIds []string, fromDomain DomainWithRequiredReference, toDomains []DomainWithRequiredReference, relatedDomains []Domain, system bool) *email {
	return &email{
		BccEmails:        bccEmails,
		CcEmails:         ccEmails,
		FromDomain:       fromDomain,
		HtmlContent:      htmlContent,
		Category:         emailCategory,
		PlainTextContent: plainTextContent,
		RelatedDomains:   relatedDomains,
		StorageObjectIds: storageObjectIds,
		Subject:          subject,
		System:           system,
		ToDomains:        toDomains,
	}
	//revive:enable:argument-limit
}

type email struct {
	BccEmails        []string                      `json:"bcc_emails,omitempty"`
	Category         string                        `json:"category"`
	CcEmails         []string                      `json:"cc_emails,omitempty"`
	FromDomain       DomainWithRequiredReference   `json:"from_domain"`
	HtmlContent      string                        `json:"html_content"`
	PlainTextContent string                        `json:"plain_text_content"`
	RelatedDomains   []Domain                      `json:"related_domains,omitempty"`
	StorageObjectIds []string                      `json:"storage_object_ids,omitempty"`
	Subject          string                        `json:"subject"`
	System           bool                          `json:"system"`
	ToDomains        []DomainWithRequiredReference `json:"to_domains"`
}

func (e email) Map() map[string]interface{} {
	m := map[string]interface{}{
		"category":           e.Category,
		"from_domain":        e.FromDomain,
		"html_content":       lib_log.TruncateField(e.HtmlContent),
		"plain_text_content": lib_log.TruncateField(e.PlainTextContent),
		"subject":            e.Subject,
		"system":             e.System,
	}
	if len(e.BccEmails) > 0 {
		m["bcc_emails"] = e.BccEmails
	}
	if len(e.CcEmails) > 0 {
		m["cc_emails"] = e.CcEmails
	}
	if len(e.StorageObjectIds) > 0 {
		m["storage_object_ids"] = e.StorageObjectIds
	}
	if len(e.RelatedDomains) > 0 {
		m["related_domains"] = e.RelatedDomains
	}
	if len(e.ToDomains) > 0 {
		m["to_domains"] = e.ToDomains
	}
	return m
}

type Email struct {
	Value string  `json:"value"`
	Name  *string `json:"name,omitempty"`
}

// TODO: This needs to return 'csr-alerts@goboxer.com' or 'csr-alerts-dev|uat@goboxer.com' when the brand is Boxer and that email is available
func CsrAlertsEmail(ctx context.Context, env lib_env.Env, brandId string) (*Email, error) {
	name, err := csrAlertsEmailName(env, brandId)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Fail generating csr alerts email name")
	}
	if env.Production() && !lib_context.Test(ctx) {
		email := NewEmailWithOptionalName("csr-alerts@tomwang.com", &name)
		return &email, nil
	}
	email := NewEmailWithOptionalName("developer-test-user@tomwang.com", &name)
	return &email, nil
}

func csrAlertsEmailName(env lib_env.Env, brandId string) (string, error) {
	brandName, err := lib_brand.Name(brandId)
	if err != nil {
		return "", lib_errors.Wrap(err, "Brand ID not recognized")
	}
	if env.Prd() {
		return fmt.Sprintf("%s CSR Alerts", brandName), nil
	}
	return fmt.Sprintf("%s CSR Alerts %s", brandName, strings.ToUpper(env.Id)), nil
}

// TODO: This needs to return 'support@goboxer.com' or 'support-dev|uat@goboxer.com' when the brand is Boxer and that email is available
func SupportEmail(env lib_env.Env, brandId string) (*Email, error) {
	name, err := supportEmailName(env, brandId)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Fail generating csr email name")
	}
	if env.Production() {
		email := NewEmailWithOptionalName("support@tomwang.com", &name)
		return &email, nil
	}
	email := NewEmailWithOptionalName(fmt.Sprintf("support-%s@tomwang.com", env.Id), &name)
	return &email, nil
}

func supportEmailName(env lib_env.Env, brandId string) (string, error) {
	brandName, err := lib_brand.Name(brandId)
	if err != nil {
		return "", lib_errors.Wrap(err, "Brand ID not recognized")
	}
	if env.Prd() {
		return fmt.Sprintf("%s Support", brandName), nil
	}
	return fmt.Sprintf("%s Support %s", brandName, strings.ToUpper(env.Id)), nil
}

func NoReplyEmail(env lib_env.Env, brandId string) (*Email, error) {
	brand, emailDomain, err := noReplyEmailBrandAndEmailDomain(env, brandId)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Fail generating no reply email brand and email name")
	}

	var emailValue string
	if env.Production() {
		emailValue = fmt.Sprintf("noreply@%s", emailDomain)

	} else {
		emailValue = fmt.Sprintf("noreply-%s@%s", env.Id, emailDomain)
	}
	email := NewEmailWithOptionalName(emailValue, &brand)

	return &email, nil
}

func noReplyEmailBrandAndEmailDomain(env lib_env.Env, brandId string) (brand, domain string, err error) {
	brand, err = lib_brand.Name(brandId)
	if err != nil {
		return "", "", lib_errors.Wrap(err, "Brand ID not recognized")
	}

	if !env.Prd() {
		brand = fmt.Sprintf("%s %s", brand, strings.ToUpper(env.Id))
	}

	domain, err = lib_brand.EmailDomain(brandId)
	if err != nil {
		return "", "", lib_errors.Wrap(err, "Brand ID not recognized")
	}

	return brand, domain, nil
}

func NewEmailWithOptionalName(value string, name *string) Email {
	if name != nil {
		n := strings.TrimSpace(*name)
		if n != "" {
			return Email{
				Value: strings.ToLower(value),
				Name:  &n,
			}
		}
		return Email{
			Value: strings.ToLower(value),
		}
	}
	return Email{
		Value: strings.ToLower(value),
	}
}

func NewEmailFromEmailWithOptionalName(e string) Email {
	s := strings.Split(e, "<")
	if len(s) == 0 {
		return Email{
			Value: e,
		}
	}

	name := strings.TrimSpace(strings.Join(s[:len(s)-1], "<"))
	if name == "" {
		return Email{
			Value: strings.TrimSuffix(s[len(s)-1], ">"),
		}
	}

	return Email{
		Value: strings.TrimSuffix(s[len(s)-1], ">"),
		Name:  &name,
	}
}

func (e Email) String() string {
	if e.Name == nil || strings.TrimSpace(*e.Name) == "" {
		return e.Value
	}
	return fmt.Sprintf("%s <%s>", strings.TrimSpace(*e.Name), e.Value)
}

func ShouldAddEnvForEmailAddress(email string) bool {
	if !IsSupportEmail(email) && !IsNoReplyEmail(email) {
		return true
	}
	return false
}

func IsNoReplyEmail(email string) bool {
	return IsNoReplyEmailForBoxer(email) || IsNoReplyEmailForTomWang(email)
}

func IsNoReplyEmailForBoxer(email string) bool {
	return email == fmt.Sprintf("noreply-%s@goboxer.com", lib_env.Dev) ||
		email == fmt.Sprintf("noreply-%s@goboxer.com", lib_env.Stg) ||
		email == fmt.Sprintf("noreply-%s@goboxer.com", lib_env.Uat) ||
		email == "noreply@goboxer.com" ||
		strings.Contains(email, fmt.Sprintf("<noreply-%s@goboxer.com>", lib_env.Dev)) ||
		strings.Contains(email, fmt.Sprintf("<noreply-%s@goboxer.com>", lib_env.Stg)) ||
		strings.Contains(email, fmt.Sprintf("<noreply-%s@goboxer.com>", lib_env.Uat)) ||
		strings.Contains(email, "<noreply@goboxer.com>")
}

func IsNoReplyEmailForTomWang(email string) bool {
	return email == fmt.Sprintf("noreply-%s@tomwang.com", lib_env.Dev) ||
		email == fmt.Sprintf("noreply-%s@tomwang.com", lib_env.Stg) ||
		email == fmt.Sprintf("noreply-%s@tomwang.com", lib_env.Uat) ||
		email == "noreply@tomwang.com" ||
		strings.Contains(email, fmt.Sprintf("<noreply-%s@tomwang.com>", lib_env.Dev)) ||
		strings.Contains(email, fmt.Sprintf("<noreply-%s@tomwang.com>", lib_env.Stg)) ||
		strings.Contains(email, fmt.Sprintf("<noreply-%s@tomwang.com>", lib_env.Uat)) ||
		strings.Contains(email, "<noreply@tomwang.com>")
}

func IsSupportEmail(email string) bool {
	return isSupportEmailForBoxer(email) || isSupportEmailForTomWang(email)
}

func isSupportEmailForBoxer(email string) bool {
	return email == fmt.Sprintf("support-%s@goboxer.com", lib_env.Dev) ||
		email == fmt.Sprintf("support-%s@goboxer.com", lib_env.Stg) ||
		email == fmt.Sprintf("support-%s@goboxer.com", lib_env.Uat) ||
		email == "support@goboxer.com" ||
		strings.Contains(email, fmt.Sprintf("<support-%s@goboxer.com>", lib_env.Dev)) ||
		strings.Contains(email, fmt.Sprintf("<support-%s@goboxer.com>", lib_env.Stg)) ||
		strings.Contains(email, fmt.Sprintf("<support-%s@goboxer.com>", lib_env.Uat)) ||
		strings.Contains(email, "<support@goboxer.com>")
}

func isSupportEmailForTomWang(email string) bool {
	return email == fmt.Sprintf("support-%s@tomwang.com", lib_env.Dev) ||
		email == fmt.Sprintf("support-%s@tomwang.com", lib_env.Stg) ||
		email == fmt.Sprintf("support-%s@tomwang.com", lib_env.Uat) ||
		email == "support@tomwang.com" ||
		strings.Contains(email, fmt.Sprintf("<support-%s@tomwang.com>", lib_env.Dev)) ||
		strings.Contains(email, fmt.Sprintf("<support-%s@tomwang.com>", lib_env.Stg)) ||
		strings.Contains(email, fmt.Sprintf("<support-%s@tomwang.com>", lib_env.Uat)) ||
		strings.Contains(email, "<support@tomwang.com>")
}

func AddEnvForEmailAddressIfNecessary(domainReference string, env lib_env.Env) string {
	result := domainReference

	if !env.Production() {
		if ShouldAddEnvForEmailAddress(domainReference) {
			values := strings.Split(domainReference, "@")
			if len(values) > 1 && !strings.HasSuffix(values[0], fmt.Sprintf("+%s", env.Id)) {
				result = values[0] + fmt.Sprintf("+%s", env.Id) + "@" + values[1]
			}

		} else {
			values := strings.Split(domainReference, "@")
			if len(values) > 1 && !strings.HasSuffix(values[0], fmt.Sprintf("-%s", env.Id)) {
				result = values[0] + fmt.Sprintf("-%s", env.Id) + "@" + values[1]
			}
		}
	}
	return result
}

func IsEmailWithOptionalNameFormat(emailWithOptionalName string) bool {
	return emailWithOptionalNameRegExpCompile.MatchString(emailWithOptionalName)
}

func IsLowercase(value string) bool {
	return value == strings.ToLower(value)
}

func IsEmailWithOptionalNameLowercase(emailWithOptionalName string) bool {
	if v := IsEmailWithOptionalNameFormat(emailWithOptionalName); !v {
		return false
	}

	return IsLowercase(NewEmailFromEmailWithOptionalName(emailWithOptionalName).Value)
}

func IsEmailFormat(email string) bool {
	return emailRegExpCompile.MatchString(email)
}

func IsEmailLowercase(email string) bool {
	if v := IsEmailFormat(email); !v {
		return false
	}

	return IsLowercase(email)
}
