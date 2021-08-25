package log

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/logging"
	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
)

// See https://cloud.google.com/appengine/docs/flexible/go/writing-application-logs

const (
	red    = 31
	green  = 32
	yellow = 33
	cyan   = 36

	colorDebug = green
	colorInfo  = cyan
	colorWarn  = yellow
	colorError = red
	colorFatal = red

	missingCorrelationIdMsg = "MISSING_CORRELATION_ID"

	// We need some space for the message and meta-data
	stackdriverFieldSizeLimitBytes         = 98 * 1024
	stackdriverCloudRunFieldSizeLimitBytes = stackdriverFieldSizeLimitBytes // We used to use '4 * 1024' but apparently the limit has been increased

	fmtStringMaxLength = 2048
)

var (
	l *logger

	// JSON doesn't support newlines and although there are various conventions these may not work everywhere e.g. in BigQuery
	disablePrettyPrint bool

	isRunningInGae      bool
	isRunningInCloudRun bool
)

func init() {
	// We use the standard logger for 'Cloud Run' but it requires JSON and so me must remove the time prefix that is auto-added to log messages
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	l = newDefaultLogger()
}

func newDefaultLogger() *logger {
	return &logger{
		debugLogger: log.New(os.Stdout, fmt.Sprintf("\x1b[%dmDEBUG ", colorDebug), log.Ldate|log.Ltime|log.Lmicroseconds),
		infoLogger:  log.New(os.Stdout, fmt.Sprintf("\x1b[%dmINFO  ", colorInfo), log.Ldate|log.Ltime|log.Lmicroseconds),
		warnLogger:  log.New(os.Stderr, fmt.Sprintf("\x1b[%dmWARN  ", colorWarn), log.Ldate|log.Ltime|log.Lmicroseconds),
		errorLogger: log.New(os.Stderr, fmt.Sprintf("\x1b[%dmERROR ", colorError), log.Ldate|log.Ltime|log.Lmicroseconds),
		fatalLogger: log.New(os.Stderr, fmt.Sprintf("\x1b[%dmFATAL ", colorFatal), log.Ldate|log.Ltime|log.Lmicroseconds),
	}
}

// logger is an interface for writing levelled logs to stackdriver stdout and stderr
// -> When running in GCP. logs are available from GCE Instances dropdown in stackdriver
type logger struct {
	Client            *logging.Client
	SvcId             string
	env               string
	gcpProjectId      string
	debug             bool
	stackdriverLogger *logging.Logger
	debugLogger       *log.Logger
	infoLogger        *log.Logger
	warnLogger        *log.Logger
	errorLogger       *log.Logger
	fatalLogger       *log.Logger
}

// ChunkStringToSize can be used on values where we need the entire content e.g. error logs from forked processes
func ChunkStringToSize(s string, chunkSize int) []string {
	if chunkSize <= 0 {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == chunkSize {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}
func ChunkString(s string) []string {
	return ChunkStringToSize(s, fmtStringMaxLength)
}
func ChunkStringReplacingNewlines(s string) []string {
	return ChunkString(strings.ReplaceAll(s, "\n", " |"))
}
func ChunkStringToSizeReplacingNewlines(s string, chunkSize int) []string {
	return ChunkStringToSize(strings.ReplaceAll(s, "\n", " |"), chunkSize)
}

func InfoOutputAsChunks(ctx context.Context, msg string, output string) {
	outputs := ChunkStringReplacingNewlines(output)
	Info(ctx, fmt.Sprintf("%s - will log output in chucks", msg), FmtInt("len(outputs)", len(outputs)))
	for i, v := range outputs {
		Info(ctx, fmt.Sprintf("%s - chunked output", msg), FmtString(fmt.Sprintf("output[%d]", i), v))
	}
}

// TruncateField can be used on fields like description, notes, HTML/text email content where we only need a summary in the log
func TruncateField(field string) string {
	truncateTo := 42
	if len(field) < truncateTo {
		return field
	}
	return field[:truncateTo]
}

// TruncateFieldLenient can be used on fields where we only need a summary in the log but that summary should be as long as possible e.g. output from a child process 'exec.CommandContext'
func TruncateFieldLenient(field string) string {
	truncateTo := 50176
	if len(field) < truncateTo {
		return field
	}
	return field[:truncateTo]
}

// Init initializes the logger based on environment preferences
func Init(ctx context.Context, env lib_env.Env) error {
	log.Println("Initializing logger")

	client, err := logging.NewClient(ctx, env.GcpProjectId)
	if err != nil {
		return lib_errors.Wrap(err, "Failed creating stackdriver client")
	}

	l.Client = client
	l.SvcId = env.SvcId
	l.env = env.Id
	l.gcpProjectId = env.GcpProjectId
	l.debug = env.Debug
	l.stackdriverLogger = l.Client.Logger(l.SvcId)

	isRunningInGae = env.Gae
	isRunningInCloudRun = env.CloudRun

	disablePrettyPrint = isRunningInGae || isRunningInCloudRun

	Info(ctx, "Starting in environment", FmtAny("env", env), FmtStrings("os.Environ()", os.Environ()))
	Info(ctx, "Initialized logger")
	return nil
}

// Debug outputs logs
// Function is variadic in order to support passing no field parameter however ONLY ONE field is supported
//   i.e. use an array, or map if you need multiple fields, a struct will be logged as JSON
// Inside GCP a stackdriver structured logger is used with level Debug
// Outside GCP a stdout logger is used with level DEBUG
func Debug(ctx context.Context, message string, fields ...Field) {
	if l.debug {
		doLog(ctx, logging.Debug, message, fields)
	}
}

// Info outputs logs
// Function is variadic in order to support passing no field parameter however ONLY ONE field is supported
//   i.e. use an array, or map if you need multiple fields, a struct will be logged as JSON
// Inside GCP a stackdriver structured logger is used with level Info
// Outside GCP a stdout logger is used with level INFO
func Info(ctx context.Context, message string, fields ...Field) {
	doLog(ctx, logging.Info, message, fields)
}

// Notice outputs logs
// Function is variadic in order to support passing no field parameter however ONLY ONE field is supported
//   i.e. use an array, or map if you need multiple fields, a struct will be logged as JSON
// Inside GCP a stackdriver structured logger is used with level Notice
// Outside GCP a stdout logger is used with level INFO
func Notice(ctx context.Context, message string, fields ...Field) {
	doLog(ctx, logging.Notice, message, fields)
}

// Warn outputs logs
// Function is variadic in order to support passing no field parameter however ONLY ONE field is supported
//   i.e. use an array, or map if you need multiple fields, a struct will be logged as JSON
// Inside GCP a stackdriver structured logger is used with level Warning
// Outside GCP a stdout logger is used with level WARN
func Warn(ctx context.Context, message string, fields ...Field) {
	doLog(ctx, logging.Warning, message, fields)
}

// Error outputs logs
// Function is variadic in order to support passing no field parameter however ONLY ONE field is supported
//   i.e. use an array, or map if you need multiple fields, a struct will be logged as JSON
// Inside GCP a stackdriver structured logger is used with level Error
// Outside GCP a stderr logger is used with level ERROR
func Error(ctx context.Context, message string, fields ...Field) {
	doLog(ctx, logging.Error, message, fields)
}

// Fatal outputs logs and terminates the process
// Function is variadic in order to support passing no field parameter however ONLY ONE field is supported
//   i.e. use an array, or map if you need multiple fields, a struct will be logged as JSON
// Inside GCP a stackdriver structured logger is used with level Emergency
// Outside GCP a stderr logger is used with level FATAL
func Fatal(ctx context.Context, message string, fields ...Field) {
	doLog(ctx, logging.Emergency, message, fields)
}

// Field is the log Field type used to ensure only formatted log fields are used for logging
type Field string

type Formatter interface {
	Fmt(m Marshal) ([]byte, error)
}

type Marshal func(v interface{}) ([]byte, error)

// FmtAny formats a struct into a Field
func FmtAny(name string, value interface{}) Field {
	return FmtAnyWithLargeContentRedactions(name, value, false)
}

func FmtAnyWithLargeContentRedactions(name string, value interface{}, withLargeContentRedactions bool) Field {
	if e, ok := value.(error); ok {
		return FmtError(e)
	}
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.Zero(reflect.TypeOf(value)) == reflect.ValueOf(value)) {
		return Field(fmt.Sprintf("%q: null", name))
	}

	valType := fmt.Sprintf("%T", value)
	var field Field
	if f, ok := value.(Formatter); ok {
		if disablePrettyPrint {
			blob, err := f.Fmt(json.Marshal)
			if err != nil {
				Warn(context.Background(), "NOT FORMATTABLE", FmtString("name", name), FmtString("valType", valType), FmtError(lib_errors.Wrap(err, "Failed formatting value"))) // NEVER FmtAny here to avoid recursion
				field = Field(fmt.Sprintf("%q:{\"type\":%q,\"value\":\"NOT FORMATTABLE\",\"error\":\"%+v\"}", name, valType, err))
			} else {
				field = Field(fmt.Sprintf("%q:{\"type\":%q,\"value\":%s}", name, valType, redactLargeContentIfNecessary(string(blob), withLargeContentRedactions)))
			}
		} else {
			blob, err := f.Fmt(func(v interface{}) ([]byte, error) {
				return json.MarshalIndent(value, "\t\t", "\t")
			})
			if err != nil {
				Warn(context.Background(), "NOT FORMATTABLE", FmtString("name", name), FmtString("valType", valType), FmtError(lib_errors.Wrap(err, "Failed formatting value"))) // NEVER FmtAny here to avoid recursion
				field = Field(fmt.Sprintf("%q: {\n\t\t\"type\": %q,\n\t\t\"value\": \"NOT FORMATTABLE\",\n\t\t\"error\": \"%+v\"\n\t}", name, valType, err))
			} else {
				field = Field(fmt.Sprintf("%q: {\n\t\t\"type\": %q,\n\t\t\"value\": %s\n\t}", name, valType, redactLargeContentIfNecessary(string(blob), withLargeContentRedactions)))
			}
		}
	} else if disablePrettyPrint {
		blob, err := json.Marshal(value)
		if err != nil {
			Warn(context.Background(), "NOT JSON MARSHALLABLE", FmtString("name", name), FmtString("valType", valType), FmtError(lib_errors.Wrap(err, "Failed json marshalling value"))) // NEVER FmtAny here to avoid recursion
			field = Field(fmt.Sprintf("%q:{\"type\":%q,\"value\":\"NOT JSON MARSHALLABLE\",\"error\":\"%+v\"}", name, valType, err))
		} else {
			field = Field(fmt.Sprintf("%q:{\"type\":%q,\"value\":%s}", name, valType, redactLargeContentIfNecessary(string(blob), withLargeContentRedactions)))
		}
	} else {
		blob, err := json.MarshalIndent(value, "\t\t", "\t")
		if err != nil {
			Warn(context.Background(), "NOT JSON MARSHALLABLE", FmtString("name", name), FmtString("valType", valType), FmtError(lib_errors.Wrap(err, "Failed json marshalling value"))) // NEVER FmtAny here to avoid recursion
			field = Field(fmt.Sprintf("%q: {\n\t\t\"type\": %q,\n\t\t\"value\": \"NOT JSON MARSHALLABLE\",\n\t\t\"error\": \"%+v\"\n\t}", name, valType, err))
		} else {
			field = Field(fmt.Sprintf("%q: {\n\t\t\"type\": %q,\n\t\t\"value\": %s\n\t}", name, valType, redactLargeContentIfNecessary(string(blob), withLargeContentRedactions)))
		}
	}

	return field
}

func redactLargeContentIfNecessary(content string, withLargeContentRedactions bool) string {
	if !withLargeContentRedactions {
		return content
	}
	if len(content) > 1024 {
		return fmt.Sprintf("{\"content\":%q,\"length\":%d,\"truncated\":true}", content[:1024], len(content))
	}
	return content
}

// FmtBool formats a bool into a Field
func FmtBool(name string, value bool) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%t", name, value))
	}
	return Field(fmt.Sprintf("%q: %t", name, value))
}

// FmtByte formats a byte into a Field
func FmtByte(name string, value byte) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%q", name, string(value)))
	}
	return Field(fmt.Sprintf("%q: %q", name, string(value)))
}

// FmtBytes formats a byte slice into a Field
func FmtBytes(name string, value []byte) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%q", name, string(value)))
	}
	return Field(fmt.Sprintf("%q: %q", name, string(value)))
}

// FmtDuration formats a string into a  field
func FmtDuration(name string, value time.Duration) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%q", name, value.String()))
	}
	return Field(fmt.Sprintf("%q: %q", name, value.String()))
}

// FmtError formats an error into a Field
func FmtError(value error) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:\"%+v\"", "error", value))
	}
	return Field(fmt.Sprintf("%q: \"%+v\"", "error", value))
}

// FmtFloat32 formats an float32 into a Field
func FmtFloat32(name string, value float32) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%6f", name, value))
	}
	return Field(fmt.Sprintf("%q: %6f", name, value))
}

// FmtFloat64 formats an float64 into a Field
func FmtFloat64(name string, value float64) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%s", name, strconv.FormatFloat(value, 'f', 6, 64)))
	}
	return Field(fmt.Sprintf("%q: %s", name, strconv.FormatFloat(value, 'f', 6, 64)))
}

// FmtInt formats an int into a Field
func FmtInt(name string, value int) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%s", name, strconv.Itoa(value)))
	}
	return Field(fmt.Sprintf("%q: %s", name, strconv.Itoa(value)))
}

// FmtInt64 formats an int64 into a Field
func FmtInt64(name string, value int64) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%s", name, strconv.FormatInt(value, 10)))
	}
	return Field(fmt.Sprintf("%q: %s", name, strconv.FormatInt(value, 10)))
}

// FmtString formats a string into a Field
func FmtString(name string, value string) Field {
	lenOfValue := len(value)
	if lenOfValue > fmtStringMaxLength {
		value = fmt.Sprintf("%s, TRUNCATED %d characters", value[0:fmtStringMaxLength], lenOfValue-fmtStringMaxLength)
	}
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%q", name, value))
	}
	return Field(fmt.Sprintf("%q: %q", name, value))
}

// FmtStringAndRedact redacts and formats a string into a Field
func FmtStringAndRedact(name string, value string) Field {
	return FmtString(name, redactString(value))
}

func redactString(value string) string {
	if len(value) == 0 {
		return "NULL"
	} else if len(value) == 1 {
		return "*"
	} else if len(value) == 2 {
		return fmt.Sprintf("%s*", string(value[:1]))
	}
	return fmt.Sprintf("%s...%s", string(value[:1]), string(value[len(value)-1:]))
}

// FmtStrings formats a string slice into a Field
func FmtStrings(name string, values []string) Field {
	f := make([]string, len(values))
	for i, v := range values {
		f[i] = fmt.Sprintf("%q", v)
	}
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:[%s]", name, strings.Join(f, ",")))
	}
	return Field(fmt.Sprintf("%q: [\n\t\t%s\n\t]", name, strings.Join(f, ",\n\t\t")))
}

// FmtTime formats a time into a Field
func FmtTime(name string, value time.Time) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%q", name, value.Format(time.RFC3339)))
	}
	return Field(fmt.Sprintf("%q: %q", name, value.Format(time.RFC3339)))
}

// FmtUint formats a uint into a Field
func FmtUint(name string, value uint) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%s", name, strconv.FormatUint(uint64(value), 10)))
	}
	return Field(fmt.Sprintf("%q: %s", name, strconv.FormatUint(uint64(value), 10)))
}

// FmtUint64 formats a uint64 into a Field
func FmtUint64(name string, value uint64) Field {
	if disablePrettyPrint {
		return Field(fmt.Sprintf("%q:%s", name, strconv.FormatUint(value, 10)))
	}
	return Field(fmt.Sprintf("%q: %s", name, strconv.FormatUint(value, 10)))
}

func doLog(ctx context.Context, severity logging.Severity, message string, fields []Field) {
	fileName, functionName, line := runtimeCaller(2)

	var formattedFields string
	if len(fields) > 0 {
		formattedFields, severity = formatFields(fields, severity)
	} else {
		formattedFields = ""
	}

	if isRunningInGae || isRunningInCloudRun {
		logToStackdriver(ctx, severity, message, formattedFields, functionName, fileName, strconv.Itoa(line))
	} else {
		logToStd(ctx, severity, message, formattedFields, functionName, fmt.Sprintf("line %d", line), fileName)
	}
}

func MaxPayload(size int) bool {
	return size > MaxPayloadSize()
}

func MaxPayloadSize() int {
	if isRunningInGae {
		return stackdriverFieldSizeLimitBytes

	} else if isRunningInCloudRun {
		return stackdriverCloudRunFieldSizeLimitBytes

	} else {
		return int(^uint(0) >> 1)
	}
}

func formatFields(fields []Field, severity logging.Severity) (formattedFields string, newSeverity logging.Severity) {
	prefix := "{"
	suffix := "}"
	separator := ","

	if !disablePrettyPrint {
		prefix = prefix + "\n\t"
		suffix = "\n" + suffix
		separator = separator + "\n\t"
	}

	var fs []string
	for _, field := range fields {
		fs = append(fs, string(field))
	}
	formattedFields += fmt.Sprintf("%s%s%s", prefix, strings.Join(fs, separator), suffix)

	// Cannot log large (100KB -> https://cloud.google.com/logging/quotas) payloads to Stackdriver
	// -> Locally we still log the fields but swap to an warn level
	if MaxPayload(len(formattedFields) - 1000) { // Reduce the threshold in order not to trigger the "logging too large" condition since there is always additional metadata included by stackdriver
		var envDeployed string
		if isRunningInGae {
			envDeployed = "in STACKDRIVER"
		} else if isRunningInCloudRun {
			envDeployed = "in CLOUD RUN"
		}
		formattedFields = fmt.Sprintf("FIELDS OMITTED DUE TO TOO LARGE FOR STACKDRIVER %s, %d BYTES > %d BYTES: %s", envDeployed, len(formattedFields), stackdriverCloudRunFieldSizeLimitBytes, formattedFields[:MaxPayloadSize()/2])

		if severity == logging.Default || severity == logging.Debug || severity == logging.Info || severity == logging.Notice {
			severity = logging.Warning
		}
	}

	newSeverity = severity

	return
}

func logToStdHelper(ctx context.Context, logger *log.Logger, message, formattedFields, functionName, line, fileName string) {
	logger.Println(fmtStdLog(ctx, message, functionName, line, fileName, formattedFields))
}

func logToStdHelperWithPanic(ctx context.Context, logger *log.Logger, message, formattedFields, functionName, line, fileName string) {
	logger.Panicln(fmtStdLog(ctx, message, functionName, line, fileName, formattedFields))
}

func fmtStdLog(ctx context.Context, message, functionName, line, fileName, formattedFields string) string {
	message = fmt.Sprintf("%s %s %s %s %s", message, lib_context.CorrelationId(ctx), functionName, line, fileName)
	if formattedFields != "" {
		message = fmt.Sprintf("%s %s", message, formattedFields)
	}
	if lib_context.IntegrationTest(ctx) {
		message = fmt.Sprintf("INTEGRATION_TEST_REQUEST: %s", message)
	} else if lib_context.Test(ctx) {
		message = fmt.Sprintf("TEST_REQUEST: %s", message)
	}
	return fmt.Sprintf("%s\x1b[0m", message)
}

func logToStd(ctx context.Context, severity logging.Severity, message, formattedFields, functionName, line, fileName string) {
	var logger *log.Logger
	switch severity {
	case logging.Debug:
		logger = l.debugLogger
	case logging.Info, logging.Notice:
		logger = l.infoLogger
	case logging.Warning:
		logger = l.warnLogger
	case logging.Error:
		logger = l.errorLogger
	case logging.Emergency:
		logger = l.fatalLogger
		logToStdHelperWithPanic(ctx, logger, message, formattedFields, functionName, line, fileName)
		return
	default:
		logger = l.errorLogger
		message = "MISSING LOG LEVEL, USING ERROR => " + message
	}
	logToStdHelper(ctx, logger, message, formattedFields, functionName, line, fileName)
}

func logToStackdriver(ctx context.Context, severity logging.Severity, message, formattedFields, functionName, fileName, lineNumber string) {
	payload := make(map[string]string)
	correlationId := lib_context.CorrelationId(ctx)
	message = fmtMessage(ctx, correlationId, fileName, functionName, lineNumber, message)
	payload = map[string]string{
		"correlation_id": correlationId,
		"file_name":      fileName,
		"function_name":  functionName,
		"line_number":    lineNumber,
		"message":        message,
		"severity":       strings.ToUpper(severity.String()),
		"svc_id":         l.SvcId,
		"value":          formattedFields,
	}

	if isRunningInCloudRun {
		b, err := json.Marshal(payload)
		if err != nil {
			log.Println("Error preparing Cloud Run log entry", payload, err)
		} else {
			log.Println(string(b))
		}

	} else {
		if severity == logging.Emergency || severity == logging.Notice {
			entry := logging.Entry{Severity: severity, Payload: payload}
			err := l.stackdriverLogger.LogSync(context.Background(), entry)
			if err != nil {
				fileName, functionName, line := runtimeCaller(0)
				lineNumber := strconv.Itoa(line)
				formattedFields, severity := formatFields([]Field{FmtError(err)}, logging.Error)
				l.stackdriverLogger.Log(logging.Entry{Severity: logging.Error, Payload: map[string]string{
					"correlation_id": correlationId,
					"file_name":      fileName,
					"function_name":  functionName,
					"line_number":    lineNumber,
					"message":        fmtMessage(ctx, correlationId, fileName, functionName, lineNumber, "Failed log syncing with stackdriver logger"),
					"severity":       strings.ToUpper(severity.String()),
					"svc_id":         l.SvcId,
					"value":          formattedFields,
				}})
				l.stackdriverLogger.Log(entry)
			}

		} else {
			l.stackdriverLogger.Log(logging.Entry{Severity: severity, Payload: payload})
		}
	}

	// Also log any 'Emergency' AKA fatal errors to the stdout in order to carry out the PANIC
	if severity == logging.Emergency {
		logToStdHelperWithPanic(ctx, l.fatalLogger, message, formattedFields, functionName, fmt.Sprintf("line %s", lineNumber), fileName)
	}
}

func runtimeCaller(skip int) (fileName, functionName string, line int) {
	pc, fileName, line, _ := runtime.Caller(skip + 1)
	functionName = runtime.FuncForPC(pc).Name()
	return
}

func fmtMessage(ctx context.Context, correlationId, fileName, functionName, lineNumber, message string) string {
	message = fmt.Sprintf("%s # %s # %s # %s # %s:%s", l.SvcId, strings.Split(correlationId, ",")[0], message, functionName, fileName, lineNumber)
	if lib_context.IntegrationTest(ctx) {
		message = fmt.Sprintf("INTEGRATION_TEST_REQUEST: %s", message)
	} else if lib_context.Test(ctx) {
		message = fmt.Sprintf("TEST_REQUEST: %s", message)
	}
	return message
}
