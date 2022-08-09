package user

import (
	"context"
	"errors"
	"time"
)

const (
	DefaultPerPage = 25
)

var (
	ErrNotFound      = errors.New("not_found")
	ErrInvalid       = errors.New("invalid")
	ErrAlreadyExists = errors.New("already exists")
)

type User struct {
	ID              string    `json:"id"`
	FirstName       string    `json:"first_name" db:"first_name"`
	LastName        string    `json:"last_name" db:"last_name"`
	Nickname        string    `json:"nickname" db:"nickname"`
	Password        string    `json:"password,omitempty"`
	EncodedPassword string    `json:"-" db:"encoded_password"`
	Email           string    `json:"email"`
	Country         string    `json:"country"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type ListOptions struct {
	Page    uint64 `schema:"page"`
	PerPage uint64 `schema:"per_page"`

	Country string `schema:"country"`
	// Search is used for text search in the email field for now
	Search string `schema:"search"`
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		PerPage: uint64(DefaultPerPage),
	}
}

type List struct {
	Users    []*User `json:"users"`
	Total    uint64  `json:"total"`
	PrevPage *uint64 `json:"prev_page"`
	NextPage *uint64 `json:"next_page"`
}

type Service interface {
	Get(_ context.Context, id string) (*User, error)
	List(context.Context, *ListOptions) (*List, error)
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, *User) error
}

//go:generate mockgen -package mock -mock_names Storage=Storage -destination mock/storage.go github.com/cadicallegari/user Storage
type Storage interface {
	Get(_ context.Context, id string) (*User, error)
	List(context.Context, *ListOptions) (*List, error)
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, *User) error
}

//go:generate mockgen -package mock -mock_names EventService=EventService -destination mock/event.go github.com/cadicallegari/user EventService
type EventService interface {
	UserCreated(context.Context, *User) error
	UserUpdated(context.Context, *User) error
	UserDeleted(context.Context, *User) error
}
