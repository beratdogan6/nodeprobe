// Package server wires the nodeprobe HTTP routes to the collectors.
package server

import (
	"fmt"
	"net/http"

	"github.com/beratdogan6/nodeprobe/internal/collector"
	"github.com/beratdogan6/nodeprobe/internal/metrics"
)

// New returns the handler serving the nodeprobe endpoints:
// GET /healthz and GET /metrics.
func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.HandleFunc("GET /metrics", handleMetrics)
	return mux
}

// handleHealthz reports that the process is up. It touches no
// collectors, so it stays cheap enough for a liveness probe.
func handleHealthz(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "ok")
}

// handleMetrics collects load, memory and disk stats on every request
// and writes them in the Prometheus text format.
func handleMetrics(w http.ResponseWriter, r *http.Request) {
	load, err := collector.ReadLoadAvg()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mem, err := collector.ReadMemInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	disk, err := collector.ReadDisk("/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", metrics.ContentType)

	metrics.WriteGauge(w, "nodeprobe_load1", "1-minute load average.", load.Load1)
	metrics.WriteGauge(w, "nodeprobe_load5", "5-minute load average.", load.Load5)
	metrics.WriteGauge(w, "nodeprobe_load15", "15-minute load average.", load.Load15)

	metrics.WriteGauge(w, "nodeprobe_memory_total_bytes", "Total usable memory.", float64(mem.Total))
	metrics.WriteGauge(w, "nodeprobe_memory_free_bytes", "Memory not used at all.", float64(mem.Free))
	metrics.WriteGauge(w, "nodeprobe_memory_available_bytes", "Memory available for programs.", float64(mem.Available))

	metrics.WriteGauge(w, "nodeprobe_disk_total_bytes", "Total size of the root filesystem.", float64(disk.Total))
	metrics.WriteGauge(w, "nodeprobe_disk_free_bytes", "Free space, including the root reserve.", float64(disk.Free))
	metrics.WriteGauge(w, "nodeprobe_disk_available_bytes", "Space available to non-root users.", float64(disk.Available))
	metrics.WriteGauge(w, "nodeprobe_disk_used_bytes", "Used space on the root filesystem.", float64(disk.Used))
}
