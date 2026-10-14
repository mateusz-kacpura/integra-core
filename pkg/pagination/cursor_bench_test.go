package pagination

import (
	"testing"
)

func BenchmarkCursorEncode(b *testing.B) {
	for b.Loop() {
		CursorEncode("550e8400-e29b-41d4-a716-446655440000", 1700000000)
	}
}

func BenchmarkCursorDecode(b *testing.B) {
	encoded := CursorEncode("550e8400-e29b-41d4-a716-446655440000", 1700000000)
	for b.Loop() {
		_, _, err := CursorDecode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}
