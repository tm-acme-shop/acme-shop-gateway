package handlers

import (
	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
)

type Handlers struct {
	Users         *UsersHandler
	Orders        *OrdersHandler
	Payments      *PaymentsHandler
	Notifications *NotificationsHandler
	Health        *HealthHandler
	Auth          *AuthHandler
}

func NewHandlers(proxyClient *proxy.Client, cfg *config.Config) *Handlers {
	return &Handlers{
		Users:         NewUsersHandler(proxyClient),
		Orders:        NewOrdersHandler(proxyClient),
		Payments:      NewPaymentsHandler(proxyClient),
		Notifications: NewNotificationsHandler(proxyClient),
		Health:        NewHealthHandler(),
		Auth:          NewAuthHandler(cfg),
	}
}
