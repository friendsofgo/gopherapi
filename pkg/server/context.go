package server

import (
	"context"
)

var (
	contextKeyServerName = contextKey("Name")
	contextKeyHttpAddr   = contextKey("HttpAddr")
)

type contextKey string

func (c contextKey) String() string {
	return "server" + string(c)
}

// Name gets the name server from context
func Name(ctx context.Context) (string, bool) {
	n, ok := ctx.Value(contextKeyServerName).(string)
	return n, ok
}

// HttpAddr gets the http address server from context
func HttpAddr(ctx context.Context) (string, bool) {
	httpAddrStr, ok := ctx.Value(contextKeyHttpAddr).(string)
	return httpAddrStr, ok
}
