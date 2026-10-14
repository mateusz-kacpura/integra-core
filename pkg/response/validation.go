package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidationError represents a single field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors sends a 422 response with field-level errors
func ValidationErrors(c *gin.Context, errors []ValidationError) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Code:    CodeValidation,
		Message: "validation failed",
		Data:    errors,
	})
}

// BindAndValidate binds JSON body and returns validation errors if any
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		BadRequest(c, fmt.Sprintf("invalid request body: %s", err.Error()))
		return false
	}
	return true
}
