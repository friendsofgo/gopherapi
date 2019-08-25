package logrus

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/friendsofgo/gopherapi/pkg/log"
	"github.com/friendsofgo/gopherapi/pkg/server"
)

// Event stores messages to log later, from our standard interface
type Event struct {
	id      string
	message string
}

// Logger centralize log messages format
type logger struct {
	*logrus.Logger
}

// NewLogger initializes the standard logger
func NewLogger(hooks ...logrus.Hook) log.Logger {
	return &logger{logrus.New()}
}

// Declare variables to store log messages as new Events
var (
	unexpectedErrorMessage = Event{"01DK2XFX9PQ85ZPZ5CP68P108Y", "Unexpected error: %v"}
)

func (l *logger) UnexpectedError(ctx context.Context, err error) {
	l.WithDefaultFields(ctx).WithField("logid", unexpectedErrorMessage.id).
		Errorf(unexpectedErrorMessage.message, err)
}

func (l *logger) WithDefaultFields(ctx context.Context) *logrus.Entry {
	hostname, _ := server.Name(ctx)
	httpAddr, _ := server.HttpAddr(ctx)

	return l.WithFields(logrus.Fields{
		"hostname": hostname,
		"httpAddr": httpAddr,
	})

}
