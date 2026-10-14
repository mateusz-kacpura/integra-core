package uuidgen

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// New generates a new UUIDv7 (RFC 9562). Panics on failure (crypto/rand unavailable).
func New() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

// ToPgtype converts uuid.UUID to pgtype.UUID for use with pgx queries.
func ToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// ToPgtypePtr converts *uuid.UUID to pgtype.UUID. Returns invalid UUID if nil.
func ToPgtypePtr(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

// FromPgtype converts pgtype.UUID to uuid.UUID. Returns uuid.Nil if invalid.
func FromPgtype(p pgtype.UUID) uuid.UUID {
	if !p.Valid {
		return uuid.Nil
	}
	return p.Bytes
}

// Parse wraps uuid.Parse for convenience at API boundaries.
func Parse(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// MustParse wraps uuid.MustParse. Panics on invalid input.
func MustParse(s string) uuid.UUID {
	return uuid.MustParse(s)
}
