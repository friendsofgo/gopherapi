package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Event stores messages to log later, from our standard interface
type Event struct {
	id      string
	message string
}

// Logger centralize log messages format
type Logger struct {
	*logrus.Logger
}

// NewLogger initializes the standard logger
func NewLogger(hooks ...logrus.Hook) *Logger {
	logger := logrus.New()

	for _, hook := range hooks {
		logger.AddHook(hook)
	}

	return &Logger{logger}
}

// Declare variables to store log messages as new Events
var (
	unexpectedErrorMessage = Event{"01DK2XFX9PQ85ZPZ5CP68P108Y", "Unexpected error: %v"}
)

// ElementNotFound is a standard error message for elements not found
func (l *Logger) UnexpectedError(ctx context.Context, err error) {
	l.WithContext(ctx).WithField("logid", unexpectedErrorMessage.id).
		Errorf(unexpectedErrorMessage.message, err)
}
