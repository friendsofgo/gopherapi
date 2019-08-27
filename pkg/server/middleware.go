package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
)

type handler struct {
	serverID string
	next     http.Handler
}

func newServerMiddleware(serverID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := &handler{
			serverID: serverID,
			next:     next,
		}
		return h
	}
}

// ServeHTTP implements http.Handler.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := h.createRequestContext(r)
	h.next.ServeHTTP(w, r.WithContext(ctx))
}

func (h handler) createRequestContext(req *http.Request) context.Context {
	ctx := req.Context()

	var (
		xForwardedFor   = req.Header.Get("X-FORWARDED-FOR")
		xForwardedProto = req.Header.Get("X-FORWARDED-PROTO")
	)

	if xForwardedFor != "" {
		ctx = context.WithValue(ctx, contextKeyXForwardedFor, xForwardedFor)
	}
	if xForwardedProto != "" {
		ctx = context.WithValue(ctx, contextKeyXForwardedProto, xForwardedProto)
	}

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	ctx = context.WithValue(ctx, contextKeyClientIP, ip)
	ctx = context.WithValue(ctx, contextKeyEndpoint, req.URL.RequestURI())

	ctx = context.WithValue(ctx, contextKeyServerID, h.serverID)
	zipkinSpanHttpName(ctx, req)
	return ctx
}

func zipkinSpanHttpName(ctx context.Context, req *http.Request) {
	if span := zipkin.SpanFromContext(ctx); span != nil {
		if currentRoute := mux.CurrentRoute(req); currentRoute != nil {
			if routePath, err := currentRoute.GetPathTemplate(); err == nil {
				zipkin.TagHTTPRoute.Set(span, routePath)
				span.SetName(fmt.Sprintf("%s %s", req.Method, routePath))

			}
		}
	}
}
