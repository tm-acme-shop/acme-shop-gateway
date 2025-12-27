package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                    string
	UsersServiceURL         string
	OrdersServiceURL        string
	PaymentsServiceURL      string
	NotificationsServiceURL string
	JWTSecret               string
	EnableNewAuth           bool
	EnableV1API             bool
	RateLimitRPS            int
	RequestTimeout          int
}

func Load() *Config {
	return &Config{
		Port:                    getEnv("GATEWAY_PORT", "8080"),
		UsersServiceURL:         getEnv("USERS_SERVICE_URL", "http://localhost:8081"),
		OrdersServiceURL:        getEnv("ORDERS_SERVICE_URL", "http://localhost:8082"),
		PaymentsServiceURL:      getEnv("PAYMENTS_SERVICE_URL", "http://localhost:8083"),
		NotificationsServiceURL: getEnv("NOTIFICATIONS_SERVICE_URL", "http://localhost:8084"),
		JWTSecret:               getEnv("JWT_SECRET", "your-secret-key"),
		EnableNewAuth:           getEnvBool("ENABLE_NEW_AUTH", true),
		EnableV1API:             getEnvBool("ENABLE_V1_API", true),
		RateLimitRPS:            getEnvInt("RATE_LIMIT_RPS", 100),
		RequestTimeout:          getEnvInt("REQUEST_TIMEOUT_SECONDS", 30),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return b
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		i, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return i
	}
	return defaultValue
}
