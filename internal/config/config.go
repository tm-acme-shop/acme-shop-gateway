package config

import (
	"os"
)

type Config struct {
	Port                    string
	UsersServiceURL         string
	OrdersServiceURL        string
	PaymentsServiceURL      string
	NotificationsServiceURL string
}

func Load() *Config {
	return &Config{
		Port:                    getEnv("GATEWAY_PORT", "8080"),
		UsersServiceURL:         getEnv("USERS_SERVICE_URL", "http://localhost:8081"),
		OrdersServiceURL:        getEnv("ORDERS_SERVICE_URL", "http://localhost:8082"),
		PaymentsServiceURL:      getEnv("PAYMENTS_SERVICE_URL", "http://localhost:8083"),
		NotificationsServiceURL: getEnv("NOTIFICATIONS_SERVICE_URL", "http://localhost:8084"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
