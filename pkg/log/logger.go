package log

import (
	"context"
)

// Logger determine the way to centralize log messages format
type Logger interface {
	// UnexpectedError is a standard error message for unexpected errors
	UnexpectedError(ctx context.Context, err error)
}
