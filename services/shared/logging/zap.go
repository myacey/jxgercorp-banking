package logging

import "go.uber.org/zap"

func ConfigureLogger() (*zap.SugaredLogger, error) {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	return sugar, nil
}
