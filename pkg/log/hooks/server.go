package hooks

import (
	"github.com/sirupsen/logrus"

	"github.com/friendsofgo/gopherapi/pkg/server"
)

type serverHook struct{}

func NewServerInformationHook() logrus.Hook {
	return &serverHook{}
}

func (h *serverHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context

	hostname, _ := server.Name(ctx)
	httpAddr, _ := server.HttpAddr(ctx)

	entry.Data["hostname"] = hostname
	entry.Data["httpAddr"] = httpAddr

	return nil
}

func (h *serverHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
