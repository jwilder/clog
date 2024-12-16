package clog

import (
	"net/http"
	"strconv"
	"time"
)

// CanonicalLogger is a middleware that logs the canonical logging context at the end of the request.
type CanonicalLogger struct {
	wrapped http.Handler
	logFn   func(string)
}

func NewCanonicalLogger(wrapped http.Handler, logFn func(string)) http.Handler {
	if logFn == nil {
		panic("logFn cannot be nil")
	}
	if wrapped == nil {
		panic("wrapped cannot be nil")
	}
	return &CanonicalLogger{wrapped: wrapped, logFn: logFn}
}

func (cl *CanonicalLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(Init(r.Context()))
	SetString(r.Context(), "http.request.method", r.Method)
	SetString(r.Context(), "http.request.path", r.URL.Path)

	start := time.Now()
	resp := &loggingResponseWriter{ResponseWriter: w}
	cl.wrapped.ServeHTTP(resp, r)
	duration := time.Since(start)

	SetInt(r.Context(), "http.response.duration_ms", int(duration.Milliseconds()))

	requestSize, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	SetInt(r.Context(), "http.request.body_bytes", requestSize)

	responseSize, _ := strconv.Atoi(w.Header().Get("Content-Length"))
	SetInt(r.Context(), "http.response.body_bytes", responseSize)
	SetInt(r.Context(), "http.response.status_code", resp.statusCode)

	cl.logFn(MarshalJSON(r.Context()))
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
