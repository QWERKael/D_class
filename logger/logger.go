package logger

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

var Logger = CreateLogger()

func CreateLogger() *logrus.Logger {
	logger := &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.InfoLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}
	return logger
}

func LogLevel(logLevel string) {
	var lvl logrus.Level
	switch logLevel {
	case "debug":
		lvl = logrus.DebugLevel
	case "info":
		lvl = logrus.InfoLevel
	case "warn":
		lvl = logrus.WarnLevel
	case "error":
		lvl = logrus.ErrorLevel
	case "fatal":
		lvl = logrus.FatalLevel
	case "panic":
		lvl = logrus.PanicLevel
	default:
		lvl = logrus.InfoLevel
	}
	Logger.Level = lvl
}
