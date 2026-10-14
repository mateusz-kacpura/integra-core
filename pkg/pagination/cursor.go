package pagination

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

// DefaultPageSize is the default number of items per page for offset-based pagination.
const DefaultPageSize = 20

// MaxPageSize is the maximum allowed items per page for offset-based pagination.
const MaxPageSize = 100

// Params represents pagination parameters
type Params struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// Normalize ensures pagination params are within valid bounds
func (p *Params) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
}

// Offset returns the SQL offset for the current page
func (p *Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// TotalPages calculates total pages from total items count
func TotalPages(totalItems int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	pages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		pages++
	}
	return pages
}

func parseTimestamp(b []byte) (int64, error) {
	if len(b) == 0 {
		return 0, fmt.Errorf("empty timestamp")
	}
	neg := false
	if b[0] == '-' {
		neg = true
		b = b[1:]
		if len(b) == 0 {
			return 0, fmt.Errorf("invalid timestamp")
		}
	}
	var n int64
	for _, c := range b {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid digit: %c", c)
		}
		d := int64(c - '0')
		if n > (1<<63-1-d)/10 {
			return 0, fmt.Errorf("timestamp overflow")
		}
		n = n*10 + d
	}
	if neg {
		return -n, nil
	}
	return n, nil
}

// CursorEncode encodes a cursor from id and timestamp.
// Uses a stack-allocated buffer for the common case (UUID id) — 1 allocation (return string).
func CursorEncode(id string, timestamp int64) string {
	needed := len(id) + 1 + 21
	if needed <= 64 {
		var rawBuf [64]byte
		raw := append(rawBuf[:0], id...)
		raw = append(raw, '|')
		raw = strconv.AppendInt(raw, timestamp, 10)

		var encBuf [96]byte
		base64.URLEncoding.Encode(encBuf[:], raw)
		return string(encBuf[:base64.URLEncoding.EncodedLen(len(raw))])
	}
	raw := make([]byte, 0, needed)
	raw = append(raw, id...)
	raw = append(raw, '|')
	raw = strconv.AppendInt(raw, timestamp, 10)
	return base64.URLEncoding.EncodeToString(raw)
}

// CursorDecode decodes a cursor into id and timestamp.
// Avoids strings.SplitN and extra string allocations by decoding into stack buffer and parsing integer directly.
func CursorDecode(cursor string) (id string, timestamp int64, err error) {
	if len(cursor) <= 128 {
		var rawBuf [128]byte
		var srcBuf [128]byte
		copy(srcBuf[:], cursor)
		n, decErr := base64.URLEncoding.Decode(rawBuf[:], srcBuf[:len(cursor)])
		if decErr != nil {
			return "", 0, fmt.Errorf("invalid cursor: %w", decErr)
		}
		raw := rawBuf[:n]
		pipe := -1
		for i, b := range raw {
			if b == '|' {
				pipe = i
				break
			}
		}
		if pipe < 0 {
			return "", 0, fmt.Errorf("malformed cursor")
		}
		ts, tsErr := parseTimestamp(raw[pipe+1:])
		if tsErr != nil {
			return "", 0, fmt.Errorf("invalid cursor timestamp: %w", tsErr)
		}
		return string(raw[:pipe]), ts, nil
	}

	raw, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return "", 0, fmt.Errorf("invalid cursor: %w", err)
	}

	pipe := -1
	for i, b := range raw {
		if b == '|' {
			pipe = i
			break
		}
	}
	if pipe < 0 {
		return "", 0, fmt.Errorf("malformed cursor")
	}

	ts, tsErr := parseTimestamp(raw[pipe+1:])
	if tsErr != nil {
		return "", 0, fmt.Errorf("invalid cursor timestamp: %w", tsErr)
	}

	return string(raw[:pipe]), ts, nil
}

// CursorWhere returns SQL WHERE clause parameters for cursor-based pagination.
// Returns the decoded timestamp and ID for use in: WHERE (created_at, id) < ($ts, $id) ORDER BY created_at DESC, id DESC LIMIT $limit
// If cursor is empty, returns zero values (meaning "start from beginning").
func CursorWhere(cursor string) (timestamp time.Time, id string, hasCursor bool, err error) {
	if cursor == "" {
		return time.Time{}, "", false, nil
	}

	decodedID, ts, err := CursorDecode(cursor)
	if err != nil {
		return time.Time{}, "", false, err
	}

	return time.Unix(0, ts*int64(time.Millisecond)).UTC(), decodedID, true, nil
}

// CursorResponse is a generic envelope for cursor-paginated REST responses.
type CursorResponse struct {
	Items      any    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}
