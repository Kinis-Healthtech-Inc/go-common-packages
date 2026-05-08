package request_logging

type Config struct {
	ServiceName string
	GetUserFunc GetUserFunc
	Logger      Logger
	SkipPaths   []string
}
