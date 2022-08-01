package mem

import (
	"context"

	"github.com/cadicallegari/user"
)

type memStorage struct {
}

func NewStorage() *memStorage {
	return &memStorage{}
}

func (s *memStorage) List(context.Context, *user.ListOptions) (*user.List, error) {

	return nil, nil
}

func (s *memStorage) Get(_ context.Context, id string) (*user.User, error) {
	return nil, nil
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
