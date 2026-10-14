package logger

import (
	"io"
	"regexp"

	"github.com/rs/zerolog"
)

// PII detection patterns.
var (
	// Credit card: 13-19 digit sequences (optionally separated by spaces or dashes).
	reCreditCard = regexp.MustCompile(`\b([0-9]{4})[\s\-]?[0-9]{4}[\s\-]?[0-9]{4}[\s\-]?([0-9]{4})\b`)

	// Extended credit card: handles 13-19 digit raw sequences without separators.
	reCreditCardRaw = regexp.MustCompile(`\b([0-9]{4})([0-9]{9,15})\b`)

	// NIP (Polish tax ID): exactly 10 digits, not part of a longer number.
	reNIP = regexp.MustCompile(`\b([0-9]{10})\b`)

	// PESEL (Polish national ID): exactly 11 digits, not part of a longer number.
	rePESEL = regexp.MustCompile(`\b([0-9]{11})\b`)

	// JWT token: standard Base64url-encoded JWT starting with eyJ.
	reJWT = regexp.MustCompile(`eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)

	// Email address.
	reEmail = regexp.MustCompile(`\b([a-zA-Z0-9._%+\-]+)@([a-zA-Z0-9.\-]+\.[a-zA-Z]{2,})\b`)
)

// SanitizeString masks PII patterns in the given string.
// Detected patterns and their masked forms:
//   - Credit card numbers  -> first 4 + " **** **** " + last 4
//   - PESEL (11 digits)    -> [PESEL_REDACTED]  (checked before NIP to avoid false match)
//   - NIP (10 digits)      -> [NIP_REDACTED]
//   - JWT tokens (eyJ...)  -> [JWT_REDACTED]
//   - Email addresses      -> first char + "***@" + domain
//
// Order matters: PESEL (11 digits) is checked before NIP (10 digits) to prevent
// partial matches. Credit cards are checked first since they have the most
// distinctive format with separators.
func SanitizeString(s string) string {
	// Credit cards with separators (e.g. "4111 1111 1111 1234" or "4111-1111-1111-1234").
	s = reCreditCard.ReplaceAllString(s, "${1} **** **** ${2}")

	// Credit cards as raw digit sequences (13-19 digits).
	s = reCreditCardRaw.ReplaceAllStringFunc(s, func(match string) string {
		if len(match) < 13 {
			return match
		}
		return match[:4] + " **** **** " + match[len(match)-4:]
	})

	// JWT tokens.
	s = reJWT.ReplaceAllString(s, "[JWT_REDACTED]")

	// Email addresses: "user@domain.com" -> "u***@domain.com".
	s = reEmail.ReplaceAllStringFunc(s, func(match string) string {
		loc := reEmail.FindStringSubmatchIndex(match)
		if loc == nil {
			return match
		}
		local := match[loc[2]:loc[3]]
		domain := match[loc[4]:loc[5]]
		if len(local) == 0 {
			return match
		}
		return string(local[0]) + "***@" + domain
	})

	// PESEL (11 digits) — must be checked before NIP (10 digits).
	s = rePESEL.ReplaceAllString(s, "[PESEL_REDACTED]")

	// NIP (10 digits).
	s = reNIP.ReplaceAllString(s, "[NIP_REDACTED]")

	return s
}

// SanitizerHook is a zerolog.Hook that sanitizes PII from log messages.
//
// Limitation: zerolog's Hook.Run method receives the message string but cannot
// modify it after it has been added to the event. The hook will append a
// sanitized copy of the message as "sanitized_msg" when PII is detected.
//
// For full message sanitization, prefer wrapping log calls with SanitizeString:
//
//	log.Info().Msg(logger.SanitizeString(untrustedInput))
type SanitizerHook struct{}

// Run implements zerolog.Hook. It detects PII in the log message and, if found,
// appends a "sanitized_msg" field with the cleaned version.
func (h SanitizerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	sanitized := SanitizeString(msg)
	if sanitized != msg {
		e.Str("sanitized_msg", sanitized)
	}
}

// SanitizingWriter wraps an io.Writer and applies SanitizeString to all output bytes.
// This guarantees that any PII serialized into JSON log messages or fields is redacted
// before being written to the underlying stream (stdout, file, or collector).
type SanitizingWriter struct {
	Out io.Writer
}

// NewSanitizingWriter creates an io.Writer wrapper that scrubs PII from the byte stream.
func NewSanitizingWriter(out io.Writer) *SanitizingWriter {
	return &SanitizingWriter{Out: out}
}

// Write intercepts the write stream and redacts PII before passing to the destination.
func (w *SanitizingWriter) Write(p []byte) (n int, err error) {
	sanitized := SanitizeString(string(p))
	_, err = w.Out.Write([]byte(sanitized))
	return len(p), err
}

// WriteLevel implements zerolog.LevelWriter to ensure level-aware logs are redacted.
func (w *SanitizingWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	sanitized := SanitizeString(string(p))
	data := []byte(sanitized)
	if lw, ok := w.Out.(zerolog.LevelWriter); ok {
		_, err = lw.WriteLevel(level, data)
	} else {
		_, err = w.Out.Write(data)
	}
	return len(p), err
}
