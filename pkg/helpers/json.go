package helpers

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch"

	"github.com/markuh/utils/pkg/apperrors"
)

// ComputeJSONDiff build RFC7396 merge patch between two JSON objects.
// If documents are structurally equal, return empty map.
func ComputeJSONDiff(oldData, newData map[string]any) (map[string]any, error) {
	oldJSON, err := json.Marshal(oldData)
	if err != nil {
		return nil, apperrors.Wrap(err, "can't compute json diff")
	}

	newJSON, err := json.Marshal(newData)
	if err != nil {
		return nil, apperrors.Wrap(err, "can't compute json diff")
	}

	if jsonpatch.Equal(oldJSON, newJSON) {
		return map[string]any{}, nil
	}

	patch, err := jsonpatch.CreateMergePatch(oldJSON, newJSON)
	if err != nil {
		return nil, apperrors.Wrap(err, "can't compute json diff")
	}

	var out map[string]any
	if err := json.Unmarshal(patch, &out); err != nil {
		return nil, apperrors.Wrap(err, "can't compute json diff")
	}
	if out == nil {
		out = map[string]any{}
	}
	return out, nil
}

// JSONToMap serialize value to JSON and parse it as map (keys are json-tags of fields).
func JSONToMap(v any) (map[string]any, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, apperrors.Wrap(err, "can't serialize to json")
	}

	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, apperrors.Wrap(err, "can't parse json to map")
	}
	if out == nil {
		out = map[string]any{}
	}

	return out, nil
}
