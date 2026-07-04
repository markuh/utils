package helpers

import "testing"

func TestComputeJSONDiff_identical(t *testing.T) {
	t.Parallel()

	old := map[string]any{"name": "a", "width": float64(1)}
	diff, err := ComputeJSONDiff(old, map[string]any{"name": "a", "width": float64(1)})
	if err != nil {
		t.Fatalf("ComputeJSONDiff: %v", err)
	}
	if len(diff) != 0 {
		t.Fatalf("expected empty diff, got %#v", diff)
	}
}

func TestComputeJSONDiff_addRemoveReplace(t *testing.T) {
	t.Parallel()

	old := map[string]any{"name": "a", "width": float64(1), "removed": true}
	newData := map[string]any{"name": "b", "width": float64(1), "added": float64(2)}

	diff, err := ComputeJSONDiff(old, newData)
	if err != nil {
		t.Fatalf("ComputeJSONDiff: %v", err)
	}

	if diff["name"] != "b" {
		t.Fatalf("name: %#v", diff["name"])
	}
	if diff["added"] != float64(2) {
		t.Fatalf("added: %#v", diff["added"])
	}
	if diff["removed"] != nil {
		t.Fatalf("removed: %#v", diff["removed"])
	}
}

func TestComputeJSONDiff_nestedAndList(t *testing.T) {
	t.Parallel()

	old := map[string]any{
		"data": map[string]any{
			"items": []any{float64(1), float64(2)},
		},
	}
	newData := map[string]any{
		"data": map[string]any{
			"items": []any{float64(1), float64(3), float64(4)},
		},
	}

	diff, err := ComputeJSONDiff(old, newData)
	if err != nil {
		t.Fatalf("ComputeJSONDiff: %v", err)
	}

	data, ok := diff["data"].(map[string]any)
	if !ok {
		t.Fatalf("data: %#v", diff["data"])
	}
	items, ok := data["items"].([]any)
	if !ok || len(items) != 3 {
		t.Fatalf("items: %#v", data["items"])
	}
}
