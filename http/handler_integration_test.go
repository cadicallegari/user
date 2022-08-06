//go:build integration

package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"testing"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/cadicallegari/user"
	userHttp "github.com/cadicallegari/user/http"
	"github.com/cadicallegari/user/mem"
	"github.com/cadicallegari/user/mysql"
	"github.com/cadicallegari/user/pkg/xdatabase/xsql"
	"github.com/cadicallegari/user/pkg/xdatabase/xsql/xmysql"
	"github.com/cadicallegari/user/pkg/xhttp"
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
	logger := logrus.StandardLogger()
	logger.SetOutput(io.Discard)

	log := logger.WithFields(nil)
	ctx := xlogger.SetLogger(context.TODO(), log)

	svc := user.NewService(s.storage, mem.NewEventService())
	r := xhttp.NewRouter(log)

	_ = userHttp.NewUserHandler(r, svc)

	buf, err := json.Marshal(user.User{
		FirstName: "first name",
		Email:     "email",
	})
	if err != nil {
		s.T().Fatalf("unable to marshal bid %q", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(buf))
	if err != nil {
		s.T().Fatalf("unable to create request %q", err)
	}

	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		s.T().Fatalf("got %d status creating user, want %d", resp.StatusCode, http.StatusCreated)
	}
}

func (s *UserStorageSuite) Test_Get_NotFound() {
	logger := logrus.StandardLogger()
	logger.SetOutput(io.Discard)

	log := logger.WithFields(nil)
	ctx := xlogger.SetLogger(context.TODO(), log)

	svc := user.NewService(s.storage, mem.NewEventService())
	r := xhttp.NewRouter(log)

	_ = userHttp.NewUserHandler(r, svc)

	req, err := http.NewRequest(http.MethodGet, "/v1/users/notfound", nil)
	if err != nil {
		s.T().Fatalf("unable to create request %q", err)
	}

	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("RESPONSE:\n%s", string(respDump))

	if resp.StatusCode != http.StatusNotFound {
		fmt.Println()
		s.T().Fatalf("got %d status creating user, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
