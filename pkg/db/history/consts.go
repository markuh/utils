package history

import "time"

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
