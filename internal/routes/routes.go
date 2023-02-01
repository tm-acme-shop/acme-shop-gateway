package routes

import (
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/handlers"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/middleware"
)

func Setup(h *handlers.Handlers, authMW *middleware.AuthMiddleware, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	correlationMW := middleware.NewCorrelationMiddleware()
	loggingMW := middleware.NewLoggingMiddleware()
	rateLimitMW := middleware.NewRateLimitMiddleware(cfg)

	mux.HandleFunc("GET /health", h.Health.Health)
	mux.HandleFunc("GET /ready", h.Health.Ready)
	mux.HandleFunc("GET /metrics", h.Health.Metrics)

	mux.HandleFunc("POST /auth/login", h.Auth.Login)
	mux.HandleFunc("POST /auth/refresh", h.Auth.Refresh)
	mux.HandleFunc("POST /auth/logout", h.Auth.Logout)

	if cfg.EnableLegacyAuth {
		mux.HandleFunc("POST /auth/login/legacy", h.Auth.LoginLegacy)
	}

	mux.Handle("GET /api/v2/users/{id}", authMW.Authenticate(http.HandlerFunc(h.Users.GetUser)))
	mux.Handle("POST /api/v2/users", authMW.Authenticate(authMW.RequireRole("admin")(http.HandlerFunc(h.Users.CreateUser))))
	mux.Handle("PUT /api/v2/users/{id}", authMW.Authenticate(http.HandlerFunc(h.Users.UpdateUser)))
	mux.Handle("DELETE /api/v2/users/{id}", authMW.Authenticate(authMW.RequireRole("admin")(http.HandlerFunc(h.Users.DeleteUser))))
	mux.Handle("GET /api/v2/users", authMW.Authenticate(authMW.RequireRole("admin")(http.HandlerFunc(h.Users.ListUsers))))

	mux.Handle("GET /api/v2/orders/{id}", authMW.Authenticate(http.HandlerFunc(h.Orders.GetOrder)))
	mux.Handle("POST /api/v2/orders", authMW.Authenticate(http.HandlerFunc(h.Orders.CreateOrder)))
	mux.Handle("PATCH /api/v2/orders/{id}/status", authMW.Authenticate(http.HandlerFunc(h.Orders.UpdateOrderStatus)))
	mux.Handle("GET /api/v2/orders", authMW.Authenticate(http.HandlerFunc(h.Orders.ListUserOrders)))

	mux.Handle("POST /api/v2/payments", authMW.Authenticate(http.HandlerFunc(h.Payments.ProcessPayment)))
	mux.Handle("GET /api/v2/payments/{id}", authMW.Authenticate(http.HandlerFunc(h.Payments.GetPayment)))
	mux.Handle("POST /api/v2/payments/{id}/refund", authMW.Authenticate(http.HandlerFunc(h.Payments.RefundPayment)))

	mux.Handle("POST /api/v2/notifications", authMW.Authenticate(http.HandlerFunc(h.Notifications.SendNotification)))
	mux.Handle("GET /api/v2/notifications/{id}", authMW.Authenticate(http.HandlerFunc(h.Notifications.GetNotification)))
	mux.Handle("POST /api/v2/notifications/email", authMW.Authenticate(http.HandlerFunc(h.Notifications.SendEmail)))
	mux.Handle("POST /api/v2/notifications/sms", authMW.Authenticate(http.HandlerFunc(h.Notifications.SendSMS)))

	if cfg.EnableV1API {
		// TODO(TEAM-GATEWAY): Remove legacy v1 routes after clients migrate
		mux.Handle("GET /api/v1/users/{id}", authMW.AuthenticateLegacy(http.HandlerFunc(h.Users.GetUserV1)))
		mux.Handle("POST /api/v1/users", authMW.AuthenticateLegacy(http.HandlerFunc(h.Users.CreateUserV1)))
		mux.Handle("GET /api/v1/orders/{id}", authMW.AuthenticateLegacy(http.HandlerFunc(h.Orders.GetOrderV1)))
		mux.Handle("POST /api/v1/payments", authMW.AuthenticateLegacy(http.HandlerFunc(h.Payments.ProcessPaymentV1)))
		mux.Handle("POST /api/v1/email/send", authMW.AuthenticateLegacy(http.HandlerFunc(h.Notifications.SendEmailLegacy)))
	}

	var handler http.Handler = mux
	handler = loggingMW.Log(handler)
	handler = rateLimitMW.Limit(handler)
	handler = correlationMW.AddRequestID(handler)

	return handler
}
