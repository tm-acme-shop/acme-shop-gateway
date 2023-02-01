package handlers

import (
	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
)

type Handlers struct {
	Auth          *AuthHandler
	Health        *HealthHandler
	Users         *UsersHandler
	Orders        *OrdersHandler
	Payments      *PaymentsHandler
	Notifications *NotificationsHandler
}

func New(cfg *config.Config, client *proxy.Client) *Handlers {
	return &Handlers{
		Auth:          NewAuthHandler(cfg),
		Health:        NewHealthHandler(),
		Users:         NewUsersHandler(cfg, client),
		Orders:        NewOrdersHandler(cfg, client),
		Payments:      NewPaymentsHandler(cfg, client),
		Notifications: NewNotificationsHandler(cfg, client),
	}
}
