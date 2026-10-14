package response

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPStatusFromCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected int
	}{
		{"OK", CodeOK, http.StatusOK},
		{"BadRequest", CodeBadRequest, http.StatusBadRequest},
		{"Unauthorized", CodeUnauthorized, http.StatusUnauthorized},
		{"Forbidden", CodeForbidden, http.StatusForbidden},
		{"NotFound", CodeNotFound, http.StatusNotFound},
		{"Conflict", CodeConflict, http.StatusConflict},
		{"Unprocessable", CodeUnprocessable, http.StatusUnprocessableEntity},
		{"Internal", CodeInternal, http.StatusInternalServerError},
		{"ServiceDown", CodeServiceDown, http.StatusServiceUnavailable},
		{"Timeout", CodeTimeout, http.StatusGatewayTimeout},
		{"RateLimited", CodeRateLimited, http.StatusTooManyRequests},
		{"Validation", CodeValidation, http.StatusUnprocessableEntity},
		{"Unknown defaults to 500", "UNKNOWN_CODE", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, HTTPStatusFromCode(tt.code))
		})
	}
}
