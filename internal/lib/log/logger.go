package log

import "go.uber.org/zap"

var Logger *zap.SugaredLogger = nil

func InitializeLogging() {
	productionLogger, _ := zap.NewProduction()
	productionLogger.Sync()
	Logger = productionLogger.Sugar()
}
