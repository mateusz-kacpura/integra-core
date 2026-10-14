package response

const (
	// CodeOK indicates successful request completion without errors.
	CodeOK = "OK"

	// CodeBadRequest indicates client request validation or malformed payload errors.
	CodeBadRequest = "BAD_REQUEST"

	// CodeUnauthorized indicates authentication failure, missing token, or revoked credentials.
	CodeUnauthorized = "UNAUTHORIZED"

	// CodeForbidden indicates insufficient authorization or permission denial.
	CodeForbidden = "FORBIDDEN"

	// CodeNotFound indicates that the requested resource does not exist.
	CodeNotFound = "NOT_FOUND"

	// CodeConflict indicates a unique constraint violation or concurrent modification conflict.
	CodeConflict = "CONFLICT"

	// CodeUnprocessable indicates semantically invalid entities or state transition errors.
	CodeUnprocessable = "UNPROCESSABLE_ENTITY"

	// CodeInternal indicates an unexpected internal server or execution error.
	CodeInternal = "INTERNAL_ERROR"

	// CodeServiceDown indicates dependency unavailability or temporary maintenance outage.
	CodeServiceDown = "SERVICE_UNAVAILABLE"

	// CodeTimeout indicates that operation processing exceeded its allocated deadline.
	CodeTimeout = "TIMEOUT"

	// CodeRateLimited indicates that client request frequency exceeded rate limits.
	CodeRateLimited = "RATE_LIMITED"

	// CodeValidation indicates that request parameters failed schema validation constraints.
	CodeValidation = "VALIDATION_ERROR"
)
