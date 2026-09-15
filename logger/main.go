package logger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	sentrygo "github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Logger interface {
	LogRequest(ctx *fiber.Ctx, fields RequestLogFields)
	LogErrorDeprecated(ctx *fiber.Ctx, errorContext string, errorMessage string)
	LogInfoDeprecated(ctx *fiber.Ctx, infoContext string, infoMessage string)
	ServiceName() string
	Logger(ctx *fiber.Ctx) *slog.Logger
}

type client struct {
	serviceName string
	logger      *slog.Logger
}

func Init(lc fx.Lifecycle) (Logger, error) {
	cfg := loadSentryConfig()
	if cfg == nil || cfg.DSN == "" || !strings.HasPrefix(cfg.DSN, "https://") {
		log.Println("logger disabled: DSN is empty or not configured")
		return nil, errors.New("logger disabled: DSN is empty or not configured")
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
		return nil, err
	}

	ctx := context.Background()

	consoleHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Set minimum level for console printing
	})
	sentryHandler := sentryslog.Option{
		LogLevel:  []slog.Level{slog.LevelInfo, slog.LevelWarn, slog.LevelError},
		AddSource: true,
	}.NewSentryHandler(ctx)
	multiHandler := slog.NewMultiHandler(sentryHandler, consoleHandler)
	slogLogger := slog.New(multiHandler)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			sentrygo.Flush(2 * time.Second)
			return nil
		},
	})
	return &client{
		serviceName: cfg.ServiceName,
		logger:      slogLogger,
	}, nil
}

func (c *client) ServiceName() string {
	return c.serviceName
}

func (c *client) Logger(ctx *fiber.Ctx) *slog.Logger {
	if ctx != nil {
		if reqLog, ok := ctx.Locals("request_logger").(*slog.Logger); ok {
			return reqLog
		}
	}
	return c.logger
}

func (c *client) LogRequest(ctx *fiber.Ctx, fields RequestLogFields) {
	msg := fmt.Sprintf("HTTP %s %s - %d", fields.Method, fields.Endpoint, fields.StatusCode)
	attrs := StructToAttrs(fields)
	//attrs := []any{
	//	slog.String("context", fields.Context),
	//	slog.String("event", fields.Event),
	//	slog.String("type", string(LogType(RequestLogType))),
	//	slog.String("service_name", c.serviceName),
	//	slog.String("timestamp", fields.Timestamp),
	//	slog.String("endpoint", fields.Endpoint),
	//	slog.String("method", fields.Method),
	//	slog.Int("status_code", fields.StatusCode),
	//	slog.String("request_id", fields.RequestID),
	//	slog.String("url", fields.URL),
	//	slog.String("client_ip", fields.ClientIP),
	//	slog.String("user_agent", fields.UserAgent),
	//	slog.String("user_id", fields.UserID),
	//	slog.String("organization_id", fields.OrganizationID),
	//	slog.Duration("latency", fields.Latency),
	//}

	if fields.ErrorMessage != "" {
		attrs = append(attrs, slog.String("error_message", fields.ErrorMessage))
		c.logger.ErrorContext(ctx.Context(), msg, attrs...)
		return
	}

	c.logger.InfoContext(ctx.Context(), msg, attrs...)
}
func (c *client) LogErrorDeprecated(ctx *fiber.Ctx, errorContext string, errorMessage string) {
	logInfo := c.Logger(ctx).
		With(slog.String("context", errorContext)).
		With(slog.String("log_type", string(LogType(ApplicationLogType))))
	if ctx == nil {
		logInfo.Error(errorMessage)
	} else {
		logInfo.ErrorContext(ctx.Context(), errorMessage)
	}
}
func (c *client) LogInfoDeprecated(ctx *fiber.Ctx, infoContext string, message string) {
	logInfo := c.logger.With(slog.String("context", infoContext)).With(slog.String("log_type", string(LogType(ApplicationLogType))))
	if ctx == nil {
		logInfo.Info(message)
	} else {
		logInfo.InfoContext(ctx.Context(), message)
	}
}
