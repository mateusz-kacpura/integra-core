package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	sjson "github.com/mateusz-kacpura/integra-core/pkg/json"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) Response {
	t.Helper()
	var resp Response
	err := sjson.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	return resp
}

func TestOK(t *testing.T) {
	c, w := setupTestContext()
	OK(c, map[string]string{"key": "value"})

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w)
	assert.True(t, resp.Success)
	assert.Equal(t, CodeOK, resp.Code)
	assert.NotNil(t, resp.Data)
}

func TestCreated(t *testing.T) {
	c, w := setupTestContext()
	Created(c, map[string]string{"id": "123"})

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := decodeResponse(t, w)
	assert.True(t, resp.Success)
	assert.Equal(t, CodeOK, resp.Code)
}

func TestNoContent(t *testing.T) {
	c, w := setupTestContext()
	NoContent(c)

	// gin's c.Status() sets the status on the writer; flush to recorder
	w.Flush()
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())
	assert.Empty(t, w.Body.Bytes())
}

func TestBadRequest(t *testing.T) {
	c, w := setupTestContext()
	BadRequest(c, "invalid input")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := decodeResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, CodeBadRequest, resp.Code)
	assert.Equal(t, "invalid input", resp.Message)
}

func TestUnauthorized(t *testing.T) {
	c, w := setupTestContext()
	Unauthorized(c, "missing token")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	resp := decodeResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, CodeUnauthorized, resp.Code)
}

func TestForbidden(t *testing.T) {
	c, w := setupTestContext()
	Forbidden(c, "no access")

	assert.Equal(t, http.StatusForbidden, w.Code)
	resp := decodeResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, CodeForbidden, resp.Code)
}

func TestNotFound(t *testing.T) {
	c, w := setupTestContext()
	NotFound(c, "resource not found")

	assert.Equal(t, http.StatusNotFound, w.Code)
	resp := decodeResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, CodeNotFound, resp.Code)
}

func TestInternalError(t *testing.T) {
	c, w := setupTestContext()
	InternalError(c, "something broke")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	resp := decodeResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, CodeInternal, resp.Code)
}

func TestConflict(t *testing.T) {
	c, w := setupTestContext()
	Conflict(c, "resource already exists")

	assert.Equal(t, http.StatusConflict, w.Code)
	resp := decodeResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, CodeConflict, resp.Code)
	assert.Equal(t, "resource already exists", resp.Message)
}

func TestOKWithMeta(t *testing.T) {
	c, w := setupTestContext()
	meta := &Meta{
		Page:       1,
		PageSize:   10,
		TotalItems: 100,
		TotalPages: 10,
	}
	OKWithMeta(c, []string{"a", "b"}, meta)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Meta)
	assert.Equal(t, 1, resp.Meta.Page)
	assert.Equal(t, 10, resp.Meta.PageSize)
	assert.Equal(t, int64(100), resp.Meta.TotalItems)
	assert.Equal(t, 10, resp.Meta.TotalPages)
}

func TestWriteOKStandardHTTP(t *testing.T) {
	rec := httptest.NewRecorder()
	err := WriteOK(rec, map[string]string{"foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))

	var resp Response
	err = sjson.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, CodeOK, resp.Code)
}

func TestWriteErrorStandardHTTP(t *testing.T) {
	rec := httptest.NewRecorder()
	err := WriteError(rec, http.StatusBadRequest, CodeBadRequest, "invalid payload")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp Response
	err = sjson.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, CodeBadRequest, resp.Code)
	assert.Equal(t, "invalid payload", resp.Message)
}
