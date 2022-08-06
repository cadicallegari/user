package xmysqltest

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/suite"

	"github.com/cadicallegari/user/pkg/xdatabase/xsql"
	"github.com/cadicallegari/user/pkg/xdatabase/xsql/xmysql"
)

var debug = flag.Bool("debug", false, "do not remove the db used for the test")

func init() {
	rand.Seed(time.Now().UnixNano())
}

type MysqlTestSuite struct {
	suite.Suite
	DB     *xsql.DB
	DBName string
}

func (s *MysqlTestSuite) SetupTest(mysqlURL, migrationDir string) {
	mysqlCfg, err := mysqldriver.ParseDSN(mysqlURL)
	if err != nil {
		s.FailNow("unable to parse mysql url")
	}
	mysqlCfg.DBName += "-" + fmt.Sprint(rand.Int())
	s.DBName = mysqlCfg.DBName

	if *debug {
		s.T().Logf("Using table %s", s.DBName)
	}

	cfg := &xmysql.Config{}
	cfg.URL = mysqlCfg.FormatDSN()
	cfg.MigrationsDir = migrationDir
	cfg.RunMigration = true

	s.DB, err = xmysql.Connect(cfg)
	if !s.NoError(err) {
		s.FailNow("unable to connect to mysql")
	}
}

func (s *MysqlTestSuite) TearDownTest() {
	if !*debug {
		s.DB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", s.DBName))
	}
	s.DB.Close()
}
