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

func (h *NotificationsHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	log.Printf("SendNotification")
	h.client.ProxyRequest(w, r, h.config.NotificationsServiceURL+"/v2/notifications")
}

func (h *NotificationsHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("GetNotification: id=%s", id)
	h.client.ProxyRequest(w, r, h.config.NotificationsServiceURL+"/v2/notifications/"+id)
}

func (h *NotificationsHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	log.Printf("SendEmail")
	h.client.ProxyRequest(w, r, h.config.NotificationsServiceURL+"/v2/notifications/email")
}

func (h *NotificationsHandler) SendSMS(w http.ResponseWriter, r *http.Request) {
	log.Printf("SendSMS")
	h.client.ProxyRequest(w, r, h.config.NotificationsServiceURL+"/v2/notifications/sms")
}

func (h *NotificationsHandler) SendEmailLegacy(w http.ResponseWriter, r *http.Request) {
	log.Printf("SendEmailLegacy")
	h.client.ProxyRequest(w, r, h.config.NotificationsServiceURL+"/email/send")
}
