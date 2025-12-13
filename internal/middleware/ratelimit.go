package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-shared-go/logging"
)

type RateLimitMiddleware struct {
	config   *config.Config
	clients  map[string]*clientLimit
	mu       sync.RWMutex
	cleanupC chan struct{}
}

type clientLimit struct {
	tokens     int
	lastRefill time.Time
}

func NewRateLimitMiddleware(cfg *config.Config) *RateLimitMiddleware {
	rl := &RateLimitMiddleware{
		config:   cfg,
		clients:  make(map[string]*clientLimit),
		cleanupC: make(chan struct{}),
	}

	go rl.cleanup()

	return rl
}

func (m *RateLimitMiddleware) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		if !m.allow(clientIP) {
			logging.Warn("Rate limit exceeded", logging.Fields{
				"client_ip": clientIP,
				"path":      r.URL.Path,
			})
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

func (m *RateLimitMiddleware) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			threshold := time.Now().Add(-10 * time.Minute)
			cleanedCount := 0
			for ip, client := range m.clients {
				if client.lastRefill.Before(threshold) {
					delete(m.clients, ip)
					cleanedCount++
				}
			}
			// TODO(TEAM-PLATFORM): Migrate to structured logging
			log.Printf("Rate limiter cleanup: removed %d stale clients", cleanedCount)
			m.mu.Unlock()
		case <-m.cleanupC:
			return
		}
	}
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
