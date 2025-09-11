package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (q *Queries) BeginTx(ctx context.Context) (*Queries, pgx.Tx, error) {
	tx, ok := q.db.(interface {
		Begin(context.Context) (pgx.Tx, error)
	})
	if !ok {
		return nil, nil, errors.New("db does not support Begin")
	}

	pgxTx, err := tx.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	return q.WithTx(pgxTx), pgxTx, nil
}
