package pagination

import "strconv"

const (
	// GRPCDefaultPageSize is the default page size used for gRPC list requests.
	GRPCDefaultPageSize = 50
	// GRPCMaxPageSize is the upper bound on allowed page size for gRPC list requests.
	GRPCMaxPageSize = 200
)

// ParsePagination extracts pageSize and offset from raw values (proto-agnostic).
// Clamps pageSize to [1, GRPCMaxPageSize] and defaults to GRPCDefaultPageSize.
func ParsePagination(pageSize int32, pageToken string) (size, offset int) {
	size = int(pageSize)
	if size <= 0 {
		size = GRPCDefaultPageSize
	}
	if size > GRPCMaxPageSize {
		size = GRPCMaxPageSize
	}
	offset = TokenToOffset(pageToken)
	return
}

// TokenToOffset converts a page token string to an integer offset.
func TokenToOffset(token string) int {
	if token == "" {
		return 0
	}
	n, err := strconv.Atoi(token)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

// OffsetToToken converts an offset to a page token string.
func OffsetToToken(offset int) string {
	if offset <= 0 {
		return ""
	}
	return strconv.Itoa(offset)
}
