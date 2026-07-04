package helpers

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/markuh/utils/pkg/apperrors"
)

// PackData serialize map to JSON, compress with zlib (level 9) and encode to base64.
func PackData(data map[string]any) (string, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return "", apperrors.Wrap(err, "can't pack data")
	}

	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return "", apperrors.Wrap(err, "can't pack data")
	}
	if _, err := w.Write(raw); err != nil {
		return "", apperrors.Wrap(err, "can't pack data")
	}
	if err := w.Close(); err != nil {
		return "", apperrors.Wrap(err, "can't pack data")
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// UnpackData decode string created by PackData.
func UnpackData(packed string) (map[string]any, error) {
	raw, err := base64.StdEncoding.DecodeString(packed)
	if err != nil {
		return nil, apperrors.Wrap(err, "can't unpack data")
	}

	r, err := zlib.NewReader(bytes.NewReader(raw))
	if err != nil {
		return nil, apperrors.Wrap(err, "can't unpack data")
	}
	defer r.Close()

	decoded, err := io.ReadAll(r)
	if err != nil {
		return nil, apperrors.Wrap(err, "can't unpack data")
	}

	var out map[string]any
	if err := json.Unmarshal(decoded, &out); err != nil {
		return nil, apperrors.Wrap(err, "can't unpack data")
	}
	if out == nil {
		out = map[string]any{}
	}
	return out, nil
}
