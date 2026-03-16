# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**sck-online-store** is a microservices-based e-commerce application used for workshop teaching. It consists of three main services:

- **store-service** (Go/Gin) — Main backend API on port 8000
- **store-web** (Next.js 14/React 18) — Frontend on port 3000, uses TailwindCSS + DaisyUI, Zustand for state
- **point-service** (NestJS/TypeORM) — Points/rewards service on port 8001
- **thirdparty** (Node.js) — Mock payment (8882) and shipping (8883) gateways

Nginx reverse proxy sits in front on port 80. MySQL 8.3 with Liquibase migrations for the database.

## Common Commands

### Start/Stop All Services
```bash
make start_all          # Docker Compose up all services (with build)
make down               # Docker Compose down
```

### Run Backend (Go) in Dev Mode
```bash
make store_service_dev_mode   # Runs store-service locally with env vars
```

### Unit Tests
```bash
# Go backend unit tests (generates report.xml)
make backend_unit_test

# Run a single Go test file/package directly
cd store-service && go test -v ./internal/order/...

# Go code coverage
make code-coverage

# NestJS point-service tests
cd point-service && npm test

# Next.js component tests (Cypress)
cd store-web && npm run test:component
```

### Integration Tests
```bash
make backend_integration_test   # Starts DB+thirdparty, runs integration tests, then docker down
```

### Full ATDD Test Suite
```bash
make start_test_suite       # Start all services for testing
make run_newman             # API tests (Newman/Postman)
make run_robot              # UI tests (Robot Framework)
make stop_test_suite        # Tear down

# Run specific test suites
make run_newman_authentication
make run_newman_order_summary_pdf
make run_robot_authentication
make run_robot_order_summary_pdf
```

### Code Analysis
```bash
make code_analysis_frontend    # npm run lint (store-web)
make code_analysis_backend     # go vet ./... (store-service)
```

### Development Workflow (run before commit)
```bash
make test_all              # Full pipeline: analysis → unit → ATDD (API + UI)
make unit_test_all         # All unit tests: Go + Jest + Cypress component
make code_analysis_all     # All linting: go vet + npm run lint
```

### Build & Generate
```bash
make build_backend         # Docker build store-service
make build_frontend        # Docker build store-web
make gen-swagger           # Generate Swagger docs from Go annotations
```

## Architecture

```
                    ┌──────────┐
                    │  nginx   │ :80
                    └────┬─────┘
              ┌──────────┼──────────┐
              ▼                     ▼
       ┌─────────────┐     ┌──────────────┐
       │  store-web   │     │store-service │
       │  (Next.js)   │     │   (Go/Gin)   │
       │    :3000     │     │    :8000     │
       └─────────────┘     └──┬───┬───┬───┘
                              │   │   │
                    ┌─────────┘   │   └─────────┐
                    ▼             ▼              ▼
             ┌───────────┐ ┌──────────┐  ┌────────────┐
             │point-svc  │ │   MySQL  │  │ thirdparty │
             │ (NestJS)  │ │  :3306   │  │ :8882/8883 │
             │  :8001    │ └──────────┘  └────────────┘
             └───────────┘
```

### store-service (Go) internal structure
- `cmd/main.go` — Entry point
- `internal/auth/` — JWT authentication, user management
- `internal/order/` — Order service + repository
- `internal/cart/` — Shopping cart
- `internal/product/` — Product catalog
- `internal/payment/` — Payment processing (calls thirdparty bank gateway)
- `internal/shipping/` — Shipping (calls thirdparty shipping gateway)
- `internal/point/` — Points integration (calls point-service)
- `internal/middleware/` — HTTP middleware
- Key deps: Gin, SQLx, gin-swagger, Elastic APM

### store-web (Next.js) structure
- `src/app/` — App router pages
- `src/components/` — React components
- `src/services/` — API clients (Axios)
- `src/hooks/` — Custom hooks
- Key deps: Zustand, Axios, DaisyUI, HeroIcons

### Database
- Schema managed via Liquibase: `db/changelog-master.yaml` + `db/changelogs/*.yaml`
- Seed data: `tearup/store/init.sql`, `tearup/point/init.sql`

### ATDD Tests
- API tests: `atdd/api/collections/` (Postman/Newman)
- UI tests: `atdd/ui/` (Robot Framework + SeleniumLibrary)

## Naming Conventions

### store-web (TypeScript/React)
- Types & Components: PascalCase (`HomeType`, `Homepage()`)
- Business logic functions: camelCase (`calculateTotalPrice()`)
- HTML element IDs: kebab-case (`receiver-name`, `total-amount`)
- Files: kebab-case (`order-list.ts`)
- Directories: lowercase
- Array variables: append "List" (`orderList`)
- Constants: UPPERCASE (`HOUR`, `MINUTE`)
- No semicolons

### store-service (Go)
- Functions: PascalCase (`CalculateTotalPrice()`)
- Files: PascalCase (`OrderService.go`, `OrderService_test.go`)
- Packages/directories: lowercase
- Test functions: Snake_Case (`Test_CalculateAge_Input_Birth_Date_18042003_Should_be_16`)
- Variables: camelCase, constants UPPERCASE

### Commit Messages
Use prefix tags: `[Created]`, `[Edited]`, `[Added]`, `[Deleted]` — include details about what changed and where.
