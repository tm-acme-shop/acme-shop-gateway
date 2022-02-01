package handlers

import (
	"log"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
)

type NotificationsHandler struct {
	config *config.Config
	client *proxy.Client
}

func NewNotificationsHandler(cfg *config.Config, client *proxy.Client) *NotificationsHandler {
	return &NotificationsHandler{
		config: cfg,
		client: client,
	}
}

func (h *NotificationsHandler) SendEmailLegacy(w http.ResponseWriter, r *http.Request) {
	log.Printf("SendEmailLegacy")
	h.client.ProxyRequest(w, r, h.config.NotificationsServiceURL+"/email/send")
}
