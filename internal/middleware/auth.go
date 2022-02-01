package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
)

type AuthMiddleware struct {
	config *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: cfg,
	}
}

// AuthenticateLegacy uses the X-Legacy-User-Id header for authentication.
func (m *AuthMiddleware) AuthenticateLegacy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-Legacy-User-Id")
		if userID == "" {
			http.Error(w, "Missing X-Legacy-User-Id header", http.StatusUnauthorized)
			return
		}

		log.Printf("Legacy auth: user=%s", userID)

		ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext retrieves the user ID from the request context.
func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(ContextKeyUserID).(string); ok {
		return userID
	}
	return ""
}
