# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this directory.

## Overview

Acceptance Test-Driven Development (ATDD) suite with API tests (Newman/Postman) and UI tests (Robot Framework/Selenium).

## Commands (from project root)

```bash
# Start all services for testing
make start_test_suite           # Local browser
make start_test_suite_grid      # With Selenium Grid (headless)

# API tests (Newman)
make run_newman                    # All API test suites
make run_newman_authentication     # Auth tests only
make run_newman_order_summary_pdf  # Order tests only

# UI tests (Robot Framework)
make run_robot                     # All UI test suites
make run_robot_authentication      # Auth tests only
make run_robot_order_summary_pdf   # Order tests only

# Stop
make stop_test_suite
```

## API Tests (`api/`)

```
api/
  collections/
    001-Authentication.postman_collection.json
    002-Order-Summary-PDF.postman_collection.json
  data/
    001-Authentication/    — TSS-AUTH-001..003, TSA-AUTH-001..003 (JSON test data)
    002-Order-Summary-PDF/ — TSS-OSP-001, TSS-OSP-002
  sck-online-store.local.postman_environment.json   — Local env (localhost)
  sck-online-store.remote.postman_environment.json  — Remote env
```

Data-driven via JSON files containing user credentials, product info, pricing, shipping details, payment info, and expected results. Newman runs with `cli,junit,htmlextra` reporters.

## UI Tests (`ui/`)

```
ui/
  001-Authentication/
    TSS-AUTH-001-Login_first_time_success.robot
    TSS-AUTH-002-Logged_in_and_Re_enter_success.robot
    TSA-AUTH-001-Login_with_incorrect_username.robot
  002-Order-Summary-PDF/
    TSS-OSP-001-Order_one_product_one_unit_success.robot
    TSS-OSP-002-Order_two_products_suscess.robot
    DownloadHelper.py  — Custom library for PDF download via Chrome DevTools Protocol
  requirements.txt     — Python dependencies (robotframework, seleniumlibrary, etc.)
```

**Key variables:** `${URL}` (app URL), `${BROWSER}` (headlesschrome), `${REMOTE_HUB_URL}` (Selenium Grid hub).

Robot keywords are written in Thai for readability. Element locators use kebab-case IDs matching the frontend (e.g., `login-username-input`, `product-card-name-{id}`, `shipping-form-*-input`).

## Test Naming Convention

- `TSS-` prefix = success scenario
- `TSA-` prefix = failure/alternative scenario
- Numbered by feature area: `001-Authentication`, `002-Order-Summary-PDF`
