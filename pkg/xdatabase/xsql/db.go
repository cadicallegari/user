package xsql

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type (
	DB = sqlx.DB
)

func Open(driverName string, cfg *Config) (*DB, error) {
	if cfg == nil {
		cfg = new(Config)
	}
	cfg.setDefaults()

	db, err := sqlx.Open(driverName, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("%s_connection_error unable to connect to %s: %w", driverName, driverName, err)
	}

	// set config values
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	if err := db.Ping(); err != nil {
		// return db to caller be able to try to ping
		// later, e.g. using WaitOK
		return db, err
	}

	return db, nil
}

func WaitOK(db *DB, errfn func(_ error, attempt int) bool) error {
	var attempt int
	for {
		err := db.Ping()
		if err == nil {
			break
		}

		attempt++

		if errfn != nil {
			cont := errfn(err, attempt)
			if !cont {
				return err
			}
		}

		time.Sleep(time.Second * time.Duration(math.Pow(2, float64(attempt))))
	}

	return nil
}

var ErrMigrationDisabled = errors.New("migrations is disabled")

func Migration(driverName string, db *DB, cfg *Config) (version uint, dirty bool, err error) {
	if !cfg.RunMigration || cfg.MigrationsDir == "" || cfg.MigrationDriverFn == nil {
		return 0, false, ErrMigrationDisabled
	}

	sourceDriver, err := source.Open(cfg.MigrationsDir)
	if err != nil {
		return 0, false, fmt.Errorf("%s_migration unable to create source driver %s: %w", driverName, driverName, err)
	}
	defer sourceDriver.Close()

	migrationDriver, err := cfg.MigrationDriverFn(db)
	if err != nil {
		return 0, false, fmt.Errorf("%s_migration unable to create migration driver: %w", driverName, err)
	}

	m, err := migrate.NewWithInstance("file", sourceDriver, driverName, migrationDriver)
	if err != nil {
		return 0, false, fmt.Errorf("%s_migration unable to create migration instance: %w", driverName, err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return 0, false, err
	}

	return m.Version()
}
