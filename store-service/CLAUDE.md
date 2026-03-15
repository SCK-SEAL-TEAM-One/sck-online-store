# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this directory.

## Overview

Go backend API service using Gin framework, MySQL via SQLx, JWT authentication. Runs on port 8000.

## Commands

```bash
# Run locally (dev mode, from project root)
make store_service_dev_mode

# Unit tests (generates report.xml)
cd store-service && go test -v ./...

# Run tests for a specific package
cd store-service && go test -v ./internal/order/...

# Integration tests (requires DB + thirdparty running)
cd store-service && go test -tags=integration ./...

# Code analysis
cd store-service && go vet ./...

# Code coverage
cd store-service && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# Generate Swagger docs
cd store-service && swag init -g cmd/main.go -o cmd/docs
```

## Architecture

```
cmd/
  main.go          — Entry point, route registration, DI wiring
  api/             — HTTP handlers (Gin handlers with Swagger annotations)
  docs/            — Auto-generated Swagger docs (DO NOT EDIT)
internal/
  auth/            — JWT auth, user repository, login
  order/           — Order service, repository, PDF generation, helpers
  cart/            — Shopping cart service + repository
  product/         — Product catalog, currency conversion (USD→THB)
  payment/         — Payment service, bank gateway client
  shipping/        — Shipping service, repository, gateway client
  point/           — Points integration, gateway to point-service
  user/            — User model, bcrypt password hashing
  middleware/      — JWT auth middleware
  common/          — Shared utilities (currency, point calc, decimal formatting)
  healthcheck/     — Health check with DB validation
  seed/            — Database seeding
```

**Layered pattern per package:**
- `{package}.go` — Service with business logic
- `repository.go` — Data access (MySQL/SQLx)
- `model.go` — Domain models
- `{package}_test.go` — Unit tests (testify mocks)
- `mock_{package}_test.go` — Mock implementations
- `repository_test.go` — Integration tests (build tag `integration`)

Services define interfaces for dependencies, enabling mock-based unit testing with `github.com/stretchr/testify/mock`.

## API Routes

Base path: `/api/v1`

**Public:** `POST /login`, `GET /refreshToken`, `GET /health`
**Protected (JWT):** `GET /product`, `GET /product/:id`, `GET /cart`, `PUT /addCart`, `PUT /updateCart`, `POST /order`, `POST /order/:id/summary`, `POST /confirmPayment`, `GET /point`, `POST /point`
**Docs:** `GET /swagger/*any`

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_CONNECTION` | — | MySQL connection string (`user:password@(host:port)/store?parseTime=true`) |
| `JWT_SECRET` | — | Required, fatal if missing |
| `BANK_GATEWAY` | `thirdparty:8882` | Payment mock |
| `SHIPPING_GATEWAY` | `thirdparty:8883` | Shipping mock |
| `POINT_GATEWAY` | `point-service:8001` | Points service |
| `STORE_WEB_HOST` | `http://localhost:3000` | CORS origin whitelist |

## Naming Conventions

- Functions: PascalCase (`CalculateTotalPrice()`)
- Files: PascalCase (`OrderService.go`, `OrderService_test.go`)
- Packages/directories: lowercase
- Test functions: Snake_Case (`Test_CalculateAge_Input_Birth_Date_18042003_Should_be_16`)
- Variables: camelCase; constants: UPPERCASE

## Key Dependencies

Gin, SQLx, golang-jwt/v5, stretchr/testify, swaggo/swag, maroto/v2 (PDF), elastic APM
