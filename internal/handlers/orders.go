package handlers

import (
	"log"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
)

type OrdersHandler struct {
	config *config.Config
	client *proxy.Client
}

func NewOrdersHandler(cfg *config.Config, client *proxy.Client) *OrdersHandler {
	return &OrdersHandler{
		config: cfg,
		client: client,
	}
}

func (h *OrdersHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("GetOrder: id=%s", id)
	h.client.ProxyRequest(w, r, h.config.OrdersServiceURL+"/v2/orders/"+id)
}

func (h *OrdersHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateOrder")
	h.client.ProxyRequest(w, r, h.config.OrdersServiceURL+"/v2/orders")
}

func (h *OrdersHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("UpdateOrderStatus: id=%s", id)
	h.client.ProxyRequest(w, r, h.config.OrdersServiceURL+"/v2/orders/"+id+"/status")
}

func (h *OrdersHandler) ListUserOrders(w http.ResponseWriter, r *http.Request) {
	log.Printf("ListUserOrders")
	h.client.ProxyRequest(w, r, h.config.OrdersServiceURL+"/v2/orders")
}

func (h *OrdersHandler) GetOrderV1(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("GetOrderV1: id=%s", id)
	h.client.ProxyRequest(w, r, h.config.OrdersServiceURL+"/orders/"+id)
}
