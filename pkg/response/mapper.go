package response

import "net/http"

// HTTPStatusFromCode maps integra error codes to HTTP status codes
func HTTPStatusFromCode(code string) int {
	switch code {
	case CodeOK:
		return http.StatusOK
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeUnprocessable:
		return http.StatusUnprocessableEntity
	case CodeInternal:
		return http.StatusInternalServerError
	case CodeServiceDown:
		return http.StatusServiceUnavailable
	case CodeTimeout:
		return http.StatusGatewayTimeout
	case CodeRateLimited:
		return http.StatusTooManyRequests
	case CodeValidation:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
