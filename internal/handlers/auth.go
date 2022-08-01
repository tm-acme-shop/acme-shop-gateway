package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/jwt"
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
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt: %s", req.Username)

	token, err := h.jwtParser.Generate(req.Username, "customer")
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

func (h *AuthHandler) LoginLegacy(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Legacy login: %s", req.Username)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Legacy-User-Id", req.Username)
	json.NewEncoder(w).Encode(map[string]string{"user_id": req.Username})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	log.Printf("Token refresh")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("Logout")
	w.WriteHeader(http.StatusOK)
}
