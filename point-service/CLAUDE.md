# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this directory.

## Overview

NestJS microservice for loyalty points management. Uses TypeORM with MySQL. Runs on port 8001. API prefix: `/api/v1`.

## Commands

```bash
npm install            # Install dependencies
npm run start:dev      # Dev mode (watch)
npm run start:debug    # Debug mode
npm run build          # Compile TypeScript
npm run start:prod     # Production (compiled JS)

# Testing
npm test               # Unit tests (Jest)
npm run test:watch     # Watch mode
npm run test:cov       # Coverage report
npm run test:e2e       # End-to-end tests

# Code quality
npm run lint           # ESLint (with auto-fix)
npm run format         # Prettier
```

## Architecture

```
src/
  main.ts             — Entry point (starts OTEL SDK, then NestJS app on :8001)
  app.module.ts       — Root module (ConfigModule, TypeORM, HelloModule, PointModule)
  trace.ts            — OpenTelemetry SDK setup (OTLP gRPC exporter to lgtm:4317)
  point/
    point.module.ts   — Module (imports TypeORM for Point entity)
    point.controller.ts — GET /point, POST /point
    point.service.ts  — getPoint(), deductPoint()
    point.entity.ts   — TypeORM entity (points table: id, orgId, userId, amount, created, updated)
    point.dto.ts      — CreatePointDto (orgId, userId, amount)
    test/
      point.service.spec.ts    — Service unit tests (mocked repository)
      point.controller.spec.ts — Controller unit tests (mocked service)
  hello/              — Example module
```

## Database

TypeORM with MySQL. `synchronize: true` (auto-creates schema from entities).

Environment variables (`.env.dev`): `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`. Database name: `point`.

## Testing Pattern

Uses `@nestjs/testing` TestingModule. Mock dependencies with `jest.fn()`. Arrange-Act-Assert pattern. Test files: `*.spec.ts` in `test/` subdirectories.

## Observability

OpenTelemetry auto-instrumentation for HTTP, MySQL2, and NestJS. OTLP gRPC exporter sends traces to `lgtm:4317`. SDK initializes before NestJS app in `main.ts`.
