package uuidgen_test

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/mateusz-kacpura/integra-core/pkg/uuidgen"
)

func ExampleToPgtype() {
	id := uuid.MustParse("018f26a1-7c9b-7e44-8c8a-7e199f7d45e0")
	pgID := uuidgen.ToPgtype(id)
	restored := uuidgen.FromPgtype(pgID)

	fmt.Println("Restored UUID:", restored.String())
	// Output:
	// Restored UUID: 018f26a1-7c9b-7e44-8c8a-7e199f7d45e0
}

func ExampleMustParse() {
	id := uuidgen.MustParse("018f26a1-7c9b-7e44-8c8a-7e199f7d45e0")
	fmt.Println("Parsed Version:", id.Version())
	// Output:
	// Parsed Version: VERSION_7
}
