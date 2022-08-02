package mem

import (
	"context"
	"strconv"

	"github.com/cadicallegari/user"
)

type memStorage struct {
	users map[string]*user.User
}

func NewStorage() *memStorage {
	return &memStorage{
		users: make(map[string]*user.User),
	}
}

func (s *memStorage) List(context.Context, *user.ListOptions) (*user.List, error) {
	us := []*user.User{}
	for _, u := range s.users {
		us = append(us, u)
	}

	return &user.List{
		Total: uint64(len(us)),
		Users: us,
	}, nil
}

func (s *memStorage) Get(ctx context.Context, id string) (*user.User, error) {
	u, ok := s.users[id]
	if !ok {
		return nil, user.ErrNotFound
	}

	return u, nil
}

func (s *memStorage) Create(_ context.Context, usr *user.User) (*user.User, error) {
	if usr.ID == "" {
		usr.ID = strconv.Itoa(len(s.users) + 1)
	}

	s.users[usr.ID] = usr

	return usr, nil
}

func (s *memStorage) Update(_ context.Context, usr *user.User) (*user.User, error) {
	s.users[usr.ID] = usr

	return usr, nil
}

func (s *memStorage) Delete(_ context.Context, usr *user.User) error {
	delete(s.users, usr.ID)
	return nil
}
