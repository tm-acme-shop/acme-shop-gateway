package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/middleware"
)

func TestCorrelationMiddleware_AddRequestID(t *testing.T) {
	mw := middleware.NewCorrelationMiddleware()

	handler := mw.AddRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetRequestIDFromContext(r.Context())
		if requestID == "" {
			t.Error("expected request ID in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Header().Get("X-Acme-Request-ID") == "" {
		t.Error("expected X-Acme-Request-ID header in response")
	}
}

func TestCorrelationMiddleware_PassThroughRequestID(t *testing.T) {
	mw := middleware.NewCorrelationMiddleware()

	expectedID := "test-request-123"

	handler := mw.AddRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetRequestIDFromContext(r.Context())
		if requestID != expectedID {
			t.Errorf("expected request ID '%s', got '%s'", expectedID, requestID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Acme-Request-ID", expectedID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Header().Get("X-Acme-Request-ID") != expectedID {
		t.Errorf("expected X-Acme-Request-ID '%s', got '%s'", expectedID, w.Header().Get("X-Acme-Request-ID"))
	}
}

func TestLoggingMiddleware_Log(t *testing.T) {
	mw := middleware.NewLoggingMiddleware()

	handler := mw.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
