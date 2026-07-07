package history

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"

	"github.com/markuh/utils/pkg/apperrors"
	"github.com/markuh/utils/pkg/db"
	"github.com/markuh/utils/pkg/helpers"
)

type latestRevision struct {
	revision int
	data     string
}

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

func (s *Store[T]) LoadByEntityID(ctx context.Context, entityID int64) ([]*Revision, error) {
	if entityID <= 0 {
		return nil, apperrors.New("history: LoadByEntityID: invalid entity id")
	}

	querySQL := db.PgDialect.From(s.table.Name).
		Select(s.selectFields()...).
		Where(goqu.Ex{colEntityID: entityID}).
		Order(goqu.I(colRevision).Desc())

	query, args, err := querySQL.Prepared(true).ToSQL()
	if err != nil {
		return nil, wrapEntity(err, ErrFmtPrepareLoad, s.entityLabel)
	}

	rows, err := s.db(ctx).Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, wrapEntity(err, ErrFmtLoad, s.entityLabel)
	}
	defer rows.Close()

	var out []*Revision
	for rows.Next() {
		row, err := s.readRevisionRow(rows)
		if err != nil {
			return nil, wrapEntity(err, ErrFmtReadRow, s.entityLabel)
		}
		if row == nil {
			continue
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, wrapEntity(err, ErrFmtLoad, s.entityLabel)
	}
	return out, nil
}

func (s *Store[T]) GetLatest(ctx context.Context, entityID int64) (*Revision, error) {
	if entityID <= 0 {
		return nil, apperrors.New("history: GetLatest: invalid entity id")
	}

	querySQL := db.PgDialect.From(s.table.Name).
		Select(s.selectFields()...).
		Where(goqu.Ex{colEntityID: entityID}).
		Order(goqu.I(colRevision).Desc()).
		Limit(1)

	query, args, err := querySQL.Prepared(true).ToSQL()
	if err != nil {
		return nil, wrapEntity(err, ErrFmtPrepareLoad, s.entityLabel)
	}

	row, err := s.readRevisionRow(s.db(ctx).QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, wrapEntity(err, ErrFmtReadRow, s.entityLabel)
	}

	return row, nil
}

func (s *Store[T]) GetByID(ctx context.Context, revisionID int64) (*Revision, error) {
	if revisionID <= 0 {
		return nil, apperrors.New("history: GetByID: invalid revision id")
	}

	querySQL := db.PgDialect.From(s.table.Name).
		Select(s.selectFields()...).
		Where(goqu.Ex{colID: revisionID}).
		Limit(1)

	query, args, err := querySQL.Prepared(true).ToSQL()
	if err != nil {
		return nil, wrapEntity(err, ErrFmtPrepareLoad, s.entityLabel)
	}

	row, err := s.readRevisionRow(s.db(ctx).QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, wrapEntity(err, ErrFmtReadRow, s.entityLabel)
	}

	return row, nil
}

func (s *Store[T]) loadLatestRevision(ctx context.Context, entityID int64) (*latestRevision, error) {
	query, args, err := db.PgDialect.
		From(s.table.Name).
		Select(colRevision, colData).
		Where(goqu.Ex{colEntityID: entityID}).
		Order(goqu.I(colRevision).Desc()).
		Limit(1).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, wrapEntity(err, ErrFmtPrepareLoad, s.entityLabel)
	}

	row := s.db(ctx).QueryRow(ctx, query, args...)
	var rev int
	var data string
	if err := row.Scan(&rev, &data); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, wrapEntity(err, ErrFmtLoad, s.entityLabel)
	}

	return &latestRevision{revision: rev, data: data}, nil
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

func (s *Store[T]) selectFields() []any {
	if len(s.table.SelectFields) > 0 {
		return s.table.SelectFields
	}
	return []any{
		colID,
		colRevision,
		colEntityID,
		colDiff,
		colData,
		colCreatedAt,
		colUpdatedAt,
		colIsDeleted,
	}
}

func (s *Store[T]) entitySnapshot(entity *T) (map[string]any, error) {
	if entity == nil {
		return nil, apperrors.New("history: nil entity")
	}
	return helpers.JSONToMap(entity)
}

func (s *Store[T]) readRevisionRow(row pgx.Row) (*Revision, error) {
	var rev Revision
	if err := row.Scan(&rev.ID, &rev.Revision, &rev.EntityID, &rev.Diff, &rev.Data, &rev.CreatedAt,
		&rev.UpdatedAt, &rev.IsDeleted); err != nil {
		return nil, err
	}
	return &rev, nil
}
