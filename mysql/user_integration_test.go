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

func (s *UserStorageSuite) Test_Delete() {
	originalUser := user.User{
		FirstName:       "firstName",
		LastName:        "lastName",
		Nickname:        "nickName",
		Email:           "email@mail.com",
		EncodedPassword: "encoded",
		Country:         "BR",
	}
	err := s.storage.Delete(s.ctx, &originalUser)
	s.ErrorIs(err, user.ErrNotFound)

	createdUser, err := s.storage.Save(s.ctx, &originalUser)
	s.NoError(err)
	s.NotEmpty(createdUser.ID)

	got, err := s.storage.Get(s.ctx, createdUser.ID)
	s.NoError(err)
	s.Equal(got.ID, createdUser.ID)

	err = s.storage.Delete(s.ctx, createdUser)
	s.NoError(err)

	got, err = s.storage.Get(s.ctx, createdUser.ID)
	s.ErrorIs(err, user.ErrNotFound)
	s.Nil(got)
}

func (s *UserStorageSuite) Test_List() {
	users := s.createUsers([]string{
		"DE", "DE", "BR", "UK", "UK", "UK", "ES", "PT",
	})

	lr, err := s.storage.List(s.ctx, &user.ListOptions{})
	if s.NoError(err) {
		s.Equal(len(users), int(lr.Total))
		s.Nil(lr.PrevPage)
		s.Nil(lr.NextPage)
	}

	lr, err = s.storage.List(s.ctx, &user.ListOptions{
		Country: "UK",
	})
	if s.NoError(err) {
		s.Equal(3, int(lr.Total))
		s.Nil(lr.PrevPage)
		s.Nil(lr.NextPage)
	}

	lr, err = s.storage.List(s.ctx, &user.ListOptions{
		Search: "u5",
	})
	if s.NoError(err) {
		s.Equal(1, int(lr.Total))
		s.Nil(lr.PrevPage)
		s.Nil(lr.NextPage)
	}
}

func (s *UserStorageSuite) Test_List_Pagination() {
	users := s.createUsers([]string{
		"DE", "DE", "BR", "UK", "UK", "UK", "ES", "PT",
	})

	_ = users

	lr, err := s.storage.List(s.ctx, &user.ListOptions{
		PerPage: 1,
	})
	if s.NoError(err) {
		s.Equal(8, int(lr.Total))
		s.Len((lr.Users), 1)
		s.Nil(lr.PrevPage)
		s.NotNil(lr.NextPage)
	}
}

func (s *UserStorageSuite) createUsers(countries []string) []*user.User {
	users := make([]*user.User, 0)

	for i, country := range countries {
		u, err := s.storage.Save(s.ctx, &user.User{
			FirstName:       fmt.Sprintf("%d name %s", i, s.T().Name()),
			LastName:        fmt.Sprintf("%d last %s", i, s.T().Name()),
			Nickname:        fmt.Sprintf("%d nick %s", i, s.T().Name()),
			Email:           fmt.Sprintf("u%d@%s", i, s.T().Name()),
			EncodedPassword: fmt.Sprintf("encoded-%d", i),
			Country:         country,
		})

		s.NoError(err)
		s.NotEmpty(u.ID)

		users = append(users, u)
	}

	return users
}
