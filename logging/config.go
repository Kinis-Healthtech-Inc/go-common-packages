package logging

import "go.uber.org/zap"

type Config struct {
	ServiceName string
	GetUserFunc GetUserFunc
	Logger      *zap.Logger
	SkipPaths   []string
}
