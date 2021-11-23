package logger

import "go.uber.org/zap"

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(filename string) (Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", filename}
	logger, err := config.Build(zap.AddCaller())
	if err != nil {
		return Logger{}, err
	}
	return Logger{logger.Sugar()}, err
}
