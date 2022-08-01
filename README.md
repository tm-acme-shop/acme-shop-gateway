# Acme Shop API Gateway

API Gateway for Acme Shop microservices.

## Endpoints

### v2 API (Modern)
- `GET /api/v2/users/{id}` - Get user by ID
- `POST /api/v2/users` - Create user (admin only)
- `GET /api/v2/orders/{id}` - Get order by ID
- `POST /api/v2/orders` - Create order
- `POST /api/v2/payments` - Process payment
- `POST /api/v2/notifications` - Send notification
- `POST /api/v2/notifications/email` - Send email

### v1 API (Legacy)
- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create user
- `GET /api/v1/orders/{id}` - Get order by ID
- `POST /api/v1/payments` - Process payment
- `POST /api/v1/email/send` - Send email notification

### Authentication
- `POST /auth/login` - Login with JWT
- `POST /auth/login/legacy` - Legacy login
- `POST /auth/refresh` - Refresh token
- `POST /auth/logout` - Logout

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| GATEWAY_PORT | 8080 | Server port |
| ENABLE_V1_API | true | Enable legacy v1 routes |
| ENABLE_LEGACY_AUTH | true | Enable legacy auth endpoint |
| JWT_SECRET | your-secret-key | JWT signing secret |
| RATE_LIMIT_RPS | 100 | Rate limit requests per second |

## Running

```bash
go run cmd/gateway/main.go
```
