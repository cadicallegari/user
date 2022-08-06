//go:build integration

package mysql_test

import (
	"context"
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

func (s *UserStorageSuite) Test_Get_NotFound() {
	got, err := s.storage.Get(s.ctx, "inexistent")
	s.ErrorIs(err, user.ErrNotFound)
	s.Nil(got)
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
	if s.NoError(err) {
		s.NotEmpty(gotUser.ID)
		s.False(gotUser.CreatedAt.IsZero())
		s.False(gotUser.UpdatedAt.IsZero())
		s.Equal(firstName, gotUser.FirstName)
		s.Equal(lastname, gotUser.LastName)
		s.Equal(nickName, gotUser.Nickname)
		s.Equal(email, gotUser.Email)
		s.Equal(encoded, gotUser.EncodedPassword)
		s.Equal(country, gotUser.Country)
	}

	got, err := s.storage.Get(s.ctx, gotUser.ID)
	if s.NoError(err) {
		s.Equal(gotUser.ID, got.ID)
		s.Equal(firstName, got.FirstName)
		s.Equal(lastname, got.LastName)
		s.Equal(nickName, got.Nickname)
		s.Equal(email, got.Email)
	}
}

func (s *UserStorageSuite) Test_Update() {
	originalEmail := "email@mail.com"
	originalUser := user.User{
		FirstName:       "firstName",
		LastName:        "lastName",
		Nickname:        "nickName",
		Email:           originalEmail,
		EncodedPassword: "encoded",
		Country:         "BR",
	}

	gotUser, err := s.storage.Save(s.ctx, &originalUser)
	if s.NoError(err) {
		s.NotEmpty(gotUser.ID)
		s.Equal(gotUser.FirstName, originalUser.FirstName)
		s.Equal(gotUser.LastName, originalUser.LastName)
		s.Equal(gotUser.Nickname, originalUser.Nickname)
		s.Equal(gotUser.Email, originalEmail)
		s.Equal(gotUser.EncodedPassword, originalUser.EncodedPassword)
		s.Equal(gotUser.Country, originalUser.Country)
		s.False(gotUser.CreatedAt.IsZero())
		s.False(gotUser.UpdatedAt.IsZero())
	}

	userToUpdate := user.User{
		ID:              gotUser.ID,
		FirstName:       "updated firstName",
		LastName:        "updated lastName",
		Nickname:        "updated nickname",
		EncodedPassword: "updated encoded",
		Email:           "email_is_not@updated.com",
		Country:         "DE",
	}

	gotUser, err = s.storage.Save(s.ctx, &userToUpdate)
	if s.NoError(err) {
		s.Equal(gotUser.ID, userToUpdate.ID)
		s.Equal(gotUser.FirstName, userToUpdate.FirstName)
		s.Equal(gotUser.LastName, userToUpdate.LastName)
		s.Equal(gotUser.Nickname, userToUpdate.Nickname)
		s.Equal(gotUser.Email, originalEmail)
		s.Equal(gotUser.EncodedPassword, userToUpdate.EncodedPassword)
		s.Equal(gotUser.Country, userToUpdate.Country)
		s.False(gotUser.UpdatedAt.IsZero())
		s.False(gotUser.CreatedAt.IsZero())
	}

	got, err := s.storage.Get(s.ctx, gotUser.ID)
	if s.NoError(err) {
		s.Equal(got.ID, userToUpdate.ID)
		s.Equal(got.FirstName, userToUpdate.FirstName)
		s.Equal(got.LastName, userToUpdate.LastName)
		s.Equal(got.Nickname, userToUpdate.Nickname)
		s.Equal(got.Email, originalEmail)
	}

}
