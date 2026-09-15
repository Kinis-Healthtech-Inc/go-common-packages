package logger

import (
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

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

// StructToAttrs converts any struct (including embedded fields and json tags) into []any for slog
func StructToAttrs(obj any) []any {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	var attrs []any
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)

		// Skip unexported fields
		if fieldTyp.PkgPath != "" && !fieldTyp.Anonymous {
			continue
		}

		// Handle embedded structs recursively (e.g., LogFields, UserInfo)
		if fieldTyp.Anonymous && fieldVal.Kind() == reflect.Struct {
			attrs = append(attrs, StructToAttrs(fieldVal.Interface())...)
			continue
		}

		// Extract key from JSON tag, fallback to field name
		tag := fieldTyp.Tag.Get("json")
		key := strings.Split(tag, ",")[0]
		if key == "" || key == "-" {
			key = fieldTyp.Name
		}

		// Skip empty omitempty fields
		if strings.Contains(tag, "omitempty") && fieldVal.IsZero() {
			continue
		}

		// Map primitive and specialized types to correct slog typed Attrs
		switch v := fieldVal.Interface().(type) {
		case string:
			if v != "" {
				attrs = append(attrs, slog.String(key, v))
			}
		case int:
			attrs = append(attrs, slog.Int(key, v))
		case int64:
			attrs = append(attrs, slog.Int64(key, v))
		case float64:
			attrs = append(attrs, slog.Float64(key, v))
		case bool:
			attrs = append(attrs, slog.Bool(key, v))
		case time.Duration:
			attrs = append(attrs, slog.Duration(key, v))
		case LogType:
			attrs = append(attrs, slog.String(key, string(v)))
		default:
			if !fieldVal.IsZero() {
				attrs = append(attrs, slog.Any(key, fieldVal.Interface()))
			}
		}
	}
	return attrs
}
