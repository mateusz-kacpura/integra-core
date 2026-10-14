// Package jsoncodec provides high-performance JSON encoding and decoding using sonic,
// eliminating package name collisions with the standard library encoding/json.
package jsoncodec

import (
	"encoding/json"
	"io"

	"github.com/bytedance/sonic"
)

// RawMessage re-exports encoding/json.RawMessage so consumers don't need both imports.
type RawMessage = json.RawMessage

// Number re-exports encoding/json.Number so consumers don't need both imports.
type Number = json.Number

var (
	// Marshal serializes v to JSON bytes using sonic.
	Marshal = sonic.ConfigStd.Marshal
	// Unmarshal deserializes JSON bytes into v using sonic.
	Unmarshal = sonic.ConfigStd.Unmarshal
	// MarshalIndent serializes v to indented JSON bytes using sonic.
	MarshalIndent = sonic.ConfigStd.MarshalIndent
	// Pretouch pre-compiles the JSON codec for the given reflect.Type at startup.
	// Call in init() for heavy domain structs to avoid lazy compilation overhead.
	Pretouch = sonic.Pretouch
)

// NewDecoder creates a sonic JSON decoder.
func NewDecoder(r io.Reader) sonic.Decoder {
	return sonic.ConfigStd.NewDecoder(r)
}

// NewEncoder creates a sonic JSON encoder.
func NewEncoder(w io.Writer) sonic.Encoder {
	return sonic.ConfigStd.NewEncoder(w)
}
