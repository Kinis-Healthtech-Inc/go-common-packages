package sentry

import (
	"time"

	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
)

// Handler is a named type around fiber.Handler so the sentryfiber middleware
// can be injected through fx by type without colliding with other unrelated
// fiber.Handler providers in the DI graph.
type Handler fiber.Handler

// NewSentryMiddleware returns the sentryfiber middleware that attaches a
// per-request *sentry.Hub to the Fiber context and auto-captures transactions
// and panics.
//
// Per the official Sentry-Fiber docs:
//   - Repanic must be true so Fiber's Recover middleware (registered before
//     this one) can finalize the HTTP response after Sentry records the panic.
//   - WaitForDelivery is false because Recover does not restart the app; we do
//     not need to block the response on event delivery.
//   - Timeout caps how long delivery may take before the request continues.
//
// IMPORTANT: This middleware must be registered AFTER recover.New() and BEFORE
// any middleware that wants to enrich the per-request hub (auth, logger, etc.),
// otherwise sentryfiber.GetHubFromContext(c) will return nil for those layers.
func NewSentryMiddleware() Handler {
	return Handler(sentryfiber.New(sentryfiber.Options{
		Repanic:         true,
		WaitForDelivery: false,
		Timeout:         5 * time.Second,
	}))
}
