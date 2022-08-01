package xlogger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey struct {
	key string
}

func (ctx contextKey) String() string {
	return "xlogger: " + ctx.key
}

var (
	loggerKey = &contextKey{"logger"}
)

func contextLogger(ctx context.Context) *logrus.Entry {
	log, _ := ctx.Value(loggerKey).(*logrus.Entry)

	return log
}

// Logger get the logger entry from context
func Logger(ctx context.Context) *logrus.Entry {
	if entry := contextLogger(ctx); entry != nil {
		return entry
	}

	return nil
}

// SetLogger save the given log entry into the context if it's not there yet
func SetLogger(ctx context.Context, log *logrus.Entry) context.Context {
	entry := contextLogger(ctx)
	if entry == nil {
		return context.WithValue(ctx, loggerKey, log)
	}

	return ctx
}
