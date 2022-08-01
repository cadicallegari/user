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

func (s *service) List(ctx context.Context, opts *ListOptions) (*List, error) {
	return s.storage.List(ctx, opts)
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

func (s *service) Update(ctx context.Context, usr *User) (*User, error) {
	u, err := s.storage.Update(ctx, usr)
	if err != nil {
		return nil, err
	}

	err = s.eventService.UserUpdated(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) Delete(ctx context.Context, usr *User) error {
	err := s.storage.Delete(ctx, usr)
	if err != nil {
		return err
	}

	err = s.eventService.UserDeleted(ctx, usr)
	if err != nil {
		return err
	}

	return nil
}
