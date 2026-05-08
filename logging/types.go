package logging

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ContextAuthKey struct {
	UserID           string
	Role             string
	OrganizationID   string
	OrganizationName string
}

type GetUserFunc func(c *fiber.Ctx) (*ContextAuthKey, error)

type Logger interface {
	Error(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
}
