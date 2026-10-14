package response

import (
	"bytes"

	sjson "github.com/mateusz-kacpura/integra-core/pkg/json"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationErrors(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	errs := []ValidationError{
		{Field: "email", Message: "is required"},
		{Field: "name", Message: "too short"},
	}

	ValidationErrors(c, errs)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var resp Response
	err := sjson.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, CodeValidation, resp.Code)
	assert.Equal(t, "validation failed", resp.Message)
	assert.NotNil(t, resp.Data)
}

type testPayload struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func TestBindAndValidateSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"name":"Alice","email":"alice@example.com"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	var payload testPayload
	ok := BindAndValidate(c, &payload)

	assert.True(t, ok)
	assert.Equal(t, "Alice", payload.Name)
	assert.Equal(t, "alice@example.com", payload.Email)
}

func TestBindAndValidateInvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{not json`))
	c.Request.Header.Set("Content-Type", "application/json")

	var payload testPayload
	ok := BindAndValidate(c, &payload)

	assert.False(t, ok)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp Response
	err := sjson.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, CodeBadRequest, resp.Code)
}
