package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/tgdrive/teldrive/internal/version"
)

type HealthResponse struct {
	Status    string           `json:"status"`
	Version   string           `json:"version"`
	Commit    string           `json:"commit"`
	BuildTime string           `json:"build_time"`
	Uptime    time.Duration    `json:"uptime"`
	GoVersion string           `json:"go_version"`
	Memory    runtime.MemStats `json:"memory"`
}

var startTime = time.Now()

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	v := version.Get()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	resp := HealthResponse{
		Status:    "healthy",
		Version:   v.Version,
		Commit:    v.Commit,
		BuildTime: v.BuildTime,
		Uptime:    time.Since(startTime),
		GoVersion: runtime.Version(),
		Memory:    memStats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
