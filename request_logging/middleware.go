package request_logging

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// New creates a new request logger middleware with custom configuration
func New(cfg Config) fiber.Handler {
	if cfg.ServiceName == "" {
		cfg.ServiceName = "service"
	}
	if cfg.SkipPaths == nil {
		cfg.SkipPaths = []string{"/health"}
	}
	if cfg.Logger == nil {
		panic("Logger must be set on the cfg")
	}

	return func(c *fiber.Ctx) error {
		url := c.OriginalURL()

		// Check if path should be skipped
		if shouldSkipPath(url, cfg.SkipPaths) {
			return c.Next()
		}

		start := time.Now()
		method := c.Method()
		clientIP := c.IP()
		userAgent := c.Get("User-Agent")
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
			c.Set("X-Request-ID", requestID)
		}
		queryParams := c.Queries()

		// execute handler
		err := c.Next()

		statusCode := c.Response().StatusCode()
		// If no route was matched, set 404 status
		if c.Response().StatusCode() == fiber.StatusOK && err != nil {
			if e, ok := err.(*fiber.Error); ok {
				if e.Code == fiber.StatusNotFound {
					statusCode = fiber.StatusNotFound
				}
			}
		}
		latency := time.Since(start)
		var userID string
		var organizationID string

		// Extract user context if GetUserFunc is provided
		if cfg.GetUserFunc != nil {
			authCtx, authErr := cfg.GetUserFunc(c)
			if authErr == nil {
				userID = authCtx.UserID
				organizationID = authCtx.OrganizationID
			}
		}

		fields := []zap.Field{
			zap.String("timestamp", time.Now().UTC().Format(time.RFC3339)),
			zap.String("endpoint", c.Path()),
			zap.String("method", method),
			zap.Int("status_code", statusCode),
			zap.Any("query_params", queryParams),
			zap.String("request_id", requestID),
			zap.String("url", url),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.String("user_id", userID),
			zap.String("organization_id", organizationID),
			zap.Duration("latency", latency),
		}

		if err != nil {
			fields = append(fields, zap.String("error_message", err.Error()))
		}

		// Log based on status code or error
		switch {
		case statusCode >= 500:
			cfg.Logger.ReqLogError(cfg.ServiceName, fields...)
		case statusCode >= 400:
			cfg.Logger.ReqLogWarn(cfg.ServiceName, fields...)
		default:
			cfg.Logger.ReqLogInfo(cfg.ServiceName, fields...)
		}

		return err
	}
}

// shouldSkipPath checks if the URL matches any skip path pattern
func shouldSkipPath(url string, skipPaths []string) bool {
	for _, path := range skipPaths {
		if strings.Contains(url, path) {
			return true
		}
	}
	return false
}
