package jwt_test

import (
	"testing"

	"github.com/tm-acme-shop/acme-shop-gateway/internal/jwt"
)

func TestParser_GenerateAndParse(t *testing.T) {
	parser := jwt.NewParser("test-secret")

	token, err := parser.Generate("user123", "test@example.com", "customer")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}

	claims, err := parser.Parse(token)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if claims.UserID != "user123" {
		t.Errorf("expected user_id 'user123', got '%s'", claims.UserID)
	}

	if claims.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", claims.Email)
	}

	if claims.Role != "customer" {
		t.Errorf("expected role 'customer', got '%s'", claims.Role)
	}
}

func TestParser_ParseInvalidToken(t *testing.T) {
	parser := jwt.NewParser("test-secret")

	_, err := parser.Parse("invalid-token")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestParser_ParseWrongSecret(t *testing.T) {
	parser1 := jwt.NewParser("secret1")
	parser2 := jwt.NewParser("secret2")

	token, err := parser1.Generate("user123", "test@example.com", "customer")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = parser2.Parse(token)
	if err == nil {
		t.Error("expected error when parsing with wrong secret")
	}
}
