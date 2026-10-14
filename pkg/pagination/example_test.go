package pagination_test

import (
	"fmt"

	"github.com/mateusz-kacpura/integra-core/pkg/pagination"
)

func ExampleCursorEncode() {
	cursor := pagination.CursorEncode("order-987", 1700000000000)
	id, ts, err := pagination.CursorDecode(cursor)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decoded ID: %s, Timestamp: %d\n", id, ts)
	// Output:
	// Decoded ID: order-987, Timestamp: 1700000000000
}

func ExampleTotalPages() {
	fmt.Println("100 items with page size 20:", pagination.TotalPages(100, 20))
	fmt.Println("105 items with page size 20:", pagination.TotalPages(105, 20))
	// Output:
	// 100 items with page size 20: 5
	// 105 items with page size 20: 6
}
