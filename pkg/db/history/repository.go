package history

import (
	"context"

	"github.com/markuh/utils/pkg/apperrors"
	"github.com/markuh/utils/pkg/db"
)

const (
	colID        = "id"
	colRevision  = "revision"
	colEntityID  = "entity_id"
	colDiff      = "diff"
	colData      = "data"
	colCreatedAt = "created_at"
	colUpdatedAt = "updated_at"
	colIsDeleted = "is_deleted"
)

type Config[T any] struct {
	Table    db.TableConfig[T]
	EntityID func(*T) int64
}

type Store[T any] struct {
	entityLabel string
	table       db.TableConfig[T]
	entityID    func(*T) int64
	db          func(context.Context) db.IQuery
}

func NewStore[T any](db func(context.Context) db.IQuery, cfg Config[T]) (*Store[T], error) {
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	if db == nil {
		return nil, apperrors.New("history: nil db getter")
	}

	return &Store[T]{
		entityLabel: cfg.Table.Name,
		table:       cfg.Table,
		entityID:    cfg.EntityID,
		db:          db,
	}, nil
}

func validateConfig[T any](cfg *Config[T]) error {
	if cfg == nil {
		return apperrors.New("history: nil Config")
	}
	if cfg.Table.Name == "" {
		return apperrors.New("history: Config.Table.Name is empty")
	}
	if cfg.EntityID == nil {
		return apperrors.New("history: incomplete Config")
	}
	return nil
}
