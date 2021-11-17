package logger

import "go.uber.org/zap"

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger() (Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout", "csgostate-server.log"}
	logger, err := config.Build(zap.AddCaller())
	if err != nil {
		return Logger{}, err
	}
	return Logger{logger.Sugar()}, err
}