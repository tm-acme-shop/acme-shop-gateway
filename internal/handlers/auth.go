package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/jwt"
	"github.com/tm-acme-shop/acme-shop-shared-go/logging"
)

type AuthHandler struct {
	config    *config.Config
	jwtParser *jwt.Parser
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		config:    cfg,
		jwtParser: jwt.NewParser(cfg.JWTSecret),
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	UserID    string `json:"user_id"`
}

type RefreshRequest struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logging.Info("Login attempt", logging.Fields{"email": req.Email})

	userID := "user_" + req.Email[:3]
	role := "customer"

	token, err := h.jwtParser.Generate(userID, req.Email, role)
	if err != nil {
		logging.Error("Failed to generate token", logging.Fields{"error": err.Error()})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		Token:     token,
		ExpiresIn: 86400,
		UserID:    userID,
	}

	logging.Info("Login successful", logging.Fields{
		"user_id": userID,
		"email":   req.Email,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, err := h.jwtParser.Parse(req.Token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	newToken, err := h.jwtParser.Generate(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		Token:     newToken,
		ExpiresIn: 86400,
		UserID:    claims.UserID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logging.Info("User logged out")
	w.WriteHeader(http.StatusNoContent)
}

// LoginLegacy handles login using the old authentication method.
// Deprecated: Use Login instead.
// TODO(TEAM-SEC): Remove after migration to new auth
func (h *AuthHandler) LoginLegacy(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logging.Warnf("Legacy login used for: %s", req.Email)

	userID := "legacy_" + req.Email[:3]

	resp := map[string]interface{}{
		"user_id":    userID,
		"session_id": "sess_" + userID,
		"deprecated": true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-API-Deprecated", "true")
	w.Header().Set("X-Legacy-User-Id", userID)
	json.NewEncoder(w).Encode(resp)
}
