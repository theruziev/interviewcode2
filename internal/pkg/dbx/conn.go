package dbx

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Querier interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type PostgresOpt struct {
	DSN string `help:"postgresql dsn" required:"" env:"DSN"`
}

type Dbx struct {
	pool *pgxpool.Pool
	mu   sync.Mutex
}

func NewDbx() *Dbx {
	return &Dbx{}
}

func (db *Dbx) Connect(ctx context.Context, dsn string) error {
	if db.pool != nil {
		return nil
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}
	db.pool = pool

	return nil
}

func (db *Dbx) GetPool() *pgxpool.Pool {
	return db.pool
}

func (db *Dbx) Close(_ context.Context) error {
	if db.pool == nil {
		return nil
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	db.pool.Close()
	return nil
}

func (db *Dbx) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.pool.Begin(ctx)
}

func (db *Dbx) BeginFunc(ctx context.Context, f func(pgx.Tx) error) (err error) {
	return pgx.BeginFunc(ctx, db.pool, f)
}

func (db *Dbx) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.pool.Exec(ctx, sql, args...)
}

func (db *Dbx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return db.pool.Query(ctx, sql, args...)
}

func (db *Dbx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return db.pool.QueryRow(ctx, sql, args...)
}

func (db *Dbx) GetConn(ctx context.Context) Querier {
	return GetConnOrTx(ctx, db)
}
