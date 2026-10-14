package json

import (
	"bytes"
	stdjson "encoding/json"
	"testing"
)

// Persisted events must remain readable and deterministic after changing codecs.
func TestPersistedPayloadCompatibility(t *testing.T) {
	values := []any{
		map[string]any{"z": "<>&\u2028\u2029", "a": "quote\"\n\\", "amount": 12.50},
		struct {
			ID      string     `json:"id"`
			Payload RawMessage `json:"payload"`
		}{ID: "event-1", Payload: RawMessage(`{"amount":12.5}`)},
		[]string{"zażółć", "", "東京"},
	}
	for _, value := range values {
		want, err := stdjson.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		got, err := Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("payload changed: got %s, want %s", got, want)
		}
		var decoded any
		if err := Unmarshal(want, &decoded); err != nil {
			t.Fatal(err)
		}
		var stream bytes.Buffer
		if err := NewEncoder(&stream).Encode(value); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(bytes.TrimSpace(stream.Bytes()), want) {
			t.Fatalf("stream payload changed: %s", stream.Bytes())
		}
	}
}
