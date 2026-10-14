package health

import (
	"context"

	sjson "github.com/mateusz-kacpura/integra-core/pkg/json"

	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func serveHandler(handler gin.HandlerFunc) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(w)
	engine.GET("/test", handler)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	engine.ServeHTTP(w, req)
	return w
}

func TestLivenessHandler(t *testing.T) {
	h := New()
	w := serveHandler(h.LivenessHandler())

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	err := sjson.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "alive", body["status"])
}

func TestReadinessAllHealthy(t *testing.T) {
	h := New()
	h.Register("db", func(ctx context.Context) error { return nil })
	h.Register("cache", func(ctx context.Context) error { return nil })

	w := serveHandler(h.ReadinessHandler())

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	err := sjson.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "ready", body["status"])

	checks := body["checks"].(map[string]interface{})
	assert.Equal(t, "ok", checks["db"])
	assert.Equal(t, "ok", checks["cache"])
}

func TestReadinessOneUnhealthy(t *testing.T) {
	h := New()
	h.Register("db", func(ctx context.Context) error { return nil })
	h.Register("cache", func(ctx context.Context) error { return errors.New("connection refused") })

	w := serveHandler(h.ReadinessHandler())

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var body map[string]interface{}
	err := sjson.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "not_ready", body["status"])

	checks := body["checks"].(map[string]interface{})
	assert.Equal(t, "ok", checks["db"])
	assert.Equal(t, "connection refused", checks["cache"])
}

func TestReadinessNoCheckers(t *testing.T) {
	h := New()

	w := serveHandler(h.ReadinessHandler())

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	err := sjson.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "ready", body["status"])
}

func TestRegisterMultiple(t *testing.T) {
	h := New()
	h.Register("db", func(ctx context.Context) error { return nil })
	h.Register("cache", func(ctx context.Context) error { return nil })
	h.Register("queue", func(ctx context.Context) error { return nil })

	w := serveHandler(h.ReadinessHandler())

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	err := sjson.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)

	checks := body["checks"].(map[string]interface{})
	assert.Len(t, checks, 3)
	assert.Equal(t, "ok", checks["db"])
	assert.Equal(t, "ok", checks["cache"])
	assert.Equal(t, "ok", checks["queue"])
}
