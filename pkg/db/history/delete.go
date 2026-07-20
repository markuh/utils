package history

import (
	"context"

	"github.com/markuh/utils/pkg/apperrors"
	"github.com/markuh/utils/pkg/helpers"
)

func (s *Store[T]) OnDelete(ctx context.Context, entity *T) error {
	if entity == nil {
		return apperrors.New("history: OnDelete: nil entity")
	}

	entityID := s.entityID(entity)
	if entityID <= 0 {
		return apperrors.New("history: OnDelete: invalid entity id")
	}

	snapshot, err := s.entitySnapshot(entity)
	if err != nil {
		return wrapEntity(err, ErrFmtSnapshot, s.entityLabel)
	}

	packed, err := helpers.PackData(snapshot)
	if err != nil {
		return wrapEntity(err, ErrFmtPack, s.entityLabel)
	}

	nextRevision := 1
	latest, err := s.loadLatestRevision(ctx, entityID)
	if err != nil {
		return err
	}
	if latest != nil {
		nextRevision = latest.revision + 1
	}

	return s.insertRevision(ctx, entityID, nextRevision, map[string]any{}, packed, true)
}
