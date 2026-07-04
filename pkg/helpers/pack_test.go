package helpers

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"reflect"
	"strings"
	"testing"
)

func TestPackData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    map[string]any
		wantErr bool
		errMsg  string
	}{
		{
			name: "empty map",
			data: map[string]any{},
		},
		{
			name: "simple values",
			data: map[string]any{
				"name":  "test",
				"count": float64(42),
				"ok":    true,
			},
		},
		{
			name: "nested map and slice",
			data: map[string]any{
				"data": map[string]any{
					"items": []any{float64(1), float64(2)},
				},
			},
		},
		{
			name:    "unsupported json type",
			data:    map[string]any{"ch": make(chan int)},
			wantErr: true,
			errMsg:  "can't pack data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := PackData(tt.data)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Fatalf("error %q, want substring %q", err.Error(), tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("PackData: %v", err)
			}
			if got == "" {
				t.Fatal("expected non-empty packed string")
			}

			unpacked, err := UnpackData(got)
			if err != nil {
				t.Fatalf("UnpackData round-trip: %v", err)
			}
			if !reflect.DeepEqual(unpacked, tt.data) {
				t.Fatalf("round-trip mismatch: got %#v, want %#v", unpacked, tt.data)
			}
		})
	}
}

func TestUnpackData(t *testing.T) {
	t.Parallel()

	packedEmpty, err := PackData(map[string]any{})
	if err != nil {
		t.Fatalf("PackData setup: %v", err)
	}

	packedSimple, err := PackData(map[string]any{"key": "value", "n": float64(1)})
	if err != nil {
		t.Fatalf("PackData setup: %v", err)
	}

	var invalidJSONBuf bytes.Buffer
	w, err := zlib.NewWriterLevel(&invalidJSONBuf, zlib.BestCompression)
	if err != nil {
		t.Fatalf("zlib setup: %v", err)
	}
	if _, err := w.Write([]byte("not json")); err != nil {
		t.Fatalf("zlib write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("zlib close: %v", err)
	}
	packedInvalidJSON := base64.StdEncoding.EncodeToString(invalidJSONBuf.Bytes())

	tests := []struct {
		name      string
		packed    string
		want      map[string]any
		wantErr   bool
		errSubstr string
	}{
		{
			name:   "empty map",
			packed: packedEmpty,
			want:   map[string]any{},
		},
		{
			name:   "simple map",
			packed: packedSimple,
			want:   map[string]any{"key": "value", "n": float64(1)},
		},
		{
			name:      "invalid base64",
			packed:    "!!!not-base64!!!",
			wantErr:   true,
			errSubstr: "can't unpack data",
		},
		{
			name:      "not zlib payload",
			packed:    base64.StdEncoding.EncodeToString([]byte("plain text")),
			wantErr:   true,
			errSubstr: "can't unpack data",
		},
		{
			name:      "invalid json after decompress",
			packed:    packedInvalidJSON,
			wantErr:   true,
			errSubstr: "can't unpack data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := UnpackData(tt.packed)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr) {
					t.Fatalf("error %q, want substring %q", err.Error(), tt.errSubstr)
				}
				return
			}
			if err != nil {
				t.Fatalf("UnpackData: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
