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

func (s *service) Create(ctx context.Context, usr *User) (*User, error) {
	u, err := s.storage.Create(ctx, usr)
	if err != nil {
		return nil, err
	}

	err = s.eventService.UserCreated(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) Update(context.Context, *User) (*User, error) {
	return nil, nil
}

func (s *service) Delete(context.Context, *User) error {
	return nil
}
