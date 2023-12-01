package middleware

import (
	"context"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-shared-go/logging"
	"github.com/tm-acme-shop/acme-shop-shared-go/utils"
)

const (
	HeaderRequestID       = "X-Acme-Request-ID"
	HeaderLegacyRequestID = "X-Request-ID"
	ContextKeyRequestID   contextKey = "request_id"
)

type CorrelationMiddleware struct{}

func NewCorrelationMiddleware() *CorrelationMiddleware {
	return &CorrelationMiddleware{}
}

func (m *CorrelationMiddleware) AddRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(HeaderRequestID)

		// TODO(TEAM-API): Remove legacy header support after migration
		if requestID == "" {
			requestID = r.Header.Get(HeaderLegacyRequestID)
			if requestID != "" {
				logging.Warnf("Legacy X-Request-ID header used, migrate to X-Acme-Request-ID")
			}
		}

		if requestID == "" {
			requestID = utils.GenerateID("req")
		}

		ctx := context.WithValue(r.Context(), ContextKeyRequestID, requestID)
		ctx = context.WithValue(ctx, logging.ContextKeyRequestID, requestID)

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
