package routes

import (
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/handlers"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/middleware"
)

func Setup(h *handlers.Handlers, authMW *middleware.AuthMiddleware, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	loggingMW := middleware.NewLoggingMiddleware()

	mux.Handle("GET /api/v1/users/{id}", authMW.AuthenticateLegacy(http.HandlerFunc(h.Users.GetUserV1)))
	mux.Handle("POST /api/v1/users", authMW.AuthenticateLegacy(http.HandlerFunc(h.Users.CreateUserV1)))
	mux.Handle("GET /api/v1/orders/{id}", authMW.AuthenticateLegacy(http.HandlerFunc(h.Orders.GetOrderV1)))
	mux.Handle("POST /api/v1/payments", authMW.AuthenticateLegacy(http.HandlerFunc(h.Payments.ProcessPaymentV1)))
	mux.Handle("POST /api/v1/email/send", authMW.AuthenticateLegacy(http.HandlerFunc(h.Notifications.SendEmailLegacy)))

	var handler http.Handler = mux
	handler = loggingMW.Log(handler)

	return handler
}
