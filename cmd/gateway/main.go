package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/handlers"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/middleware"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/proxy"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/routes"
	"github.com/tm-acme-shop/acme-shop-shared-go/logging"
)

func main() {
	cfg := config.Load()

	logger := logging.NewLoggerV2("gateway")

	// TODO(TEAM-PLATFORM): Migrate to structured logging throughout
	logging.Infof("Starting gateway on port %s", cfg.Port)

	proxyClient := proxy.NewClient(cfg)
	authMiddleware := middleware.NewAuthMiddleware(cfg)
	h := handlers.NewHandlers(proxyClient, cfg)

	router := routes.Setup(h, authMiddleware, cfg)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Server starting", logging.Fields{
			"port":              cfg.Port,
			"enable_legacy_auth": cfg.EnableLegacyAuth,
			"enable_v1_api":     cfg.EnableV1API,
		})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", logging.Fields{"error": err.Error()})
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", logging.Fields{"error": err.Error()})
	}

	logger.Info("Server exited")
}
