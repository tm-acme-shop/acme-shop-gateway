# AcmeShop Gateway

API Gateway service for the AcmeShop platform. Routes requests to internal microservices and handles authentication, rate limiting, and request correlation.

## Features

- JWT-based authentication
- Request routing to internal services
- Rate limiting
- Request correlation (X-Acme-Request-ID)
- Metrics and health checks
- v1 and v2 API versioning

## Running

```bash
go run ./cmd/gateway

# With environment variables
GATEWAY_PORT=8080 ENABLE_NEW_AUTH=true go run ./cmd/gateway
```

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `GATEWAY_PORT` | HTTP port | `8080` |
| `USERS_SERVICE_URL` | Users service URL | `http://localhost:8081` |
| `ORDERS_SERVICE_URL` | Orders service URL | `http://localhost:8082` |
| `PAYMENTS_SERVICE_URL` | Payments service URL | `http://localhost:8083` |
| `ENABLE_NEW_AUTH` | Enable new auth endpoints | `false` |
| `ENABLE_V1_API` | Enable v1 API routes | `true` |

## API Endpoints

### v2 (Current)
- `GET /api/v2/users/:id` - Get user by ID
- `POST /api/v2/users` - Create user
- `GET /api/v2/orders/:id` - Get order by ID
- `POST /api/v2/orders` - Create order
- `POST /api/v2/payments` - Process payment

### v1 (Deprecated)
- `GET /api/v1/users/:id` - Legacy get user
- `POST /api/v1/users` - Legacy create user
- `GET /api/v1/orders/:id` - Legacy get order

## Architecture

```
┌─────────────────┐
│   API Gateway   │
├─────────────────┤
│ - Auth MW       │
│ - Rate Limit    │
│ - Correlation   │
└────────┬────────┘
         │
    ┌────┴────┬────────┬──────────┐
    ▼         ▼        ▼          ▼
┌───────┐ ┌───────┐ ┌────────┐ ┌──────────┐
│ Users │ │Orders │ │Payments│ │Notif.    │
└───────┘ └───────┘ └────────┘ └──────────┘
```
