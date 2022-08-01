package user

import "context"

type service struct {
	storage      Storage
	eventService EventService
}

func NewService(storage Storage, eventService EventService) *service {
	return &service{
		storage:      storage,
		eventService: eventService,
	}
}

func (s *service) List(context.Context, *ListOptions) (*List, error) {

	return nil, nil
}

func (s *service) Get(_ context.Context, id string) (*User, error) {
	return nil, nil
}

func (s *service) Create(context.Context, *User) (*User, error) {
	return nil, nil
}

func (s *service) Update(context.Context, *User) (*User, error) {
	return nil, nil
}

func (s *service) Delete(context.Context, *User) error {
	return nil
}
