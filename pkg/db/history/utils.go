package history

import (
	"github.com/markuh/utils/pkg/apperrors"
	"github.com/markuh/utils/pkg/helpers"
)

func (s *Store[T]) entitySnapshot(entity *T) (map[string]any, error) {
	if entity == nil {
		return nil, apperrors.New("history: nil entity")
	}
	return helpers.JSONToMap(entity)
}
