package logger

import (
	"auth-service/pkg/config"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	onceLog  sync.Once
	instance *logrus.Logger
)

const (
	envLocal = "local"
	envDebug = "debug"
	envProd  = "prod"
)

func SetupLogger(cfg *config.Config) *logrus.Logger {
	onceLog.Do(func() {

		env := cfg.Level

		logger := logrus.New()
		logger.SetOutput(os.Stdout)

		switch env {
		case envLocal:
			logger.SetLevel(logrus.DebugLevel)
			logger.SetFormatter(&logrus.JSONFormatter{})
			logger.SetReportCaller(true)
		case envDebug:
			logger.SetLevel(logrus.TraceLevel)
			logger.SetFormatter(&logrus.JSONFormatter{})
			logger.SetReportCaller(true)
		case envProd:
			logger.SetLevel(logrus.InfoLevel)
			logger.SetFormatter(&logrus.TextFormatter{})
			logger.SetReportCaller(false)
		}

		instance = logger
	})

	return instance
}
