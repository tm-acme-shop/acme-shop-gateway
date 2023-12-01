package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/jwt"
	"github.com/tm-acme-shop/acme-shop-shared-go/logging"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyRole   contextKey = "role"
)

type AuthMiddleware struct {
	config    *config.Config
	jwtParser *jwt.Parser
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config:    cfg,
		jwtParser: jwt.NewParser(cfg.JWTSecret),
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtParser.Parse(parts[1])
		if err != nil {
			logging.Warnf("JWT parse failed: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)

		logging.Info("Request authenticated", logging.Fields{
			"user_id": claims.UserID,
			"role":    claims.Role,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthenticateLegacy uses the old X-Legacy-User-Id header for authentication.
// Deprecated: Use Authenticate with JWT tokens instead.
// TODO(TEAM-SEC): Remove after migration to JWT auth is complete
func (m *AuthMiddleware) AuthenticateLegacy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(TEAM-SEC): This is insecure, should be removed
		userID := r.Header.Get("X-Legacy-User-Id")
		if userID == "" {
			http.Error(w, "Missing X-Legacy-User-Id header", http.StatusUnauthorized)
			return
		}

		logging.Warnf("Legacy auth used for user: %s", userID)

		ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
		ctx = context.WithValue(ctx, ContextKeyRole, "customer")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(ContextKeyRole).(string)
			if !ok {
				http.Error(w, "No role in context", http.StatusForbidden)
				return
			}

			for _, allowedRole := range roles {
				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Insufficient permissions", http.StatusForbidden)
		})
	}
}

// GetUserIDFromContext retrieves the user ID from the request context.
func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(ContextKeyUserID).(string); ok {
		return userID
	}
	return ""
}

// GetRoleFromContext retrieves the role from the request context.
func GetRoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value(ContextKeyRole).(string); ok {
		return role
	}
	return ""
}
