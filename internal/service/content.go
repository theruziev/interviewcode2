package service

import (
	"context"

	"github.com/theruziev/oson_auth/internal/db"
	"github.com/theruziev/oson_auth/internal/model"
	"github.com/theruziev/oson_auth/internal/pkg/dbx"
)

type ContentService struct {
	store *db.ContentStore
	pool  dbx.Querier // for transaction
}

func NewContentService(store *db.ContentStore, pool dbx.Querier) *ContentService {
	return &ContentService{
		store: store,
		pool:  pool,
	}
}

func (s *ContentService) Add(ctx context.Context, content *model.Content) error {
	err := s.store.Add(ctx, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *ContentService) ListNewest(ctx context.Context, lastID, limit uint64) (*model.ContentListResponse, error) {
	result, err := s.store.ListNewest(ctx, lastID, limit)
	if err != nil {
		return nil, err
	}
	return result, nil
}
