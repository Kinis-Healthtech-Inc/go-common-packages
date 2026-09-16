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
			sentrygo.Flush(3 * time.Second)
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

	// Automatically maps struct fields & JSON tags to slog attributes
	attrs := StructToAttrs(fields)

	if fields.ErrorMessage != "" {
		c.logger.ErrorContext(ctx.Context(), msg, attrs...)
		return
	}

	c.logger.InfoContext(ctx.Context(), msg, attrs...)
}

func (c *client) LogErrorDeprecated(ctx *fiber.Ctx, errorContext string, errorMessage string) {
	logInfo := c.Logger(ctx).
		With(slog.String(LogKeyContext, errorContext)).
		With(slog.String(LogKeyLogType, string(LogType(ApplicationLogType))))

	if attrs := FileSourceAttributes(2); attrs != nil {
		logInfo = logInfo.With(attrs...)
	}

	if ctx == nil {
		logInfo.Error(errorMessage)
	} else {
		logInfo.ErrorContext(ctx.Context(), errorMessage)
	}
}

func (c *client) LogInfoDeprecated(ctx *fiber.Ctx, infoContext string, message string) {
	logInfo := c.Logger(ctx).
		With(slog.String(LogKeyContext, infoContext)).
		With(slog.String(LogKeyLogType, string(LogType(ApplicationLogType))))

	if attrs := FileSourceAttributes(2); attrs != nil {
		logInfo = logInfo.With(attrs...)
	}

	if ctx == nil {
		logInfo.Info(message)
	} else {
		logInfo.InfoContext(ctx.Context(), message)
	}
}
