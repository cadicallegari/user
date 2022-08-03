package xsql

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/cadicallegari/user/pkg/xlogger"
)

type Commiter interface {
	Commit() error
}

type Rollbacker interface {
	Rollback() error
}

type Transaction interface {
	Commiter
	Rollbacker
}

func CommitOrRollback(tx Transaction, err error) error {
	if err != nil {
		rerr := tx.Rollback()
		_ = rerr // ignore for now.
		return err
	}
	return tx.Commit()
}

// AffectedRows returns a channel with the number of records affected by the query.
func AffectedRows(ctx context.Context, tx sq.BaseRunner, total *uint64, query sq.SelectBuilder) chan uint64 {
	totalCh := make(chan uint64)
	go func() {
		defer close(totalCh)
		if total == nil {
			q := sq.Select("COUNT(*)").FromSelect(query, "q")

			row := q.RunWith(tx).QueryRowContext(ctx)
			err := row.Scan(&total)
			if err != nil {
				xlogger.Logger(ctx).
					WithError(err).
					Error("unable to get total of affected rows")
				total = new(uint64)
			}
		}
		totalCh <- *total
	}()
	return totalCh
}
