package dbx

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type contextKey string

const txKey = contextKey("tx")

func FromContext(ctx context.Context) Querier {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx
	}
	return nil
}

func WithContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetConnOrTx(ctx context.Context, q Querier) Querier {
	tx := FromContext(ctx)
	if tx != nil {
		return tx
	}
	return q
}
