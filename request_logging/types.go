package request_logging

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
	ReqLogError(msg string, fields ...zap.Field)
	ReqLogInfo(msg string, fields ...zap.Field)
	ReqLogWarn(msg string, fields ...zap.Field)
}
