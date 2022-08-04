//go:build integration

package mysql_test

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/cadicallegari/user"
	"github.com/cadicallegari/user/mysql"
	"github.com/cadicallegari/user/pkg/xdatabase/xsql"
	"github.com/cadicallegari/user/pkg/xdatabase/xsql/xmysql"
	"github.com/cadicallegari/user/pkg/xlogger"
)

var debug = flag.Bool("debug", false, "do not remove the db used for the test")

func init() {
	rand.Seed(time.Now().UnixNano())

	// To make the test easier let's fix a date
	mysql.TimeNow = func() time.Time {
		return time.Date(2022, time.June, 23, 0, 30, 0, 0, time.UTC)
	}
}

type UserStorageSuite struct {
	suite.Suite
	db      *xsql.DB
	dbName  string
	storage *mysql.UserStorage
	ctx     context.Context
}

func TestUserStorage(t *testing.T) {
	suite.Run(t, new(UserStorageSuite))
}

// TODO: move it to pkg
func (s *UserStorageSuite) SetupTest() {
	mysqlURL := os.Getenv("USER_MYSQL_URL")
	if mysqlURL == "" {
		s.FailNow("envvar USER_MYSQL_URL is empty or missing")
	}
	mysqlCfg, err := mysqldriver.ParseDSN(mysqlURL)
	if err != nil {
		s.FailNow("unable to parse mysql url")
	}
	mysqlCfg.DBName += "-" + fmt.Sprint(rand.Int())
	s.dbName = mysqlCfg.DBName

	if *debug {
		s.T().Logf("Using table %s", s.dbName)
	}

	cfg := &xmysql.Config{}
	cfg.URL = mysqlCfg.FormatDSN()
	cfg.MigrationsDir = os.Getenv("USER_MYSQL_MIGRATIONS_DIR")
	cfg.RunMigration = true

	s.db, err = xmysql.Connect(cfg)
	if !s.NoError(err) {
		s.FailNow("unable to connect to mysql")
	}

	s.storage = mysql.NewStorage(s.db)

	ctx := context.Background()
	logger := logrus.StandardLogger()
	s.ctx = xlogger.SetLogger(ctx, logger.WithField("test", "test"))
}

func (s *UserStorageSuite) TearDownTest() {
	if !*debug {
		s.db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", s.dbName))
	}
	s.db.Close()
}

func (s *UserStorageSuite) Test_Create() {
	firstName := "firstName"
	lastname := "lastName"
	nickName := "nickName"
	email := "email@mail.com"
	encoded := "234kj;salkfj"
	country := "DE"

	u := user.User{
		FirstName:       firstName,
		LastName:        lastname,
		Nickname:        nickName,
		Email:           email,
		EncodedPassword: encoded,
		Country:         country,
	}

	gotUser, err := s.storage.Create(s.ctx, &u)
	if !s.NoError(err) {
		s.T().FailNow()
	}

	s.NotEmpty(gotUser.ID)
	s.False(gotUser.CreatedAt.IsZero())
	s.False(gotUser.UpdatedAt.IsZero())
	s.Equal(firstName, gotUser.FirstName)
	s.Equal(lastname, gotUser.LastName)
	s.Equal(nickName, gotUser.Nickname)
	s.Equal(email, gotUser.Email)
	s.Equal(encoded, gotUser.EncodedPassword)
	s.Equal(country, gotUser.Country)

	listResp, err := s.storage.List(s.ctx, &user.ListOptions{})
	if s.NoError(err) && s.Len(listResp.Users, 1) {
		u := listResp.Users[0]
		s.Equal(u.ID, gotUser.ID)
		s.Equal(firstName, gotUser.FirstName)
		s.Equal(lastname, gotUser.LastName)
		s.Equal(nickName, gotUser.Nickname)
		s.Equal(email, gotUser.Email)
	}
}
