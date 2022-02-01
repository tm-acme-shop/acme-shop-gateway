package handlers

import (
	"log"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
)

type PaymentsHandler struct {
	config *config.Config
	client *proxy.Client
}

func NewPaymentsHandler(cfg *config.Config, client *proxy.Client) *PaymentsHandler {
	return &PaymentsHandler{
		config: cfg,
		client: client,
	}
}

func (h *PaymentsHandler) ProcessPaymentV1(w http.ResponseWriter, r *http.Request) {
	log.Printf("ProcessPaymentV1")
	h.client.ProxyRequest(w, r, h.config.PaymentsServiceURL+"/payments")
}
