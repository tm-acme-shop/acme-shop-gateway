package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/middleware"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
	"github.com/tm-acme-shop/acme-shop-shared-go/logging"
)

type OrdersHandler struct {
	proxy *proxy.Client
}

func NewOrdersHandler(proxy *proxy.Client) *OrdersHandler {
	return &OrdersHandler{proxy: proxy}
}

func (h *OrdersHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")
	if orderID == "" {
		http.Error(w, "Order ID required", http.StatusBadRequest)
		return
	}

	logging.Info("Getting order", logging.Fields{"order_id": orderID})

	body, status, err := h.proxy.ProxyToOrders(r.Context(), "GET", "/api/v2/orders/"+orderID, nil)
	if err != nil {
		logging.Error("Failed to get order", logging.Fields{"error": err.Error()})
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *OrdersHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	req["user_id"] = userID

	logging.Info("Creating order", logging.Fields{
		"user_id": userID,
	})

	body, status, err := h.proxy.ProxyToOrders(r.Context(), "POST", "/api/v2/orders", req)
	if err != nil {
		logging.Error("Failed to create order", logging.Fields{"error": err.Error()})
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *OrdersHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")
	if orderID == "" {
		http.Error(w, "Order ID required", http.StatusBadRequest)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	body, status, err := h.proxy.ProxyToOrders(r.Context(), "PATCH", "/api/v2/orders/"+orderID+"/status", req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *OrdersHandler) ListUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	query := r.URL.RawQuery
	path := "/api/v2/users/" + userID + "/orders"
	if query != "" {
		path = path + "?" + query
	}

	body, status, err := h.proxy.ProxyToOrders(r.Context(), "GET", path, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func (h *OrdersHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")
	if orderID == "" {
		http.Error(w, "Order ID required", http.StatusBadRequest)
		return
	}

	body, status, err := h.proxy.ProxyToOrders(r.Context(), "POST", "/api/v2/orders/"+orderID+"/cancel", nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(status)
	w.Write(body)
}

// GetOrderV1 retrieves an order using the v1 API format.
// Deprecated: Use GetOrder instead.
// TODO(TEAM-API): Remove after v1 API sunset
func (h *OrdersHandler) GetOrderV1(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")
	if orderID == "" {
		http.Error(w, "Order ID required", http.StatusBadRequest)
		return
	}

	logging.Warnf("V1 API called: GetOrderV1 for order %s", orderID)

	body, status, err := h.proxy.ProxyToOrders(r.Context(), "GET", "/api/v1/orders/"+orderID, nil)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-API-Deprecated", "true")
	w.WriteHeader(status)
	w.Write(body)
}
