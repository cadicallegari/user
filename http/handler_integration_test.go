//go:build integration

package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/cadicallegari/user"
	userHttp "github.com/cadicallegari/user/http"
	"github.com/cadicallegari/user/mem"
	"github.com/cadicallegari/user/mysql"
	"github.com/cadicallegari/user/pkg/xdatabase/xsql/xmysqltest"
	"github.com/cadicallegari/user/pkg/xhttp"
	"github.com/cadicallegari/user/pkg/xlogger"
)

type HttpIntegrationTestSuite struct {
	xmysqltest.MysqlTestSuite
	storage *mysql.UserStorage
	ctx     context.Context
}

func TestHttpIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(HttpIntegrationTestSuite))
}

func (s *HttpIntegrationTestSuite) SetupTest() {
	mysqlURL := os.Getenv("USER_MYSQL_URL")
	if mysqlURL == "" {
		s.FailNow("envvar USER_MYSQL_URL is empty or missing")
	}

	s.MysqlTestSuite.SetupTest(mysqlURL, os.Getenv("USER_MYSQL_MIGRATIONS_DIR"))

	s.storage = mysql.NewStorage(s.DB)

	ctx := context.Background()
	logger := logrus.StandardLogger()
	s.ctx = xlogger.SetLogger(ctx, logger.WithField("test", "test"))
}

func (s *HttpIntegrationTestSuite) Test_Create() {
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

func (s *HttpIntegrationTestSuite) Test_Get_NotFound() {
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
