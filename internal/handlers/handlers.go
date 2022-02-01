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
}

func New(cfg *config.Config, client *proxy.Client) *Handlers {
	return &Handlers{
		Users:         NewUsersHandler(cfg, client),
		Orders:        NewOrdersHandler(cfg, client),
		Payments:      NewPaymentsHandler(cfg, client),
		Notifications: NewNotificationsHandler(cfg, client),
	}
}
