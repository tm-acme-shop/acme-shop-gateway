package middleware

import (
	"net/http"
	"time"

	"github.com/tm-acme-shop/acme-shop-shared-go/logging"
)

type LoggingMiddleware struct{}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

func (m *LoggingMiddleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// TODO(TEAM-PLATFORM): Migrate all logging to structured format
		logging.Infof("[%s] %s %s - %d (%dms)",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			rw.statusCode,
			duration.Milliseconds(),
		)

		logging.Info("Request completed", logging.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      rw.statusCode,
			"duration_ms": duration.Milliseconds(),
			"bytes":       rw.bytes,
			"user_agent":  r.UserAgent(),
		})
	})
}
