package durable

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Custom logger
type Logger struct {
	*logrus.Logger
}

// NewLogger create a new logger
func NewLogger(logger *logrus.Logger) *Logger {
	file, err := os.OpenFile("bitlygo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Warnf("Can not create log file: %v", err)
	} else {
		logger.Out = file
	}

	return &Logger{logger}
}
