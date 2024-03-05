package featureflag

import (
	"context"
	"log/slog"
)

// Logger is the interface used to log errors, it is a subset of slog.Logger.
type Logger interface {
	ErrorContext(ctx context.Context, msg string, args ...any)
}

var logger Logger = slog.Default()

// SetLogger sets the logger used to log errors.
func SetLogger(l Logger) {
	logger = l
}
