package logger

import (
	"os"

	"github.com/gaurav2721/notification-service/constants"
	"github.com/sirupsen/logrus"
)

// Configure sets up the logging configuration based on environment variables
func Configure() {
	// Configure logrus formatter
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Set log level from environment variable or default to InfoLevel
	logLevel := os.Getenv(constants.LOG_LEVEL)
	switch logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel) // Default to InfoLevel to disable debug logs
	}

	logrus.Debug("Logger configured successfully")
}

// GetLogger returns the configured logrus logger instance
func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

// SetLevel allows runtime log level changes
func SetLevel(level string) {
	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	}
}
