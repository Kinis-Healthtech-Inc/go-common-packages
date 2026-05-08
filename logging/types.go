package logging

import "github.com/gofiber/fiber/v2"

type ContextAuthKey struct {
	UserID           string
	Role             string
	OrganizationID   string
	OrganizationName string
}

type GetUserFunc func(c *fiber.Ctx) (*ContextAuthKey, error)
