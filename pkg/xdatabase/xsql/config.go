package xsql

import (
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4/database"

	"github.com/cadicallegari/user/pkg/xlogger"
)

type Config struct {
	// URL do not put the database, use the field .Database otherwise.
	URL string `envconfig:"URL"`

	Charset string

	// ConnMaxLifetime use < 0 to reuse connections forever.
	ConnMaxLifetime time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"1h"`

	// MaxIdleConns use < 0 for no idle connections.
	MaxIdleConns int `envconfig:"MAX_IDLE_CONNS" default:"1"`

	// MaxOpenConns use < 0 for no limit on the number of open connections.
	MaxOpenConns int `envconfig:"MAX_OPEN_CONNS" default:"10"`

	RunMigration    bool   `envconfig:"RUN_MIGRATIONS" default:"true"`
	MigrationsDir   string `envconfig:"MIGRATIONS_DIR"`
	MigrationsTable string `envconfig:"MIGRATIONS_TABLE"`

	MigrationDriverFn func(*DB) (database.Driver, error)
	Logger            xlogger.FieldLogger
}

func (cfg *Config) setDefaults() {
	if cfg.Charset == "" {
		cfg.Charset = "utf8mb4_unicode_ci"
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = 1 * time.Hour
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 1
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 10
	}
	if cfg.MigrationsDir != "" && !strings.HasPrefix(cfg.MigrationsDir, "file://") {
		cfg.MigrationsDir = "file://" + cfg.MigrationsDir
	}
}
