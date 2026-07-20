package history

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"

	"github.com/markuh/utils/pkg/apperrors"
	"github.com/markuh/utils/pkg/db"
)

type latestRevision struct {
	revision int
	data     string
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

	out, err := db.ScanRows(rows, s.scanRevision)
	if err != nil {
		return nil, wrapEntity(err, ErrFmtReadRow, s.entityLabel)
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

	row, err := s.scanRevision(s.db(ctx).QueryRow(ctx, query, args...))
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

	row, err := s.scanRevision(s.db(ctx).QueryRow(ctx, query, args...))
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

	var rev int
	var data string
	if err := s.db(ctx).QueryRow(ctx, query, args...).Scan(&rev, &data); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, wrapEntity(err, ErrFmtLoad, s.entityLabel)
	}

	return &latestRevision{revision: rev, data: data}, nil
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

func (s *Store[T]) scanRevision(row db.PgxScanner) (*Revision, error) {
	var rev Revision
	if err := row.Scan(&rev.ID, &rev.Revision, &rev.EntityID, &rev.Diff, &rev.Data, &rev.CreatedAt,
		&rev.UpdatedAt, &rev.IsDeleted); err != nil {
		return nil, err
	}
	return &rev, nil
}
