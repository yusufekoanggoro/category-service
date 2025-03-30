package logger

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(message, event, key string)
	Warn(message, event, key string)
	Error(message, event, key string)
	Debug(message, event, key string)
	Fatal(message, event, key string)
	Panic(message, event, key string)
	SetOutput(output *os.File)
}

type LoggerImpl struct {
	logger  *logrus.Logger
	appName string
}

var (
	once     sync.Once
	instance *LoggerImpl
)

// NewLogger creates a singleton logger instance
func NewLogger(appName string, logLevel logrus.Level, output *os.File) Logger {
	once.Do(func() {
		logger := logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: false,
			ForceQuote:    true,
		})
		logger.SetLevel(logLevel)
		logger.SetOutput(output)

		instance = &LoggerImpl{logger: logger, appName: appName}
	})

	return instance
}

// SetOutput allows changing the output of the logger
func (l *LoggerImpl) SetOutput(output *os.File) {
	l.logger.SetOutput(output)
}

func (l *LoggerImpl) logWithFields(level logrus.Level, message, event, key string) {
	fields := logrus.Fields{
		"caller": l.getCallerInfo(),
		"topic":  l.appName,
		"event":  event,
		"key":    key,
	}

	entry := l.logger.WithFields(fields)
	switch level {
	case logrus.InfoLevel:
		entry.Info(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.ErrorLevel:
		entry.Error(message)
	case logrus.DebugLevel:
		entry.Debug(message)
	case logrus.FatalLevel:
		entry.Fatal(message)
	case logrus.PanicLevel:
		entry.Panic(message)
	}
}

func (l *LoggerImpl) getCallerInfo() string {
	if pc, file, line, ok := runtime.Caller(2); ok {
		return fmt.Sprintf("%s:%d %s", file, line, runtime.FuncForPC(pc).Name())
	}
	return "Unknown Caller"
}

func (l *LoggerImpl) Info(message, event, key string) {
	l.logWithFields(logrus.InfoLevel, message, event, key)
}

func (l *LoggerImpl) Warn(message, event, key string) {
	l.logWithFields(logrus.WarnLevel, message, event, key)
}

func (l *LoggerImpl) Error(message, event, key string) {
	l.logWithFields(logrus.ErrorLevel, message, event, key)
}

func (l *LoggerImpl) Debug(message, event, key string) {
	l.logWithFields(logrus.DebugLevel, message, event, key)
}

func (l *LoggerImpl) Fatal(message, event, key string) {
	l.logWithFields(logrus.FatalLevel, message, event, key)
}

func (l *LoggerImpl) Panic(message, event, key string) {
	l.logWithFields(logrus.PanicLevel, message, event, key)
}
