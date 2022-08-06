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

type UserStorage struct {
	db *sqlx.DB
}

var TimeNow = func() time.Time {
	return time.Now().UTC()
}

var baseSelect = sq.Select(
	"u.id",
	"u.first_name",
	"u.last_name",
	"u.nickname",
	"u.email",
	"u.encoded_password",
	"u.country",
	"u.created_at",
	"u.updated_at",
).From("users u")

func NewStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) Save(ctx context.Context, usr *user.User) (*user.User, error) {
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
		).
		Suffix(`AS u ON DUPLICATE KEY UPDATE
			first_name = u.first_name,
			last_name = u.last_name,
			nickname = u.nickname,
			encoded_password = u.encoded_password,
			country = u.country
		`)

	_, err := q.RunWith(s.db).ExecContext(ctx)
	if err != nil {
		xlogger.Logger(ctx).
			WithField("query", sq.DebugSqlizer(q)).
			WithError(err).
			Error("unable to save user")
		return nil, err
	}

	return s.Get(ctx, usr.ID)
}

func (s *UserStorage) List(ctx context.Context, opts *user.ListOptions) (*user.List, error) {
	q := baseSelect.OrderBy("email ASC")

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
	q := baseSelect.Where(sq.Eq{"u.id": id})

	query, args := q.MustSql()

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u user.User

	for rows.Next() {
		err := rows.StructScan(&u)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if u.ID == "" {
		return nil, user.ErrNotFound
	}

	return &u, nil
}

func (s *UserStorage) Delete(ctx context.Context, usr *user.User) error {
	q := sq.Delete("users").Where(sq.Eq{"id": usr.ID})

	res, err := q.RunWith(s.db).ExecContext(ctx)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		err = user.ErrNotFound
	}

	return err
}
