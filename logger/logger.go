package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wraps zap.SugaredLogger with contextual helpers
type Logger struct {
	*zap.SugaredLogger
}

// ReqLogInfo logs request information
func (l *Logger) ReqLogInfo(msg string, fields ...zap.Field) {
	l.Desugar().Info(msg, fields...)
}

// ReqLogError logs request errors
func (l *Logger) ReqLogError(msg string, fields ...zap.Field) {
	l.Desugar().Error(msg, fields...)
}

// ReqLogWarn logs request warnings
func (l *Logger) ReqLogWarn(msg string, fields ...zap.Field) {
	l.Desugar().Warn(msg, fields...)
}

// LogError logs an error with context
func (l *Logger) LogError(context string, userID string, message string) {
	fields := []zap.Field{
		zap.String("type", "user_action"),
		zap.String("context", context),
		zap.String("user_id", userID),
		zap.String("data", message),
	}
	l.Desugar().Error("api_service", fields...)
}

// LogInfo logs information with context
func (l *Logger) LogInfo(context string, userID string, message string) {
	fields := []zap.Field{
		zap.String("type", "user_action"),
		zap.String("context", context),
		zap.String("user_id", userID),
		zap.String("data", message),
	}
	l.Desugar().Info("api_service", fields...)
}

// LogBusinessInfo logs clinician business actions
// Example: LogBusinessInfo("68a9e2fe2238315d674133bd", "viewed patient profile with ID 68abb8b22df564bcd22f8db0")
func (l *Logger) LogBusinessInfo(clinicianID string, message string) {
	formattedMessage := fmt.Sprintf("(business)-%s %s", clinicianID, message)
	l.Infow(formattedMessage)
}

// NewLogger creates a new structured logger based on configs
func NewLogger(cfg *Config) (*Logger, error) {
	var zapCfg zap.Config
	zapCfg = zap.NewProductionConfig()
	encoder := zapcore.NewJSONEncoder(zapCfg.EncoderConfig)
	
	// stdout (for ECS / CloudWatch)
	stdout := zapcore.AddSync(os.Stderr)

	// rotating file writer
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   os.Getenv("APP_LOGS_PATH") + "/api.log",
		MaxSize:    200, // MB
		MaxBackups: 50,
		MaxAge:     7, // keep logs for 7 days
		Compress:   true,
	})

	allLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	// Core for all logs (Info, Warn, Error)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, stdout, allLevelEnabler),
		zapcore.NewCore(encoder, fileWriter, allLevelEnabler),
	)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.WarnLevel),
	)

	return &Logger{SugaredLogger: logger.Sugar()}, nil
}
