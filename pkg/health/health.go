package health

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Checker is a function that checks a dependency's health
type Checker func(ctx context.Context) error

// Handler manages health check endpoints
type Handler struct {
	mu       sync.RWMutex
	checkers map[string]Checker
}

// New creates a new health handler
func New() *Handler {
	return &Handler{
		checkers: make(map[string]Checker),
	}
}

// Register adds a named health checker
func (h *Handler) Register(name string, checker Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers[name] = checker
}

// LivenessHandler returns a simple liveness probe (always 200 if process is running)
func (h *Handler) LivenessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	}
}

// ReadinessHandler checks all registered dependencies
func (h *Handler) ReadinessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.mu.RLock()
		checkers := make(map[string]Checker, len(h.checkers))
		for k, v := range h.checkers {
			checkers[k] = v
		}
		h.mu.RUnlock()

		// Execute checks WITHOUT holding lock
		results := make(map[string]string)
		allHealthy := true

		for name, checker := range checkers {
			if err := checker(c.Request.Context()); err != nil {
				results[name] = err.Error()
				allHealthy = false
			} else {
				results[name] = "ok"
			}
		}

		status := http.StatusOK
		if !allHealthy {
			status = http.StatusServiceUnavailable
		}

		c.JSON(status, gin.H{
			"status": map[bool]string{true: "ready", false: "not_ready"}[allHealthy],
			"checks": results,
		})
	}
}
