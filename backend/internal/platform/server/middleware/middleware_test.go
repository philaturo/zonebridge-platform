package middleware

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func TestRequestID(t *testing.T) {
	t.Parallel()

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID, ok := r.Context().Value(requestIDKey).(string)
		if !ok || reqID == "" {
			t.Error("expected request ID in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	if rr.Header().Get("X-Request-Id") == "" {
		t.Error("expected X-Request-Id header to be set")
	}
}

func TestRealIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		headers  map[string]string
		remote   string
		expected string
	}{
		{"X-Forwarded-For", map[string]string{"X-Forwarded-For": "203.0.113.1"}, "10.0.0.1", "203.0.113.1"},
		{"X-Real-IP", map[string]string{"X-Real-IP": "198.51.100.1"}, "10.0.0.1", "198.51.100.1"},
		{"Fallback to RemoteAddr", map[string]string{}, "10.0.0.1:1234", "10.0.0.1:1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			handler := RealIP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ip, _ := r.Context().Value(realIPKey).(string)
				if ip != tt.expected {
					t.Errorf("expected IP %q, got %q", tt.expected, ip)
				}
			}))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tt.remote
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			handler.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}

func TestRecoverer(t *testing.T) {
	t.Parallel()

	logger := slog.Default()
	handler := Recoverer(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rr := httptest.NewRecorder()

	// Should not panic the test
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Internal Server Error") {
		t.Error("expected internal server error message")
	}
}
