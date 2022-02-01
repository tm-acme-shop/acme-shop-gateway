package main

import (
	"log"
	"net/http"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/handlers"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/middleware"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/routes"
)

func main() {
	cfg := config.Load()

	client := proxy.NewClient()
	h := handlers.New(cfg, client)
	authMW := middleware.NewAuthMiddleware(cfg)

	handler := routes.Setup(h, authMW, cfg)

	log.Printf("Starting API Gateway on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
