package sentry

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	sentrygo "github.com/getsentry/sentry-go"
	"go.uber.org/fx"
)

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

type SentryConfig struct {
	DSN              string
	Environment      string
	Release          string
	Debug            bool
	SampleRate       float64
	TracesSampleRate float64
	EnableTracing    bool
	MaxBreadcrumbs   int
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
	}
}

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

func Init(lc fx.Lifecycle) error {
	cfg := loadSentryConfig()
	if cfg == nil || cfg.DSN == "" || !strings.HasPrefix(cfg.DSN, "https://") {
		log.Println("sentry disabled: DSN is empty or not configured")
		return nil
	}

	opts := sentrygo.ClientOptions{
		Dsn:              cfg.DSN,
		Environment:      cfg.Environment,
		Release:          cfg.Release,
		Debug:            cfg.Debug,
		SendDefaultPII:   false,
		EnableTracing:    cfg.EnableTracing,
		TracesSampleRate: cfg.TracesSampleRate,
		BeforeSend:       beforeSend,
	}
	if cfg.SampleRate > 0 {
		opts.SampleRate = cfg.SampleRate
	}
	if cfg.MaxBreadcrumbs > 0 {
		opts.MaxBreadcrumbs = cfg.MaxBreadcrumbs
	}

	if err := sentrygo.Init(opts); err != nil {
		return err
	}
	log.Printf("sentry initialized (env=%s, release=%s)", cfg.Environment, cfg.Release)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			sentrygo.Flush(2 * time.Second)
			return nil
		},
	})
	return nil
}
