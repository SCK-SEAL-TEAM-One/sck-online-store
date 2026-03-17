# Load Testing — Shopping Flow

Performance tests for the full e-commerce shopping workflow. Two implementations available: **k6** (recommended for CI) and **JMeter** (for GUI-based analysis).

Both tools execute the same 8-step user journey with realistic think time (1-3s random delay between steps).

## Test Flow

```
1. Login              POST   /api/v1/login
2. Browse Products    GET    /api/v1/product?q=&offset=0&limit=20
3. View Product       GET    /api/v1/product/{id}         (random from [1,2,3,4,5,6,8,9])
4. Add to Cart        PUT    /api/v1/addCart
5. Get Cart           GET    /api/v1/cart
6. Submit Order       POST   /api/v1/order
7. Confirm Payment    POST   /api/v1/confirmPayment
8. Order Summary      POST   /api/v1/order/{id}/summary   (Accept: application/json)
```

Product ID 7 is intentionally skipped to simulate fault injection scenarios.

## Prerequisites

- All services running (via `make start_all` or deployed to EKS)
- Database seeded with test users (`user_1` through `user_84`, password: `P@ssw0rd`)
- Reset database before each run for clean results

```bash
# Local: reset via Docker
make down && make start_all

# EKS: reset MySQL pod and re-run Liquibase migration
kubectl -n public delete pod -l app=mysql
kubectl -n public delete job liquibase-migration
kubectl -n public apply -f deploy/k8s/liquibase/job.yml
kubectl -n public rollout restart deployment store-service point-service
```

## k6

### Install

```bash
# macOS
brew install k6

# Docker (no install needed)
docker run --rm -i grafana/k6 version
```

### Run

```bash
# Against local (default: http://localhost)
k6 run atdd/load-test/shopping-flow.js

# Against remote
k6 run -e BASE_URL=http://<elb-hostname> atdd/load-test/shopping-flow.js

# Custom VU count (override stages)
k6 run --vus 50 --duration 3m -e BASE_URL=http://localhost atdd/load-test/shopping-flow.js
```

### Configuration

| Setting | Value | Description |
|---------|-------|-------------|
| Virtual Users | 20 | Ramp up over 30s, sustain 2m30s, ramp down 30s |
| Total Duration | ~3m30s | Across all stages |
| User Pool | `user_1`..`user_84` | Distributed via `VU % 84 + 1` |
| Think Time | 1-3s random | Between each step |

### Thresholds

| Metric | Threshold |
|--------|-----------|
| `http_req_duration` | p(99) < 200ms |
| `http_req_failed` | rate == 0 (zero failures) |

### Output to Grafana (k6 + Prometheus)

```bash
K6_PROMETHEUS_RW_SERVER_URL=http://prometheus:9090/api/v1/write \
  k6 run -o experimental-prometheus-rw atdd/load-test/shopping-flow.js
```

## JMeter

### Install

```bash
# macOS
brew install jmeter

# Or download from https://jmeter.apache.org/download_jmeter.cgi
```

### Run (CLI — non-GUI mode)

```bash
# Basic run with HTML report
jmeter -n \
  -t atdd/load-test/shopping-flow.jmx \
  -l atdd/load-test/results.jtl \
  -e -o atdd/load-test/report/

# With timestamped output (avoid overwrite)
TIMESTAMP=$(date +%Y%m%d%H%M%S) && \
jmeter -n \
  -t atdd/load-test/shopping-flow.jmx \
  -l "atdd/load-test/results-${TIMESTAMP}.jtl" \
  -e -o "atdd/load-test/report-${TIMESTAMP}/"
```

### Run (GUI — for debugging)

```bash
jmeter -t atdd/load-test/shopping-flow.jmx
```

### Configuration

Edit the `HTTP Request Defaults` element in the `.jmx` to change the target:

| Property | Default | Description |
|----------|---------|-------------|
| Server Name | `localhost` | Target host |
| Port | `80` | Target port |
| Protocol | `http` | HTTP or HTTPS |
| Thread Count | `1` | Number of concurrent threads (increase for load) |
| Duration | `210s` | Total test duration |

### Change Target Host

In GUI mode, update **HTTP Request Defaults** > Server Name.

In CLI mode, override with JMeter properties:

```bash
jmeter -n \
  -t atdd/load-test/shopping-flow.jmx \
  -Jhost=<elb-hostname> \
  -Jport=80 \
  -Jthreads=20 \
  -l atdd/load-test/results.jtl
```

> Note: Property overrides (`-J`) require the `.jmx` to use `${__P(host,localhost)}` syntax.
> The current `.jmx` uses hardcoded defaults — update directly in the file or GUI if needed.

## Comparison

| Aspect | k6 | JMeter |
|--------|------|--------|
| Best for | CI pipelines, scripted runs | GUI exploration, detailed analysis |
| Config | Code (JavaScript) | XML (GUI editor) |
| Parameterization | `BASE_URL` env var | HTTP Request Defaults |
| HTML Report | Built-in summary | `--report` flag generates full dashboard |
| Cloud Integration | Grafana Cloud k6 | BlazeMeter, Azure Load Testing |
| Resource Usage | Low (Go-based) | Higher (JVM-based) |

## File Structure

```
atdd/load-test/
  shopping-flow.js    — k6 load test script
  shopping-flow.jmx   — JMeter test plan
  README.md           — This file
  results*.jtl        — JMeter results (gitignored)
  report*/            — JMeter HTML reports (gitignored)
```
