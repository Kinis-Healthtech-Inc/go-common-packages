package logger

import (
	"context"
	"log/slog"
)

type contextKey string

const userContextKey contextKey = "slog_user_info"

// InjectUserIntoContext stores user info in context
func InjectUserIntoContext(ctx context.Context, info UserInfo) context.Context {
	return context.WithValue(ctx, userContextKey, info)
}

// ContextHandler intercepts logs and auto-injects context-bound user info
type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if info, ok := ctx.Value(userContextKey).(UserInfo); ok {
			r.AddAttrs(
				slog.String("user_id", info.UserID),
				slog.String("role", info.Role),
				slog.String("organization_id", info.OrganizationID),
			)
		}
	}
	return h.Handler.Handle(ctx, r)
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{Handler: h.Handler.WithGroup(name)}
}
