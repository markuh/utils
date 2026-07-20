package history

import (
	"context"
	"encoding/json"

	"github.com/doug-martin/goqu/v9"
	"github.com/markuh/utils/pkg/apperrors"
	"github.com/markuh/utils/pkg/db"
	"github.com/markuh/utils/pkg/helpers"
)

func (s *Store[T]) OnSave(ctx context.Context, entity *T) error {
	if entity == nil {
		return apperrors.New("history: OnSave: nil entity")
	}

	entityID := s.entityID(entity)
	if entityID <= 0 {
		return apperrors.New("history: OnSave: invalid entity id")
	}

	snapshot, err := s.entitySnapshot(entity)
	if err != nil {
		return wrapEntity(err, ErrFmtSnapshot, s.entityLabel)
	}

	packed, err := helpers.PackData(snapshot)
	if err != nil {
		return wrapEntity(err, ErrFmtPack, s.entityLabel)
	}

	latest, err := s.loadLatestRevision(ctx, entityID)
	if err != nil {
		return err
	}

	nextRevision := 1
	var diff map[string]any
	if latest != nil {
		nextRevision = latest.revision + 1
		oldData, err := helpers.UnpackData(latest.data)
		if err != nil {
			return wrapEntity(err, ErrFmtPack, s.entityLabel)
		}
		diff, err = helpers.ComputeJSONDiff(oldData, snapshot)
		if err != nil {
			return wrapEntity(err, ErrFmtDiff, s.entityLabel)
		}
		if len(diff) == 0 {
			return nil
		}
	}

	return s.insertRevision(ctx, entityID, nextRevision, diff, packed, false)
}

func (s *Store[T]) insertRevision(
	ctx context.Context,
	entityID int64,
	revision int,
	diff map[string]any,
	packedData string,
	isDeleted bool,
) error {

	if diff == nil {
		diff = map[string]any{}
	}

	diffJSON, err := json.Marshal(diff)
	if err != nil {
		return wrapEntity(err, ErrFmtPrepareInsert, s.entityLabel)
	}

	query, args, err := db.PgDialect.
		Insert(s.table.Name).
		Rows(goqu.Record{
			colRevision:  revision,
			colEntityID:  entityID,
			colDiff:      diffJSON,
			colData:      packedData,
			colCreatedAt: goqu.L("NOW()"),
			colUpdatedAt: goqu.L("NOW()"),
			colIsDeleted: isDeleted,
		}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return wrapEntity(err, ErrFmtPrepareInsert, s.entityLabel)
	}

	if _, err := s.db(ctx).Exec(ctx, query, args...); err != nil {
		return wrapEntity(err, ErrFmtInsert, s.entityLabel)
	}
	return nil
}
