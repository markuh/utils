package history

import (
	"fmt"

	"github.com/markuh/utils/pkg/apperrors"
)

const (
	ErrFmtSnapshot      = "can't snapshot %s"
	ErrFmtPack          = "can't pack %s history data"
	ErrFmtDiff          = "can't compute %s history diff"
	ErrFmtPrepareInsert = "can't prepare %s history insert"
	ErrFmtInsert        = "can't insert %s history revision"
	ErrFmtPrepareLoad   = "can't prepare %s history load"
	ErrFmtLoad          = "can't load %s history"
	ErrFmtReadRow       = "can't read %s history row"
)

func wrapEntity(err error, format, entity string) error {
	return apperrors.Wrap(err, fmt.Sprintf(format, entity))
}
