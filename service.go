package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	storage      Storage
	eventService EventService

	passwordCost int
}

func NewService(storage Storage, eventService EventService, passwordCost int) *service {
	return &service{
		storage:      storage,
		eventService: eventService,
		passwordCost: passwordCost,
	}
}

func (s *service) encryptPassword(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), s.passwordCost)

	return string(bytes), err
}

func (s *service) List(ctx context.Context, opts *ListOptions) (*List, error) {
	return s.storage.List(ctx, opts)
}

func (s *service) Get(ctx context.Context, id string) (*User, error) {
	return s.storage.Get(ctx, id)
}

func (s *service) Save(ctx context.Context, usr *User) (*User, error) {
	l, err := s.List(ctx, &ListOptions{Search: usr.Email})
	if err != nil {
		return nil, err
	}

	if l.Total > 0 {
		return nil, ErrAlreadyExists
	}

	if usr.Password != "" {
		encoded, err := s.encryptPassword(usr.Password)
		if err != nil {
			return nil, ErrInvalid
		}
		usr.EncodedPassword = encoded
		usr.Password = ""
	}

	u, err := s.storage.Save(ctx, usr)
	if err != nil {
		return nil, err
	}

	// dual write problem, can be solved using listen yourself or outbox pattern for example

	err = s.eventService.UserCreated(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) Update(ctx context.Context, usr *User) (*User, error) {
	if usr.Password != "" {
		encoded, err := s.encryptPassword(usr.Password)
		if err != nil {
			return nil, ErrInvalid
		}
		usr.EncodedPassword = encoded
		usr.Password = ""
	}

	u, err := s.storage.Update(ctx, usr)
	if err != nil {
		return nil, err
	}

	// dual write problem, can be solved using listen yourself or outbox pattern for example

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

	// dual write problem, can be solved using listen yourself or outbox pattern for example

	err = s.eventService.UserDeleted(ctx, usr)
	if err != nil {
		return err
	}

	return nil
}
