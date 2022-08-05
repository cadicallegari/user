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

func (s *service) Get(ctx context.Context, id string) (*User, error) {
	return s.storage.Get(ctx, id)
}

func (s *service) Save(ctx context.Context, usr *User) (*User, error) {
	u, err := s.storage.Save(ctx, usr)
	if err != nil {
		return nil, err
	}

	// TODO: add something about dual write

	err = s.eventService.UserCreated(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) Update(ctx context.Context, usr *User) (*User, error) {
	u, err := s.storage.Save(ctx, usr)
	if err != nil {
		return nil, err
	}

	// TODO: add something about dual write

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

	// TODO: add something about dual write

	err = s.eventService.UserDeleted(ctx, usr)
	if err != nil {
		return err
	}

	return nil
}
