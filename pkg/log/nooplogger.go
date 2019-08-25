package log

import "context"

type noop struct {
}

// NewNoopLogger ...
func NewNoopLogger() Logger {
	return &noop{}
}

func (l *noop) UnexpectedError(ctx context.Context, err error) {
	// nothing to do here
}
