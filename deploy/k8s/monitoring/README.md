# Monitoring Overlays

Optional manifests that add observability sidecars to app services. Apply these **after** the base app manifests in `deploy/k8s/app/`.

## Architecture: Agent-Gateway Pattern

All telemetry from app services flows through a local **OTel Gateway** on the app cluster, which forwards to the monitoring cluster's OTel Collector. This avoids multiple direct cross-cluster connections and provides local buffering/retry.

```
App Cluster                           Monitoring Cluster
  services ──► OTel Gateway ──OTLP──► OTel Collector ──► Tempo/Loki/Prometheus
  (sidecars)   (forwarder)            (spanmetrics, servicegraph, fan-out)
```

The gateway is deployed by:
- **k3d:** `k3d-deploy.sh` (Helm release `otel-gateway`)
- **EKS:** Terraform `monitoring-workshop-agents.tf` (Helm release `otel-gateway`)

## Files

| File | Replaces | What it adds |
|------|----------|-------------|
| `store-database-with-otel.yml` | `app/store-database/service.yml` | OTel Collector sidecar (MySQL metrics + slow query logs via HTTP/4318) |
| `thirdparty-with-beyla.yml` | `app/thirdparty/service.yml` | Grafana Beyla sidecar (auto-instrumentation via eBPF, gRPC/4317) |

## Setup

1. `REPLACE_OTEL_ENDPOINT` placeholders are substituted at deploy time with the local OTel Gateway address (e.g., `otel-gateway-opentelemetry-collector.monitoring`)
2. Use `make eks_deploy_monitoring` (EKS) or `make k3d_connect` (k3d) — both handle substitution automatically

## How Overlays Work

These manifests use the **same resource names** as the app versions (e.g., `Deployment/mysql`, `Deployment/thirdparty-deployment`). When you `kubectl apply` them, Kubernetes updates the existing resources in-place — adding the sidecar containers.

To revert to app-only, re-apply the base manifest:
```bash
kubectl apply -f deploy/k8s/app/store-database/service.yml
kubectl apply -f deploy/k8s/app/thirdparty/service.yml
```

## Maintenance

When updating the base MySQL or thirdparty config (env vars, probes, image tags), update **both** the app and monitoring versions to keep them in sync.
