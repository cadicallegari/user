package user

import (
	"context"
	"errors"
	"time"
)

const (
	DefaultPerPage = 25
)

// TODO: it also could be moved to pkg/xerrors
var (
	ErrNotFound      = errors.New("user not found")
	ErrAlreadyExists = errors.New("user already exists")
)

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListOptions struct {
	PerPage uint64 `schema:"per_page"`
	Cursor  string `schema:"cursor"`
	Country string `schema:"country"`
	// Search is used for text search in some fields
	Search string `schema:"search"`
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		PerPage: uint64(DefaultPerPage),
	}
}

type List struct {
	Users      []*User `json:"users"`
	Total      uint64  `json:"total"`
	PrevCursor *string `json:"prev_cursor"`
	NextCursor *string `json:"next_cursor"`
}

type Service interface {
	Get(_ context.Context, id string) (*User, error)
	List(context.Context, *ListOptions) (*List, error)
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, *User) error
}

//go:generate mockgen -package mock -mock_names Storage=Storage -destination mock/storage.go github.com/cadicallegari/user Storage
type Storage interface {
	Get(_ context.Context, id string) (*User, error)
	List(context.Context, *ListOptions) (*List, error)
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, *User) error
}

//go:generate mockgen -package mock -mock_names EventService=EventService -destination mock/event.go github.com/cadicallegari/user EventService
type EventService interface {
	UserCreated(context.Context, *User) error
	UserUpdated(context.Context, *User) error
	UserDeleted(context.Context, *User) error
}
