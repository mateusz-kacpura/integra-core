package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    Params
		expected Params
	}{
		{"zeros", Params{0, 0}, Params{1, DefaultPageSize}},
		{"negative and over max", Params{-5, 200}, Params{1, MaxPageSize}},
		{"valid values", Params{3, 50}, Params{3, 50}},
		{"page negative size ok", Params{-1, 10}, Params{1, 10}},
		{"size exactly 1", Params{1, 1}, Params{1, 1}},
		{"size exactly 100", Params{2, 100}, Params{2, 100}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.input
			p.Normalize()
			assert.Equal(t, tt.expected, p)
		})
	}
}

func TestOffset(t *testing.T) {
	tests := []struct {
		name     string
		params   Params
		expected int
	}{
		{"first page", Params{1, 20}, 0},
		{"third page size 10", Params{3, 10}, 20},
		{"second page size 50", Params{2, 50}, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.params.Offset())
		})
	}
}

func TestTotalPages(t *testing.T) {
	tests := []struct {
		name       string
		totalItems int64
		pageSize   int
		expected   int
	}{
		{"zero items", 0, 20, 0},
		{"21 items size 20", 21, 20, 2},
		{"exact fit", 20, 20, 1},
		{"one item", 1, 20, 1},
		{"zero page size", 10, 0, 0},
		{"large count", 1000, 7, 143},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, TotalPages(tt.totalItems, tt.pageSize))
		})
	}
}

func TestCursorEncodeDecode(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		timestamp int64
	}{
		{"simple", "abc123", 1700000000},
		{"uuid-like", "550e8400-e29b-41d4-a716-446655440000", 0},
		{"empty id", "", 999},
		{"negative timestamp", "id", -1},
		{"large timestamp", "id", 9999999999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := CursorEncode(tt.id, tt.timestamp)
			assert.NotEmpty(t, encoded)

			id, ts, err := CursorDecode(encoded)
			require.NoError(t, err)
			assert.Equal(t, tt.id, id)
			assert.Equal(t, tt.timestamp, ts)
		})
	}
}

func TestCursorDecodeInvalidBase64(t *testing.T) {
	_, _, err := CursorDecode("!!!not-base64!!!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid cursor")
}

func TestCursorDecodeMalformed(t *testing.T) {
	// Encode a string without a pipe separator
	_, _, err := CursorDecode("bm9waXBl") // base64 of "nopipe"
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed cursor")
}

func TestCursorDecodeInvalidTimestamp(t *testing.T) {
	// base64 of "id|abc"
	_, _, err := CursorDecode("aWR8YWJj") // base64url of "id|abc"
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid cursor timestamp")
}
