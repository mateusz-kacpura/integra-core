package pagination

import (
	"testing"
)

func FuzzCursorDecode(f *testing.F) {
	// Seed with valid cursors
	f.Add(CursorEncode("abc", 123))
	f.Add(CursorEncode("", 0))

	// Seed with invalid inputs
	f.Add("")
	f.Add("!!!invalid!!!")
	f.Add("bm9waXBl") // "nopipe"
	f.Add("aWR8YWJj") // "id|abc"

	f.Fuzz(func(t *testing.T, cursor string) {
		// Must not panic regardless of input
		_, _, _ = CursorDecode(cursor)
	})
}
