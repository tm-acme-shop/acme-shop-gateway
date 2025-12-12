module github.com/tm-acme-shop/acme-shop-gateway

go 1.21

require (
	github.com/tm-acme-shop/acme-shop-shared-go v0.0.0
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/prometheus/client_golang v1.17.0
)

replace github.com/tm-acme-shop/acme-shop-shared-go => ../acme-shop-shared-go
