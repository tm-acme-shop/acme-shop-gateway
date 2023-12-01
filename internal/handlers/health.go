package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

type HealthHandler struct {
	startTime time.Time
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services,omitempty"`
}

type ReadinessResponse struct {
	Ready    bool              `json:"ready"`
	Checks   map[string]bool   `json:"checks"`
	Details  map[string]string `json:"details,omitempty"`
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    time.Since(h.startTime).String(),
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	resp := ReadinessResponse{
		Ready: true,
		Checks: map[string]bool{
			"gateway": true,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *HealthHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := map[string]interface{}{
		"uptime_seconds":    time.Since(h.startTime).Seconds(),
		"goroutines":        runtime.NumGoroutine(),
		"memory_alloc_mb":   float64(m.Alloc) / 1024 / 1024,
		"memory_sys_mb":     float64(m.Sys) / 1024 / 1024,
		"gc_runs":           m.NumGC,
		"go_version":        runtime.Version(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
