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
	// publish into user.created topic for example
	return nil
}

func (s *memEventService) UserUpdated(context.Context, *user.User) error {
	// publish into user.updated topic for example
	return nil
}

func (s *memEventService) UserDeleted(context.Context, *user.User) error {
	// publish into user.deleted topic for example
	return nil
}
