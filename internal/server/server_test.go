package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	New().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	res := rec.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Errorf("status: got %d, want %d", res.StatusCode, http.StatusOK)
	}
	if got, want := res.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("Content-Type: got %q, want %q", got, want)
	}
	if got, want := rec.Body.String(), "ok\n"; got != want {
		t.Errorf("body: got %q, want %q", got, want)
	}
}

func TestMetrics(t *testing.T) {
	rec := httptest.NewRecorder()
	New().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	res := rec.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want %d\nbody: %s", res.StatusCode, http.StatusOK, rec.Body.String())
	}
	if got, want := res.Header.Get("Content-Type"), "text/plain; version=0.0.4; charset=utf-8"; got != want {
		t.Errorf("Content-Type: got %q, want %q", got, want)
	}

	body := rec.Body.String()
	wantNames := []string{
		"nodeprobe_load1",
		"nodeprobe_load5",
		"nodeprobe_load15",
		"nodeprobe_memory_total_bytes",
		"nodeprobe_memory_free_bytes",
		"nodeprobe_memory_available_bytes",
		"nodeprobe_disk_total_bytes",
		"nodeprobe_disk_free_bytes",
		"nodeprobe_disk_available_bytes",
		"nodeprobe_disk_used_bytes",
	}
	for _, name := range wantNames {
		if !strings.Contains(body, "# TYPE "+name+" gauge\n") {
			t.Errorf("missing TYPE line for %s", name)
		}
	}
}

// TestMethodNotAllowed pins the behaviour the Go 1.22 method patterns
// give us for free: a known path with the wrong method is a 405, not a 404.
func TestMethodNotAllowed(t *testing.T) {
	rec := httptest.NewRecorder()
	New().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/healthz", nil))

	if got, want := rec.Code, http.StatusMethodNotAllowed; got != want {
		t.Errorf("status: got %d, want %d", got, want)
	}
}

func TestUnknownPath(t *testing.T) {
	rec := httptest.NewRecorder()
	New().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/nope", nil))

	if got, want := rec.Code, http.StatusNotFound; got != want {
		t.Errorf("status: got %d, want %d", got, want)
	}
}
