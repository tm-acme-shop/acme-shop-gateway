package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/middleware"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
	"github.com/tm-acme-shop/acme-shop-shared-go/logging"
)

type PaymentsHandler struct {
	proxy *proxy.Client
}

func NewPaymentsHandler(proxy *proxy.Client) *PaymentsHandler {
	return &PaymentsHandler{proxy: proxy}
}

func (h *PaymentsHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	req["user_id"] = userID

	logging.Info("Processing payment", logging.Fields{
		"user_id":  userID,
		"order_id": req["order_id"],
		"amount":   req["amount"],
	})

	body, status, err := h.proxy.ProxyToPayments(r.Context(), "POST", "/api/v2/payments", req)
	if err != nil {
		logging.Error("Failed to process payment", logging.Fields{"error": err.Error()})
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *PaymentsHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("id")
	if paymentID == "" {
		http.Error(w, "Payment ID required", http.StatusBadRequest)
		return
	}

	body, status, err := h.proxy.ProxyToPayments(r.Context(), "GET", "/api/v2/payments/"+paymentID, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *PaymentsHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("id")
	if paymentID == "" {
		http.Error(w, "Payment ID required", http.StatusBadRequest)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	body, status, err := h.proxy.ProxyToPayments(r.Context(), "POST", "/api/v2/payments/"+paymentID+"/refund", req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *PaymentsHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("X-Stripe-Signature")
	if signature == "" {
		signature = r.Header.Get("X-Webhook-Signature")
	}

	logging.Info("Received payment webhook", logging.Fields{
		"event_type": payload["type"],
	})

	body, status, err := h.proxy.ProxyToPayments(r.Context(), "POST", "/api/v2/payments/webhook", payload)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(status)
	w.Write(body)
}

// ProcessPaymentV1 processes a payment using the v1 API format.
// Deprecated: Use ProcessPayment instead.
// TODO(TEAM-PAYMENTS): Remove after v1 API migration
func (h *PaymentsHandler) ProcessPaymentV1(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logging.Warnf("V1 API called: ProcessPaymentV1")

	body, status, err := h.proxy.ProxyToPayments(r.Context(), "POST", "/api/v1/payments", req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-API-Deprecated", "true")
	w.WriteHeader(status)
	w.Write(body)
}
