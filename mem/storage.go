package mem

import (
	"context"

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
	l := &user.List{}
	return l, nil
}

func (s *memStorage) Create(context.Context, *user.User) (*user.User, error) {

	return nil, nil
}

func (s *memStorage) Update(context.Context, *user.User) (*user.User, error) {
	return nil, nil
}

func (s *memStorage) Delete(context.Context, *user.User) error {
	return nil
}
