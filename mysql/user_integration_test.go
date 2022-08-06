//go:build integration

package mysql_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/cadicallegari/user"
	"github.com/cadicallegari/user/mysql"
	"github.com/cadicallegari/user/pkg/xdatabase/xsql/xmysqltest"
	"github.com/cadicallegari/user/pkg/xlogger"
)

type UserStorageSuite struct {
	xmysqltest.MysqlTestSuite
	storage *mysql.UserStorage
	ctx     context.Context
}

func TestUserStorage(t *testing.T) {
	suite.Run(t, new(UserStorageSuite))
}

func (s *UserStorageSuite) SetupTest() {
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

	gotUser, err := s.storage.Save(s.ctx, &u)
	fmt.Println(err)
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
		s.Equal(gotUser.ID, u.ID)
		s.Equal(firstName, u.FirstName)
		s.Equal(lastname, u.LastName)
		s.Equal(nickName, u.Nickname)
		s.Equal(email, u.Email)
	}
}

func (s *UserStorageSuite) Test_Get() {
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

	createdUser, err := s.storage.Save(s.ctx, &u)
	if !s.NoError(err) {
		s.T().FailNow()
	}

	s.NotEmpty(createdUser.ID)
	s.False(createdUser.CreatedAt.IsZero())
	s.False(createdUser.UpdatedAt.IsZero())
	s.Equal(firstName, createdUser.FirstName)
	s.Equal(lastname, createdUser.LastName)
	s.Equal(nickName, createdUser.Nickname)
	s.Equal(email, createdUser.Email)
	s.Equal(encoded, createdUser.EncodedPassword)
	s.Equal(country, createdUser.Country)

	gotUser, err := s.storage.Get(s.ctx, createdUser.ID)
	if s.NoError(err) {
		s.Equal(createdUser.ID, gotUser.ID)
		s.Equal(firstName, gotUser.FirstName)
		s.Equal(lastname, gotUser.LastName)
		s.Equal(nickName, gotUser.Nickname)
		s.Equal(email, gotUser.Email)
	}
}

func (s *UserStorageSuite) Test_Get_NotFound() {
	got, err := s.storage.Get(s.ctx, "inexistent")
	s.Nil(got)
	s.ErrorIs(err, user.ErrNotFound)
}
