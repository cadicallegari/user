package xmysql

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4/database"
	migrationmysql "github.com/golang-migrate/migrate/v4/database/mysql"

	"github.com/cadicallegari/user/pkg/xdatabase/xsql"
)

var DriverName = "mysql"

type Config struct {
	xsql.Config
}

func (cfg *Config) setDefaults() {
	if cfg.Config.MigrationDriverFn == nil {
		cfg.Config.MigrationDriverFn = func(db *xsql.DB) (database.Driver, error) {
			return migrationmysql.WithInstance(db.DB, &migrationmysql.Config{
				MigrationsTable: cfg.MigrationsTable,
			})
		}
	}
}

func Connect(cfg *Config) (*xsql.DB, error) {
	if cfg == nil {
		cfg = new(Config)
	}
	cfg.setDefaults()

	db, err := xsql.Open(DriverName, &cfg.Config)
	if err != nil {
		if errorNumber(err) == 1049 {
			// if Error 1049: Unknown database
			// create the database and try again
			err := createDB(cfg.Config)
			if err != nil {
				return nil, err
			}
			return Connect(cfg)
		}
		if db == nil {
			return nil, err
		}

		err := xsql.WaitOK(db, func(err error, attempt int) bool {
			if cfg.Logger != nil {
				cfg.Logger.WithError(err).Warn("unable to connect to mysql")
			}

			return attempt < 5
		})
		if err != nil {
			return nil, err
		}
	}

	version, dirty, err := xsql.Migration(DriverName, db, &cfg.Config)
	if err != nil && err != xsql.ErrMigrationDisabled {
		return nil, err
	}
	if version > 0 && cfg.Logger != nil {
		if dirty {
			cfg.Logger.Warnf("mysql database is dirty, version: %d", version)
		} else {
			cfg.Logger.Infof("mysql database is ok, version: %d", version)
		}
	}

	if cfg.Logger != nil {
		cfg.Logger.Info("connected to mysql")
	}

	return db, nil
}

func createDB(cfg xsql.Config) error {
	mysqlCfg, err := mysql.ParseDSN(cfg.URL)
	if err != nil {
		return err
	}
	dbname := mysqlCfg.DBName
	mysqlCfg.DBName = ""

	// update database locally and try to connect
	cfg.URL = mysqlCfg.FormatDSN()

	db, err := xsql.Open("mysql", &cfg)
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT COLLATE %s", dbname, cfg.Charset))
	if err != nil {
		return err
	}
	return db.Close()
}

func errorNumber(err error) uint16 {
	var merr *mysql.MySQLError
	if errors.As(err, &merr) {
		return merr.Number
	}
	return 0
}
