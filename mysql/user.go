package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/cadicallegari/user"
	"github.com/cadicallegari/user/pkg/xlogger"
)

var TimeNow = func() time.Time {
	return time.Now().UTC()
}

type UserStorage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) Create(ctx context.Context, usr *user.User) (*user.User, error) {
	if usr.ID == "" {
		usr.ID = uuid.NewString()
	}

	q := sq.Insert("users").
		Columns(
			"id",
			"first_name",
			"last_name",
			"nickname",
			"email",
			"encoded_password",
			"country",
		).
		Values(
			usr.ID,
			usr.FirstName,
			usr.LastName,
			usr.Nickname,
			usr.Email,
			usr.EncodedPassword,
			usr.Country,
		)

	_, err := q.RunWith(s.db).ExecContext(ctx)
	if err != nil {
		xlogger.Logger(ctx).
			WithField("query", sq.DebugSqlizer(q)).
			WithError(err).
			Error("unable to save user")
		return nil, err
	}

	return usr, nil
}

func (s *UserStorage) List(ctx context.Context, opts *user.ListOptions) (*user.List, error) {
	q := sq.
		Select(
			"u.id",
			"u.first_name",
			"u.last_name",
			"u.nickname",
			"u.email",
			"u.encoded_password",
			"u.country",
			"u.created_at",
			"u.updated_at",
		).
		From("users u").
		OrderBy("email ASC")

	// TODO: filter and paginate
	query, args := q.MustSql()

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		xlogger.Logger(ctx).
			WithField("query", sq.DebugSqlizer(q)).
			WithError(err).
			Error("unable to get users")
		return nil, err
	}
	defer rows.Close()

	us := make([]*user.User, 0)

	for rows.Next() {
		var u user.User
		err := rows.StructScan(&u)
		if err != nil {
			return nil, err
		}

		us = append(us, &u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user.List{
		Users: us,
		Total: uint64(len(us)),
	}, nil
}

func (s *UserStorage) Get(ctx context.Context, id string) (*user.User, error) {
	return &user.User{}, nil
}

func (s *UserStorage) Update(_ context.Context, usr *user.User) (*user.User, error) {
	return &user.User{}, nil
}

func (s *UserStorage) Delete(_ context.Context, usr *user.User) error {
	return nil
}
