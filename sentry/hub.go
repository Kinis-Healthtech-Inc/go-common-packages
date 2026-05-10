package sentry

import (
	"time"

	sentrygo "github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
)

// hubFromCtx returns the per-request *sentry.Hub attached by the sentryfiber
// middleware. It returns nil if Sentry is disabled, the middleware did not run
// for this request (e.g. background jobs), or the ctx is nil.
func hubFromCtx(c *fiber.Ctx) *sentrygo.Hub {
	if c == nil {
		return nil
	}
	return sentryfiber.GetHubFromContext(c)
}

// AddBreadcrumb appends a breadcrumb to the request hub. Breadcrumbs form the
// "actions before error" timeline visible in the Sentry UI and are dropped
// silently when Sentry is disabled.
func AddBreadcrumb(c *fiber.Ctx, category, message string, level sentrygo.Level, data map[string]interface{}) {
	hub := hubFromCtx(c)
	if hub == nil {
		return
	}
	hub.AddBreadcrumb(&sentrygo.Breadcrumb{
		Category: category,
		Message:  message,
		Level:    level,
		Data:     data,
	}, nil)
}

// SetUser tags subsequent events from this request with user identity. Call
// this from auth middleware once the user has been resolved.
func SetUser(c *fiber.Ctx, id, email, role string) {
	hub := hubFromCtx(c)
	if hub == nil {
		return
	}
	hub.Scope().SetUser(sentrygo.User{
		ID:    id,
		Email: email,
		Data:  map[string]string{"role": role},
	})
}

// SetTag attaches a single indexed tag to subsequent events from this request.
func SetTag(c *fiber.Ctx, key, value string) {
	if hub := hubFromCtx(c); hub != nil {
		hub.Scope().SetTag(key, value)
	}
}

// SetContext attaches a structured (non-indexed) context block to subsequent
// events from this request. Use this for richer payloads that should not be
// indexed as tags.
func SetContext(c *fiber.Ctx, key string, value map[string]interface{}) {
	if hub := hubFromCtx(c); hub != nil {
		hub.Scope().SetContext(key, value)
	}
}

// CaptureError sends an error to Sentry using the per-request hub when
// available, or the global hub as a fallback (useful for background jobs and
// non-HTTP code paths). It is safe to call with a nil error.
func CaptureError(c *fiber.Ctx, err error) {
	if err == nil {
		return
	}
	if hub := hubFromCtx(c); hub != nil {
		hub.CaptureException(err)
		hub.Flush(2 * time.Second)
		return
	}
	sentrygo.CaptureException(err)
}

// CaptureMessage sends a non-error message to Sentry at the given level using
// the per-request hub when available, or the global hub as a fallback.
func CaptureMessage(c *fiber.Ctx, msg string, level sentrygo.Level) {
	if hub := hubFromCtx(c); hub != nil {
		hub.WithScope(func(scope *sentrygo.Scope) {
			scope.SetLevel(level)
			hub.CaptureMessage(msg)
		})
		return
	}
	sentrygo.WithScope(func(scope *sentrygo.Scope) {
		scope.SetLevel(level)
		sentrygo.CaptureMessage(msg)
	})
}
