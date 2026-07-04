package history

import (
	"context"
	"fmt"

	"github.com/markuh/utils/pkg/db"
)

type Config[T any] struct {
	Table    db.TableConfig
	EntityID func(*T) int64
}

type Store[T any] struct {
	entityLabel string
	table       db.TableConfig
	entityID    func(*T) int64
	db          func(context.Context) db.IQuery
}

func NewStore[T any](db func(context.Context) db.IQuery, cfg Config[T]) (*Store[T], error) {
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("history: nil db getter")
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
		return fmt.Errorf("history: nil Config")
	}
	if cfg.Table.Name == "" {
		return fmt.Errorf("history: Config.Table.Name is empty")
	}
	if cfg.EntityID == nil {
		return fmt.Errorf("history: incomplete Config")
	}
	return nil
}

func tableName(cfg db.TableConfig) (string, error) {
	if cfg.Name == "" {
		return "", fmt.Errorf("history: table name is empty")
	}
	return cfg.Name, nil
}
