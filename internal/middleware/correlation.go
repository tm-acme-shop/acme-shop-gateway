package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	HeaderRequestID     = "X-Acme-Request-ID"
	ContextKeyRequestID contextKey = "request_id"
)

type CorrelationMiddleware struct{}

func NewCorrelationMiddleware() *CorrelationMiddleware {
	return &CorrelationMiddleware{}
}

func (m *CorrelationMiddleware) AddRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(HeaderRequestID)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), ContextKeyRequestID, requestID)

		w.Header().Set(HeaderRequestID, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(ContextKeyRequestID).(string); ok {
		return requestID
	}
	return ""
}
