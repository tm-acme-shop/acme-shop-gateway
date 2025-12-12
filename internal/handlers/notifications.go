package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
	"github.com/tm-acme-shop/acme-shop-shared-go/logging"
)

type NotificationsHandler struct {
	proxy *proxy.Client
}

func NewNotificationsHandler(proxy *proxy.Client) *NotificationsHandler {
	return &NotificationsHandler{proxy: proxy}
}

func (h *NotificationsHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logging.Info("Sending notification", logging.Fields{
		"type":      req["type"],
		"recipient": req["recipient"],
	})

	body, status, err := h.proxy.ProxyToNotifications(r.Context(), "POST", "/api/v2/notifications", req)
	if err != nil {
		logging.Error("Failed to send notification", logging.Fields{"error": err.Error()})
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *NotificationsHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	notificationID := r.PathValue("id")
	if notificationID == "" {
		http.Error(w, "Notification ID required", http.StatusBadRequest)
		return
	}

	body, status, err := h.proxy.ProxyToNotifications(r.Context(), "GET", "/api/v2/notifications/"+notificationID, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *NotificationsHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	body, status, err := h.proxy.ProxyToNotifications(r.Context(), "POST", "/api/v2/notifications/email", req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *NotificationsHandler) SendSMS(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	body, status, err := h.proxy.ProxyToNotifications(r.Context(), "POST", "/api/v2/notifications/sms", req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

// SendEmailLegacy sends an email using the legacy format.
// Deprecated: Use SendEmail instead.
// TODO(TEAM-NOTIFICATIONS): Migrate all callers to new API
func (h *NotificationsHandler) SendEmailLegacy(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logging.Warnf("Legacy email API called")

	body, status, err := h.proxy.ProxyToNotifications(r.Context(), "POST", "/api/v1/email/send", req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-API-Deprecated", "true")
	w.WriteHeader(status)
	w.Write(body)
}
