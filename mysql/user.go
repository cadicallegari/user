package mysql

import (
	"context"
	"strconv"
	"time"

	"github.com/cadicallegari/user"
	"github.com/jmoiron/sqlx"
)

var TimeNow = func() time.Time {
	return time.Now().UTC()
}

type UserStorage struct {
	db    *sqlx.DB
	users map[string]*user.User
}

func NewStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db:    db,
		users: make(map[string]*user.User),
	}
}

func (s *UserStorage) List(context.Context, *user.ListOptions) (*user.List, error) {
	us := []*user.User{}
	for _, u := range s.users {
		us = append(us, u)
	}

	return &user.List{
		Total: uint64(len(us)),
		Users: us,
	}, nil
}

func (s *UserStorage) Get(ctx context.Context, id string) (*user.User, error) {
	u, ok := s.users[id]
	if !ok {
		return nil, user.ErrNotFound
	}

	return u, nil
}

func (s *UserStorage) Create(_ context.Context, usr *user.User) (*user.User, error) {
	if usr.ID == "" {
		usr.ID = strconv.Itoa(len(s.users) + 1)
		usr.CreatedAt = TimeNow()
		usr.UpdatedAt = TimeNow()
	}

	s.users[usr.ID] = usr

	return usr, nil
}

func (s *UserStorage) Update(_ context.Context, usr *user.User) (*user.User, error) {
	s.users[usr.ID] = usr

	return usr, nil
}

func (s *UserStorage) Delete(_ context.Context, usr *user.User) error {
	delete(s.users, usr.ID)
	return nil
}
