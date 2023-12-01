package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
	"github.com/tm-acme-shop/acme-shop-shared-go/logging"
)

type UsersHandler struct {
	proxy *proxy.Client
}

func NewUsersHandler(proxy *proxy.Client) *UsersHandler {
	return &UsersHandler{proxy: proxy}
}

func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	logging.Info("Getting user", logging.Fields{"user_id": userID})

	body, status, err := h.proxy.ProxyToUsers(r.Context(), "GET", "/api/v2/users/"+userID, nil)
	if err != nil {
		logging.Error("Failed to get user", logging.Fields{"error": err.Error()})
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logging.Info("Creating user", logging.Fields{"email": req["email"]})

	body, status, err := h.proxy.ProxyToUsers(r.Context(), "POST", "/api/v2/users", req)
	if err != nil {
		logging.Error("Failed to create user", logging.Fields{"error": err.Error()})
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	body, status, err := h.proxy.ProxyToUsers(r.Context(), "PUT", "/api/v2/users/"+userID, req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	body, status, err := h.proxy.ProxyToUsers(r.Context(), "DELETE", "/api/v2/users/"+userID, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(status)
	w.Write(body)
}

func (h *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.RawQuery
	path := "/api/v2/users"
	if query != "" {
		path = path + "?" + query
	}

	body, status, err := h.proxy.ProxyToUsers(r.Context(), "GET", path, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

// GetUserV1 retrieves a user using the v1 API format.
// Deprecated: Use GetUser instead.
// TODO(TEAM-API): Remove after v1 deprecation deadline (Q1 2024)
func (h *UsersHandler) GetUserV1(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	logging.Warnf("V1 API called: GetUserV1 for user %s", userID)

	body, status, err := h.proxy.ProxyToUsersLegacy(r.Context(), "GET", "/users/"+userID, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-API-Deprecated", "true")
	w.Header().Set("X-API-Sunset", "2024-06-01")
	w.WriteHeader(status)
	w.Write(body)
}

// CreateUserV1 creates a user using the v1 API format.
// Deprecated: Use CreateUser instead.
func (h *UsersHandler) CreateUserV1(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logging.Warnf("V1 API called: CreateUserV1")

	body, status, err := h.proxy.ProxyToUsersLegacy(r.Context(), "POST", "/users", req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-API-Deprecated", "true")
	w.WriteHeader(status)
	w.Write(body)
}
