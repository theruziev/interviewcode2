package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/theruziev/oson_auth/internal/model"
	"github.com/theruziev/oson_auth/internal/pkg/dbx"
)

const contentTable = "contents"

var defaultContentFields = []string{
	"id",
	"public_id",
	"source",
	"preview_url",
	"tags",
	"description",
	"content_type",
	"status",
	"created_at",
	"updated_at",
	"views",
}

type ContentStore struct {
	db dbx.Querier
}

func NewContentStore(db dbx.Querier) *ContentStore {
	return &ContentStore{
		db: db,
	}
}

func (s *ContentStore) Add(ctx context.Context, content *model.Content) error {
	builder := pgsql.Insert(contentTable).SetMap(map[string]interface{}{
		"source":       content.Source,
		"preview_url":  content.PreviewURL,
		"preview_type": content.PreviewType,
		"description":  content.Description,
		"tags":         content.Tags,
		"content_type": content.Type,
		"status":       content.Status,
		"created_at":   content.CreatedAt,
		"updated_at":   content.UpdatedAt,
	}).Suffix("returning id")

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	err = pgxscan.Get(ctx, s.db, content, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *ContentStore) GetByID(ctx context.Context, cid uint64) (*model.Content, error) {
	builder := pgsql.Select(defaultContentFields...).
		From(contentTable).Where(sq.Eq{
		"id":     cid,
		"status": model.ContentStatusActive,
	})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var content model.Content
	err = pgxscan.Get(ctx, s.db, content, query, args...)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (s *ContentStore) ListNewest(ctx context.Context, lastID, limit uint64) (*model.ContentListResponse, error) {
	builder := pgsql.Select(defaultContentFields...).
		From(contentTable).Where(sq.Gt{
		"id":     lastID,
		"status": model.ContentStatusActive,
	}).Limit(limit).OrderBy("id desc")
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var contents []*model.Content
	err = pgxscan.Select(ctx, s.db, contents, query, args...)
	if err != nil {
		return nil, err
	}
	var newLastID uint64
	if len(contents) > 0 {
		newLastID = contents[len(contents)-1].ID
	}
	return &model.ContentListResponse{
		Result: contents,
		LastID: newLastID,
	}, nil
}
