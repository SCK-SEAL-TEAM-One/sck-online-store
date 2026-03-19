# Known Issue: "Order completed" Log Missing on EKS

**Status:** Unresolved — works locally, fails on EKS
**Date:** 2026-03-19
**Branch:** observability-workshop

## Problem

`slog.InfoContext(ctx, "Order completed", ...)` in `store-service/internal/order/order.go:181` never reaches the OTel collector on EKS. The other log `slog.InfoContext(ctx, "Order created", ...)` in `store-service/cmd/api/order.go:79` works fine on EKS. Both work locally in docker-compose.

The HTTP response returns 200 with a valid order_number, proving both log lines execute. The OTel collector debug exporter confirmed it receives "Order created" but NOT "Order completed" — the issue is in the app's OTel SDK, not the collector or Loki.

## Key Context for Debugging

**Execution order:** "Order completed" fires FIRST (inside `CreateOrder()`), then "Order created" fires SECOND (in the handler after return). Both use the same `ctx` and same `slog` handler.

**Local environment:** App sends directly to LGTM container (Loki 3.6.0) on port 4317 via gRPC.
**EKS environment:** App sends to OTel collector (via NLB) on port 4317 via gRPC, collector forwards to Loki 3.1.x.

**"Order completed" has 26 Loki labels locally** (9-field metadata map). "Order created" has fewer labels (5-field metadata map).

## What Has Been Ruled Out

- **Binary mismatch** — Confirmed deployed binary contains the code (rebuilt with `--no-cache`)
- **Loki label limits** — Increased `max_label_names_per_series` to 30, didn't help (26 < 30)
- **Batch queue overflow** — Set `OTEL_BLRP_MAX_QUEUE_SIZE=16384` + `OTEL_BLRP_SCHEDULE_DELAY=200`, didn't help
- **Nested map attributes** — Flattened `slog.Any("metadata", map)` to flat key-values, didn't help
- **Collector filtering** — Logs pipeline has no processors/filters, just `otlp → otlphttp/loki`

## What to Try Next

**IMPORTANT: Try ONE thing at a time. Test locally first, then deploy to EKS. Use `make eks_deploy_store` for EKS builds.**

### 1. Switch protocol to HTTP (env-only change, no code change)
The local LGTM uses gRPC directly. EKS also uses gRPC but through an NLB. Try HTTP to rule out gRPC-specific issues.
```yaml
# In deploy/k8s/store-service/service.yml
- name: OTEL_EXPORTER_OTLP_PROTOCOL
  value: http/protobuf
- name: OTEL_EXPORTER_OTLP_ENDPOINT
  value: http://<collector-endpoint>:4318
```

### 2. Add console logging alongside OTel (code change)
Add a multi-handler in `store-service/internal/otel/otel.go` to write logs to both stdout and OTel. This will confirm whether slog is actually emitting the record.
```go
multiHandler := slogmulti.Fanout(
    slog.NewJSONHandler(os.Stdout, nil),
    otelslog.NewHandler("store-service", otelslog.WithLoggerProvider(lp)),
)
slog.SetDefault(slog.New(multiHandler))
```

### 3. Use SimpleProcessor instead of BatchProcessor (code change)
Replace async batching with synchronous export to rule out timing/buffering issues.
```go
// In otel.go, change:
lp := otellog.NewLoggerProvider(
    otellog.WithProcessor(otellog.NewSimpleProcessor(logExporter)),
    otellog.WithResource(res),
)
```

### 4. Reduce metadata map size
"Order completed" has 9 fields in its metadata map vs 5 in "Order created". Try reducing to match and see if size matters.

## Files Involved

- `store-service/internal/order/order.go:181-198` — "Order completed" log (broken on EKS)
- `store-service/cmd/api/order.go:79-92` — "Order created" log (works on EKS)
- `store-service/internal/otel/otel.go` — OTel init with BatchProcessor
- `deploy/k8s/store-service/service.yml` — K8s deployment manifest
- `deploy/terraform/lgtm-stack.tf` — Loki + OTel collector config
