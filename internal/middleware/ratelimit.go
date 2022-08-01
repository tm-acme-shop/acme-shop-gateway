package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
)

type RateLimitMiddleware struct {
	config  *config.Config
	clients map[string]*clientLimit
	mu      sync.RWMutex
}

type clientLimit struct {
	tokens     int
	lastRefill time.Time
}

func NewRateLimitMiddleware(cfg *config.Config) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		config:  cfg,
		clients: make(map[string]*clientLimit),
	}
}

func (m *RateLimitMiddleware) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		if !m.allow(clientIP) {
			log.Printf("Rate limit exceeded for %s", clientIP)
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *RateLimitMiddleware) allow(clientIP string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	client, exists := m.clients[clientIP]

	if !exists {
		m.clients[clientIP] = &clientLimit{
			tokens:     m.config.RateLimitRPS - 1,
			lastRefill: now,
		}
		return true
	}

	elapsed := now.Sub(client.lastRefill)
	refillTokens := int(elapsed.Seconds()) * m.config.RateLimitRPS

	if refillTokens > 0 {
		client.tokens = min(m.config.RateLimitRPS, client.tokens+refillTokens)
		client.lastRefill = now
	}

	if client.tokens > 0 {
		client.tokens--
		return true
	}

	return false
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
