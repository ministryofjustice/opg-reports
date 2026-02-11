package headers

import (
	"testing"
)

func TestTabulateHeaders(t *testing.T) {

	headers := &Headers{Headers: []*Header{
		{Field: "name", Type: KEY, Default: ""},
		{Field: "team", Type: KEY, Default: ""},
		{Field: "2026-01", Type: DATA, Default: 0.0},
		{Field: "2025-12", Type: DATA, Default: 0.0},
		{Field: "trend", Type: EXTRA, Default: ""},
		{Field: "total", Type: END, Default: 0.0},
	}}
	// test keys
	keys := headers.Keys()
	if len(keys) != 2 {
		t.Errorf("keys listing failed ")
	}
}
