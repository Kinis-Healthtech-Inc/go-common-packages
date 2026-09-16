package logger

import (
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	sentrygo "github.com/getsentry/sentry-go"
)

func parseFloat(s string, fallback float64) float64 {
	if s == "" {
		return fallback
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fallback
	}
	return v
}

func parseInt(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return v
}

type SentryConfig struct {
	DSN              string
	Environment      string
	Release          string
	Debug            bool
	SampleRate       float64
	TracesSampleRate float64
	EnableTracing    bool
	MaxBreadcrumbs   int
	ServiceName      string
}

func loadSentryConfig() *SentryConfig {
	sampleRate := parseFloat(os.Getenv("SENTRY_SAMPLE_RATE"), 1.0)
	tracesSampleRate := parseFloat(os.Getenv("SENTRY_TRACES_SAMPLE_RATE"), 0.2)
	maxBreadcrumbs := parseInt(os.Getenv("SENTRY_MAX_BREADCRUMBS"), 100)
	sentryDns := os.Getenv("SENTRY_DSN")
	sentryEnvironment := os.Getenv("SENTRY_ENVIRONMENT")
	sentryDebug := os.Getenv("SENTRY_DEBUG") == "true"
	sentryEnableTraces := os.Getenv("SENTRY_ENABLE_TRACES") == "true"

	if sentryEnvironment == "" || sentryDns == "" || !strings.HasPrefix(sentryDns, "https://") {
		panic("SENTRY_DSN and SENTRY_ENVIRONMENT must be set and valid for Sentry to work")
	}

	return &SentryConfig{
		DSN:              sentryDns,
		Debug:            sentryDebug,
		Environment:      sentryEnvironment,
		SampleRate:       sampleRate,
		TracesSampleRate: tracesSampleRate,
		EnableTracing:    sentryEnableTraces,
		MaxBreadcrumbs:   maxBreadcrumbs,
		ServiceName:      os.Getenv("SENTRY_SERVICE_NAME"),
	}
}

// sensitiveHeaders lists request headers that must never reach Sentry.
// Lookups are done in lower-case.
var sensitiveHeaders = map[string]struct{}{
	"authorization":      {},
	"cookie":             {},
	"set-cookie":         {},
	"x-api-key":          {},
	"x-internal-api-key": {},
	"x-auth-token":       {},
}

// beforeSend is the SDK hook used to redact sensitive data from outgoing
// events before they are delivered to Sentry. It is wired into
// ClientOptions.BeforeSend in Init.
//
// Returning nil drops the event entirely. We currently only redact; add
// project-specific drop rules here (e.g. context.Canceled, expected 4xx
// errors) when needed.
func beforeSend(event *sentrygo.Event, _ *sentrygo.EventHint) *sentrygo.Event {
	if event == nil {
		return nil
	}
	if event.Request != nil {
		event.Request.Cookies = ""
		for k := range event.Request.Headers {
			if _, ok := sensitiveHeaders[strings.ToLower(k)]; ok {
				event.Request.Headers[k] = "[redacted]"
			}
		}
	}

	return event
}

func StructToAttrs(v any) []any {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var attrs []any
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		structField := typ.Field(i)

		// Handle embedded structs (like LogFields) recursively
		if structField.Anonymous {
			if fieldVal.Kind() == reflect.Struct {
				attrs = append(attrs, StructToAttrs(fieldVal.Interface())...)
			}
			continue
		}

		tag := structField.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		// Extract key name from json tag (e.g., "user_id,omitempty" -> "user_id")
		name := tag
		if idx := strings.Index(tag, ","); idx != -1 {
			name = tag[:idx]
		}

		// Skip zero values if "omitempty" is specified in the tag
		if fieldVal.IsZero() && strings.Contains(tag, "omitempty") {
			continue
		}

		attrs = append(attrs, slog.Any(name, fieldVal.Interface()))
	}

	return attrs
}
func FileSourceAttributes(skip int) []any {
	if _, file, line, ok := runtime.Caller(skip); ok {
		return []any{
			slog.String(LogKeySourceFile, file),
			slog.Int(LogKeySourceLine, line),
		}
	}
	return nil
}
