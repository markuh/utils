package history

import "time"

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

type Revision struct {
	ID        int64
	Revision  int
	EntityID  int64
	Diff      []byte
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}
