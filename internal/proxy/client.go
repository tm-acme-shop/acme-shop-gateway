package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/config"
	"github.com/tm-acme-shop/acme-shop-gateway/internal/middleware"
	"github.com/tm-acme-shop/acme-shop-shared-go/logging"
)

type Client struct {
	httpClient *http.Client
	config     *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.RequestTimeout) * time.Second,
		},
		config: cfg,
	}
}

func (c *Client) ProxyToUsers(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
	return c.proxy(ctx, c.config.UsersServiceURL, method, path, body)
}

func (c *Client) ProxyToOrders(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
	return c.proxy(ctx, c.config.OrdersServiceURL, method, path, body)
}

func (c *Client) ProxyToPayments(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
	return c.proxy(ctx, c.config.PaymentsServiceURL, method, path, body)
}

func (c *Client) ProxyToNotifications(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
	return c.proxy(ctx, c.config.NotificationsServiceURL, method, path, body)
}

func (c *Client) proxy(ctx context.Context, baseURL, method, path string, body interface{}) ([]byte, int, error) {
	url := baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	requestID := middleware.GetRequestIDFromContext(ctx)
	if requestID != "" {
		req.Header.Set("X-Acme-Request-ID", requestID)
	}

	userID := middleware.GetUserIDFromContext(ctx)
	if userID != "" {
		req.Header.Set("X-User-Id", userID)
		req.Header.Set("X-Legacy-User-Id", userID)
	}

	logging.Info("Proxying request", logging.Fields{
		"method":     method,
		"url":        url,
		"request_id": requestID,
	})

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logging.Error("Proxy request failed", logging.Fields{
			"method": method,
			"url":    url,
			"error":  err.Error(),
		})
		return nil, 0, fmt.Errorf("proxy request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// ProxyToUsersLegacy proxies requests using the old API format.
// Deprecated: Use ProxyToUsers instead.
// TODO(TEAM-API): Remove after v1 API deprecation
func (c *Client) ProxyToUsersLegacy(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
	log.Printf("Legacy proxy to users service: %s %s", method, path)
	return c.proxy(ctx, c.config.UsersServiceURL, method, "/v1"+path, body)
}

// ProxyToOrdersLegacy proxies requests using the old v1 API format.
// Deprecated: Use ProxyToOrders instead.
// TODO(TEAM-API): Remove after v1 API deprecation
func (c *Client) ProxyToOrdersLegacy(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
	log.Printf("Legacy proxy to orders service: %s %s", method, path)
	return c.proxy(ctx, c.config.OrdersServiceURL, method, "/v1"+path, body)
}
