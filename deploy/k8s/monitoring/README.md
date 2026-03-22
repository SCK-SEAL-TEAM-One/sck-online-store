# Monitoring Overlays

Optional manifests that add observability sidecars to app services. Apply these **after** the base app manifests in `deploy/k8s/app/`.

## Files

| File | Replaces | What it adds |
|------|----------|-------------|
| `store-database-with-otel.yml` | `app/store-database/service.yml` | OTel Collector sidecar (MySQL metrics + slow query logs via HTTP/4318) |
| `thirdparty-with-beyla.yml` | `app/thirdparty/service.yml` | Grafana Beyla sidecar (auto-instrumentation via eBPF, gRPC/4317) |

## Setup

1. Replace `REPLACE_OTEL_ENDPOINT` in each file with your OTel collector address
2. Apply: `kubectl apply -f deploy/k8s/monitoring/`

## How Overlays Work

These manifests use the **same resource names** as the app versions (e.g., `Deployment/mysql`, `Deployment/thirdparty-deployment`). When you `kubectl apply` them, Kubernetes updates the existing resources in-place — adding the sidecar containers.

To revert to app-only, re-apply the base manifest:
```bash
kubectl apply -f deploy/k8s/app/store-database/service.yml
kubectl apply -f deploy/k8s/app/thirdparty/service.yml
```

## Maintenance

When updating the base MySQL or thirdparty config (env vars, probes, image tags), update **both** the app and monitoring versions to keep them in sync.
