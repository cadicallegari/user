package user

import (
	"context"
	"time"
)

const (
	DefaultPerPage = 25
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
	Role    string `schema:"role"`
	Search  string `schema:"search"`
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
	List(context.Context, *ListOptions) (*List, error)
	Get(_ context.Context, id string) (*User, error)
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, *User) error
}

type Storage interface {
	List(context.Context, *ListOptions) (*List, error)
	Get(_ context.Context, id string) (*User, error)
	Create(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, *User) error
}

type EventService interface {
	UserCreated(context.Context, *User) error
	UserUpdated(context.Context, *User) error
	UserDeleted(context.Context, *User) error
}
