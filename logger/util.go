package logger

import (
	"fmt"
	"log/slog"
	"runtime"
)

func GetCallStack(depth int) slog.Attr {
	var frames []any
	skip := 3 // adjust based on your call stack depth

	for i := 0; i < depth; i++ {
		if _, file, line, ok := runtime.Caller(skip); ok {
			frames = append(frames, slog.Group(fmt.Sprintf("frame_%d", i),
				slog.String(LogKeySourceFile, file),
				slog.Int(LogKeySourceLine, line),
			))
			skip++
		} else {
			break
		}
	}

	return slog.Group("call_stack", frames...)
}
