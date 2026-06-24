package sentry

import (
	"time"

	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
)

var sentryMiddleware = sentryfiber.New(sentryfiber.Options{
	Repanic:         true,
	WaitForDelivery: false,
	Timeout:         5 * time.Second,
})

func NewSentryMiddleware(c *fiber.Ctx) error {
	return sentryMiddleware(c)
}
