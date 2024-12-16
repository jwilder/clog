package clog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanonicalLogger_ServeHTTP(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Length", "2")
		_, _ = w.Write([]byte("OK"))
	})
	logFn := func(log string) {
		require.JSONEq(t, `{"http":{"request":{"method":"GET","path":"/test","body_bytes":0},"response":{"duration_ms":0,"body_bytes":2,"status_code":200}}}`, log)
	}
	logger := NewCanonicalLogger(handler, logFn)

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	logger.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCanonicalLogger_ServeHTTP_NilLogFn(t *testing.T) {
	require.PanicsWithValue(t, "logFn cannot be nil", func() {
		NewCanonicalLogger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), nil)
	})
}

func TestCanonicalLogger_ServeHTTP_NilHandler(t *testing.T) {
	require.PanicsWithValue(t, "wrapped cannot be nil", func() {
		NewCanonicalLogger(nil, func(log string) {})
	})
}

func TestCanonicalLogger_ServeHTTP_InvalidContentLength(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "invalid")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	logFn := func(log string) {
		require.JSONEq(t, `{"http":{"request":{"method":"GET","path":"/test","body_bytes":0},"response":{"duration_ms":0,"body_bytes":0,"status_code":200}}}`, log)
	}
	logger := NewCanonicalLogger(handler, logFn)

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	logger.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
