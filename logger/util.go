package logger

import (
	"log/slog"
	"runtime"
)

// Here captures the exact file and line of the code line it is invoked on.
func Here() slog.Attr {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return slog.Group("source", slog.String("file", "unknown"))
	}
	return slog.Group("source",
		slog.String(LogKeySourceFile, file),
		slog.Int(LogKeySourceLine, line),
	)
}
