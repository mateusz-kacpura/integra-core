// Package json re-exports jsoncodec for backward compatibility.
// Consumers are encouraged to import "github.com/mateusz-kacpura/integra-core/pkg/jsoncodec"
// directly to prevent collision with standard encoding/json.
package json

import "github.com/mateusz-kacpura/integra-core/pkg/jsoncodec"

// RawMessage re-exports jsoncodec.RawMessage.
type RawMessage = jsoncodec.RawMessage

// Number re-exports jsoncodec.Number.
type Number = jsoncodec.Number

var (
	// Marshal serializes v to JSON bytes.
	Marshal = jsoncodec.Marshal
	// Unmarshal deserializes JSON bytes into v.
	Unmarshal = jsoncodec.Unmarshal
	// MarshalIndent serializes v to indented JSON bytes.
	MarshalIndent = jsoncodec.MarshalIndent
	// Pretouch pre-compiles the JSON codec.
	Pretouch = jsoncodec.Pretouch
	// NewDecoder creates a JSON decoder.
	NewDecoder = jsoncodec.NewDecoder
	// NewEncoder creates a JSON encoder.
	NewEncoder = jsoncodec.NewEncoder
)
