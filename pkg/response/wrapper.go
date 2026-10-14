package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sjson "github.com/mateusz-kacpura/integra-core/pkg/json"
)

// Response is the standard integra API response envelope
type Response struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta holds pagination metadata
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	TotalItems int64 `json:"total_items,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// OK sends a successful response
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    CodeOK,
		Data:    data,
	})
}

// OKWithMeta sends a successful response with pagination
func OKWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    CodeOK,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a 201 response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Code:    CodeOK,
		Data:    data,
	})
}

// NoContent sends a 204 response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code string, message string) {
	c.AbortWithStatusJSON(statusCode, Response{
		Success: false,
		Code:    code,
		Message: message,
	})
}

// BadRequest sends a 400 response
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, CodeBadRequest, message)
}

// Unauthorized sends a 401 response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

// Forbidden sends a 403 response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, CodeForbidden, message)
}

// NotFound sends a 404 response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, CodeNotFound, message)
}

// Conflict sends a 409 response
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, CodeConflict, message)
}

// InternalError sends a 500 response
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, CodeInternal, message)
}

// WriteJSON writes a Response envelope to any standard http.ResponseWriter.
func WriteJSON(w http.ResponseWriter, statusCode int, resp Response) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	data, err := sjson.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// WriteOK writes a successful Response envelope to any standard http.ResponseWriter.
func WriteOK(w http.ResponseWriter, data interface{}) error {
	return WriteJSON(w, http.StatusOK, Response{
		Success: true,
		Code:    CodeOK,
		Data:    data,
	})
}

// WriteError writes an error Response envelope to any standard http.ResponseWriter.
func WriteError(w http.ResponseWriter, statusCode int, code string, message string) error {
	return WriteJSON(w, statusCode, Response{
		Success: false,
		Code:    code,
		Message: message,
	})
}
