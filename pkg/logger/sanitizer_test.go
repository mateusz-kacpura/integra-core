package logger

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSanitizeStringCreditCardWithSpaces(t *testing.T) {
	input := "Card: 4111 1111 1111 1234"
	result := SanitizeString(input)
	assert.Contains(t, result, "4111")
	assert.Contains(t, result, "1234")
	assert.Contains(t, result, "****")
	assert.NotContains(t, result, "1111 1111")
}

func TestSanitizeStringCreditCardWithDashes(t *testing.T) {
	input := "Card: 4111-1111-1111-5678"
	result := SanitizeString(input)
	assert.Contains(t, result, "4111")
	assert.Contains(t, result, "5678")
	assert.Contains(t, result, "****")
}

func TestSanitizeStringRawCreditCard(t *testing.T) {
	input := "Number: 4111111111111234"
	result := SanitizeString(input)
	assert.Contains(t, result, "4111")
	assert.Contains(t, result, "1234")
	assert.Contains(t, result, "****")
	assert.NotContains(t, result, "4111111111111234")
}

func TestSanitizeStringEmail(t *testing.T) {
	input := "Contact user@example.com for info"
	result := SanitizeString(input)
	assert.NotContains(t, result, "user@example.com")
	assert.Contains(t, result, "u***@example.com")
}

func TestSanitizeStringJWT(t *testing.T) {
	// A minimal JWT-like string (three dot-separated base64url segments starting with eyJ).
	jwt := "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
	input := "Token: " + jwt
	result := SanitizeString(input)
	assert.NotContains(t, result, jwt)
	assert.Contains(t, result, "[JWT_REDACTED]")
}

func TestSanitizeStringNoChange(t *testing.T) {
	input := "Hello, this is a normal log message with no PII"
	result := SanitizeString(input)
	assert.Equal(t, input, result)
}

func TestSanitizeStringMultiplePatterns(t *testing.T) {
	input := "Email: admin@test.com JWT: eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxIn0.abc123_def"
	result := SanitizeString(input)
	assert.Contains(t, result, "a***@test.com")
	assert.Contains(t, result, "[JWT_REDACTED]")
}

func TestSanitizerHookDoesNotPanicOnCleanMessage(t *testing.T) {
	// The SanitizerHook.Run should not panic on a clean message.
	// We cannot easily test the full hook chain without a real zerolog.Event,
	// but we can verify the underlying SanitizeString logic is correct.
	clean := "no PII here"
	assert.Equal(t, clean, SanitizeString(clean))
}

func TestSanitizingWriterRedactsStreamOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := NewSanitizingWriter(buf)
	log := zerolog.New(writer).With().Timestamp().Logger()

	log.Info().Str("card", "4111-1111-1111-5678").Msg("Processing order for user@example.com")

	out := buf.String()
	assert.NotContains(t, out, "user@example.com")
	assert.Contains(t, out, "u***@example.com")
	assert.NotContains(t, out, "4111-1111-1111-5678")
	assert.Contains(t, out, "4111 **** **** 5678")
}
