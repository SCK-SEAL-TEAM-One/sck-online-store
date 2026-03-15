# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this directory.

## Overview

Next.js 14 (App Router) e-commerce frontend with React 18, TailwindCSS + DaisyUI, Zustand state management. Runs on port 3000.

## Commands

```bash
npm install          # Install dependencies
npm run dev          # Dev server
npm run build        # Production build
npm run lint         # ESLint
npm run format       # Prettier formatting

# Cypress component tests
npm run cy:open      # Interactive mode
npm run cy:run       # Headless mode
npm run test:component  # Run component test specs
```

## Architecture

```
src/
  app/                — Next.js App Router pages
    auth/login/       — Login page + components
    product/list/     — Product listing with search
    product/[id]/     — Product detail (dynamic route)
    cart/             — Shopping cart view
    checkout/         — Checkout with shipping/payment
    payment/          — Payment gateway (OTP flow)
    orders/completed/ — Order confirmation
    layout.tsx        — Root layout (wraps with HydrationZustand)
  components/         — Shared UI components (button, badge, input-field, etc.)
  layouts/common/     — Header, nav menu, cart offcanvas
  services/           — API client layer (Axios)
    auth.ts           — Login, token refresh
    product-*.ts      — Product list/detail
    cart/             — Cart operations (get, add, update)
    order-checkout.ts — Submit order
    confirm-payment.ts — Payment confirmation
    point.ts          — Points operations
    download-pdf.ts   — Invoice PDF download
  hooks/              — Zustand stores
    use-user-store.ts — User state (persisted to localStorage)
    use-order-store.ts — Cart/checkout state (devtools, uses Immer)
  utils/              — Helpers
    axios.ts          — Axios instance with auth interceptor + token refresh queue
    total-price.ts    — Price calculations
    shipping.ts       — Shipping fee calculations
    point.ts          — Point calculations
    credit-card-format.ts — Card validation
  config/             — App constants (image URLs, logo paths, pointRate)
  __test__/           — Cypress component tests
```

## State Management (Zustand)

**useUserStore** — Persisted user session (userId, name, username). Uses `persist` + `devtools` middleware.

**useOrderStore** — Checkout flow state (cart items, shipping, payment, points, summary). Uses `devtools` + Immer for immutable updates. Not persisted.

**HydrationZustand** — Wrapper component in root layout that prevents SSR hydration mismatch by waiting for Zustand store hydration.

## API Client Pattern

All services use a shared Axios instance (`utils/axios.ts`) that:
- Adds Bearer token from localStorage on every request
- Handles 401 responses with automatic token refresh
- Queues concurrent requests during refresh to prevent duplicate refreshes
- Redirects to login on auth failure

## Naming Conventions

- Types & Components: PascalCase (`HomeType`, `Homepage()`)
- Business logic: camelCase (`calculateTotalPrice()`)
- HTML element IDs: kebab-case (`receiver-name`, `total-amount`)
- Files: kebab-case (`order-list.ts`)
- Directories: lowercase
- Array variables: append "List" (`orderList`)
- Constants: UPPERCASE
- No semicolons

## Key Dependencies

Next.js 14, React 18, Zustand, Axios, TailwindCSS 3.4, DaisyUI 4.6, Headless UI, HeroIcons, Immer, dayjs, Cypress
