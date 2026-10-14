// Package uuidx re-exports uuidgen for backward compatibility.
// Consumers are encouraged to import "github.com/mateusz-kacpura/integra-core/pkg/uuidgen"
// directly for clear contextual domain naming.
package uuidx

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/mateusz-kacpura/integra-core/pkg/uuidgen"
)

// New generates a new UUIDv7.
func New() uuid.UUID {
	return uuidgen.New()
}

// ToPgtype converts uuid.UUID to pgtype.UUID.
func ToPgtype(id uuid.UUID) pgtype.UUID {
	return uuidgen.ToPgtype(id)
}

// ToPgtypePtr converts *uuid.UUID to pgtype.UUID.
func ToPgtypePtr(id *uuid.UUID) pgtype.UUID {
	return uuidgen.ToPgtypePtr(id)
}

// FromPgtype converts pgtype.UUID to uuid.UUID.
func FromPgtype(p pgtype.UUID) uuid.UUID {
	return uuidgen.FromPgtype(p)
}

// Parse wraps uuid.Parse.
func Parse(s string) (uuid.UUID, error) {
	return uuidgen.Parse(s)
}

// MustParse wraps uuid.MustParse.
func MustParse(s string) uuid.UUID {
	return uuidgen.MustParse(s)
}
