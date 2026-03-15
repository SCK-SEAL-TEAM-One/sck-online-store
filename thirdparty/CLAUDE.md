# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this directory.

## Overview

Mock external services using MounteBank (imposters). Provides fake payment and shipping gateways for development and testing.

## Services

### Bank Gateway (port 8882)

| Method | Path | Response |
|---|---|---|
| POST | `/payment/visa` | `{ status: "completed", payment_date, transaction_id }` (dynamic) |
| GET | `/card/information` | Static test card data (Visa 4719700591590995) |

`response.js` generates random transaction IDs and timestamps.

### Shipping Gateway (port 8883)

| Method | Path | Response |
|---|---|---|
| POST | `/shipping` | `{ tracking_number: "{PREFIX}-{RANDOM}" }` (dynamic) |

Reads `shipping_method_id` from request body. Prefix mapping: 1→KR (Kerry), 2→TH (Thai Post), 3→LM (Lineman).

## Structure

```
imposters.json           — Root config (includes both gateways)
bank-gateway/
  imposters.json         — Port 8882 stubs
  response.js            — Dynamic payment response
shipping-gateway/
  imposters.json         — Port 8883 stubs
  response.js            — Dynamic shipping response
```

Unmatched requests return HTTP 400.
