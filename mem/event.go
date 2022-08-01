package mem

import (
	"context"

	"github.com/cadicallegari/user"
)

type memEventService struct {
}

func NewEventService() *memEventService {
	return &memEventService{}
}

func (s *memEventService) UserCreated(context.Context, *user.User) error {
	return nil
}

func (s *memEventService) UserUpdated(context.Context, *user.User) error {
	return nil
}

func (s *memEventService) UserDeleted(context.Context, *user.User) error {
	return nil
}
