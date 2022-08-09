package xlogger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type (
	FieldLogger = logrus.FieldLogger
)

type contextKey struct {
	key string
}

type Config struct {
	Formatter map[string]string `envconfig:"FORMATTER"`
	Level     string            `envconfig:"LEVEL"`
}

var (
	loggerKey = &contextKey{"logger"}
)

func (cfg *Config) setDefault() {
	if cfg.Level == "" {
		cfg.Level = "info"
	}
}

func New(cfg *Config) *logrus.Logger {
	if cfg == nil {
		cfg = new(Config)
	}
	cfg.setDefault()

	log := logrus.StandardLogger()
	log.SetOutput(os.Stdout)

	var fmtter logrus.Formatter
	if typ, ok := cfg.Formatter["type"]; ok && typ == "text" {
		fmtter = &logrus.TextFormatter{}
	} else {
		fmtter = &logrus.JSONFormatter{}
	}

	log.SetFormatter(fmtter)

	if lvl, err := logrus.ParseLevel(cfg.Level); err == nil {
		log.SetLevel(lvl)
	}

	return log
}

func (ctx contextKey) String() string {
	return "xlogger: " + ctx.key
}

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
