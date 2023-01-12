package db

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/theruziev/oson_auth/internal/model"
	"github.com/theruziev/oson_auth/internal/pkg/dbx"
	"github.com/theruziev/oson_auth/internal/pkg/logging"
)

const outboxTable = "outbox"

type OutBoxStore struct {
	db *dbx.Dbx
}

func NewOutBoxStore(db *dbx.Dbx) *OutBoxStore {
	return &OutBoxStore{
		db: db,
	}
}

func (o *OutBoxStore) Add(ctx context.Context, messages ...*model.OutBox) error {
	if len(messages) == 0 {
		return nil
	}

	builder := pgsql.Insert(outboxTable)

	for _, msg := range messages {
		builder = builder.SetMap(map[string]interface{}{
			"topic":      msg.Topic,
			"status":     msg.Status,
			"msg":        msg.Data,
			"created_at": msg.CreatedAt,
			"updated_at": msg.UpdatedAt,
		})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	if _, err := o.db.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (o *OutBoxStore) getMessage(ctx context.Context) (*model.OutBox, error) {
	conn := dbx.GetConnOrTx(ctx, o.db)
	builder := pgsql.Select(
		"id",
		"topic",
		"status",
		"msg",
		"created_at",
		"updated_at",
	).From(outboxTable).Where(squirrel.Eq{
		"status": []model.OutboxStatus{model.CreatedStatus},
	}).OrderBy("created_at ASC").Limit(1).Suffix("for update skip locked")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var dest model.OutBox
	if err := pgxscan.Get(ctx, conn, &dest, query, args...); err != nil {
		if dbx.IsErrNoRows(err) {
			return nil, nil
		}
		return nil, err
	}

	return &dest, nil
}

func (o *OutBoxStore) ProcessMessage(ctx context.Context, fn func(ctx context.Context, m *model.OutBox) error) (err error) {
	logger := logging.FromContext(ctx)
	tx, err := o.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	ctx = dbx.WithContext(ctx, tx)
	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); err != nil {
				logger.Errorf("failed to rollback: %s", errRollback)
			}
		} else {
			if errCommit := tx.Commit(ctx); err != nil {
				logger.Errorf("failed to commit: %s", errCommit)
			}
		}
	}()

	message, err := o.getMessage(ctx)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	if message == nil {
		return nil
	}

	if err = fn(ctx, message); err != nil {
		if markErr := o.markError(ctx, message.ID); markErr != nil {
			return fmt.Errorf("failed to mark error for message: %w", err)
		}
		return fmt.Errorf("failed to run fn for message: %w", err)
	}

	if err := o.markDone(ctx, message.ID); err != nil {
		return fmt.Errorf("failed to mark message done: %w", err)
	}

	return nil

}

func (o *OutBoxStore) markError(ctx context.Context, oid int64) error {
	tx := dbx.FromContext(ctx)
	builder := pgsql.Update(outboxTable).SetMap(map[string]interface{}{
		"status": model.ErrorStatus,
	}).Where(squirrel.Eq{
		"id": oid,
	})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (o *OutBoxStore) markDone(ctx context.Context, oid int64) error {
	tx := dbx.FromContext(ctx)
	builder := pgsql.Update(outboxTable).SetMap(map[string]interface{}{
		"status": model.DoneStatus,
	}).Where(squirrel.Eq{
		"id": oid,
	})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
