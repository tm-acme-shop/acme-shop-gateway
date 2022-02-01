# Acme Shop API Gateway

API Gateway for Acme Shop microservices.

## Endpoints

- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create user
- `GET /api/v1/orders/{id}` - Get order by ID
- `POST /api/v1/payments` - Process payment
- `POST /api/v1/email/send` - Send email notification

## Running

```bash
go run cmd/gateway/main.go
```
