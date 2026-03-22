# Deployment

Two-cluster architecture: **app cluster** (runs the application) + **monitoring cluster** (runs the LGTM observability stack). This pattern is used on both k3d (local) and EKS (production).

## Architecture Overview

Uses an **agent-gateway pattern**: all telemetry flows through a local OTel Gateway on the app cluster, which forwards to the monitoring cluster via a single buffered connection.

```
┌──────────────────────────────────────┐
│  App Cluster (sck-workshop)          │
│                                      │
│  nginx-ingress :80                   │
│  ├─ store-web (Next.js)              │
│  ├─ store-service (Go/Gin)           │
│  ├─ point-service (NestJS)           │
│  ├─ thirdparty (mock gateways)      │
│  └─ mysql + liquibase                │
│                                      │
│  store-service ──┐                   │
│  point-service ──┤                   │
│  Beyla sidecar ──┼► OTel Gateway ────┼──┐  OTLP gRPC
│  MySQL sidecar ──┤   (forwarder)    │  │  (batch + retry)
│  node-exporter ──┤   :4317/:4318    │  │
│  kube-state-m  ──┘                   │  │
└──────────────────────────────────────┘  │
                                          │  Cross-cluster network
┌──────────────────────────────────────┐  │
│  Monitoring Cluster (sck-monitoring) │  │
│                                      │  │
│  OTel Collector :4317/:4318  ←───────┼──┘
│  ├─► Tempo      (traces)            │
│  ├─► Loki       (logs)              │
│  ├─► Prometheus (metrics)           │
│  └─► spanmetrics + servicegraph     │
│                                      │
│  Pyroscope :4040  ←──────────────────┼── (direct from services, not via gateway)
│  Grafana :80 (admin/workshop)        │
│  └─ datasources: Tempo, Loki,       │
│     Prometheus, Pyroscope            │
└──────────────────────────────────────┘
```

**Key design:** No processing duplication. The gateway is a lightweight forwarder. All connectors (spanmetrics, servicegraph) and backend routing stay in the monitoring cluster. Signal correlations (trace↔log, trace↔profile, metric→trace exemplars) are preserved — the gateway forwards full OTLP payloads without stripping attributes.

**Why Pyroscope is direct (not via gateway):** The `pyroscope-go` SDK uses Pyroscope's native HTTP push API (port 4040), not OTLP. The OTel Collector has no Pyroscope exporter — there is no way to route profiles through the OTel pipeline. Trace↔profile correlation still works because `otel-profiling-go` injects `pyroscope.profile.id` into spans (which travel through the gateway to Tempo), while profiles are pushed directly to Pyroscope. Grafana's `tracesToProfiles` datasource config links them together.

Cross-cluster connectivity:
- **k3d**: Shared Docker network (`k3d-shared`), gateway forwards via `k3d-sck-monitoring-serverlb`
- **EKS**: VPC Peering + internal NLB on monitoring cluster, gateway uses NLB DNS hostname

---

## k3d (Local Development)

### Prerequisites

- [k3d](https://k3d.io/) v5+
- [Helm](https://helm.sh/) v3+
- Docker

### Quick Start (Full Stack)

For first-time clone — run these in order:

```bash
# 1. Create both clusters on shared network + install nginx-ingress
make k3d_create_all

# 2. Build all Docker images
docker build -t sck/store-database:k3d -f deploy/k8s/app/store-database/Dockerfile.seed .
docker build -t sck/store-service:k3d store-service/
docker build -t sck/point-service:k3d point-service/
docker build -t sck/store-web:k3d store-web/
docker build -t sck/store-shipping-gateway:k3d thirdparty/
docker build -t sck/liquibase:k3d -f db/Dockerfile .
docker build -t sck/mysql-seed:k3d -f deploy/k8s/app/store-database/Dockerfile.seed .

# 3. Import images into app cluster
k3d image import \
  sck/store-database:k3d sck/store-service:k3d sck/point-service:k3d \
  sck/store-web:k3d sck/store-shipping-gateway:k3d sck/liquibase:k3d \
  sck/mysql-seed:k3d \
  -c sck-workshop

# 4. Create namespace and secrets
kubectl create namespace public --context k3d-sck-workshop
kubectl create secret generic db-credentials \
  --from-literal=mysql-root-password=root \
  --from-literal=mysql-user=user \
  --from-literal=mysql-password=password \
  -n public --context k3d-sck-workshop
kubectl create secret generic app-secrets \
  --from-literal=db-connection='user:password@tcp(mysql:3306)/store?parseTime=true' \
  --from-literal=jwt-secret=my-secret-key \
  -n public --context k3d-sck-workshop

# 5. Deploy LGTM monitoring stack
make k3d_deploy_monitoring

# 6. Deploy app with monitoring endpoints auto-configured
make k3d_connect
```

### Verify the Setup

#### Step 1: Check all pods are running

```bash
kubectl get pods -n public --context k3d-sck-workshop
kubectl get pods -n monitoring --context k3d-sck-monitoring
```

Expected: All pods `Running` or `Completed` (seed/migration jobs). Key pods on app cluster:
- `store-service-deployment` (2 replicas)
- `point-service-deployment`
- `thirdparty-deployment` (2/2 — includes Beyla sidecar)
- `mysql` (2/2 — includes OTel Collector sidecar)
- `otel-gateway-opentelemetry-collector` — **this is the new gateway**
- `node-exporter`, `kube-state-metrics`

#### Step 2: Verify endpoints point to local gateway (not remote collector)

```bash
# All should show: http://otel-gateway-opentelemetry-collector.public:4317
kubectl exec -n public --context k3d-sck-workshop deployment/store-service-deployment \
  -- printenv OTEL_EXPORTER_OTLP_ENDPOINT

kubectl exec -n public --context k3d-sck-workshop deployment/point-service-deployment \
  -- printenv OTEL_EXPORTER_OTLP_ENDPOINT

# Pyroscope should point directly to monitoring cluster (not via gateway)
kubectl exec -n public --context k3d-sck-workshop deployment/store-service-deployment \
  -- printenv PYROSCOPE_URL
# Expected: http://k3d-sck-monitoring-serverlb:4040
```

#### Step 3: Generate traffic

```bash
# Hit various endpoints to generate traces, logs, and metrics
for i in $(seq 1 10); do
  curl -s -o /dev/null http://localhost/api/v1/product
  curl -s -o /dev/null -X POST http://localhost/api/v1/login \
    -H 'Content-Type: application/json' -d '{"email":"test@test.com","password":"test"}'
done

# Wait for telemetry pipeline to flush
sleep 20
```

#### Step 4: Verify traces in Tempo

```bash
kubectl port-forward -n monitoring --context k3d-sck-monitoring svc/tempo 3100:3100 &
sleep 2

# Should show traces from store-service, point-service, thirdparty-gateway
curl -s 'http://localhost:3100/api/search?q=%7B%7D&limit=20' | python3 -c "
import sys, json
traces = json.load(sys.stdin).get('traces',[])
services = {}
for t in traces:
    svc = t.get('rootServiceName','?')
    services[svc] = services.get(svc, 0) + 1
print(f'Total traces: {len(traces)}')
for svc, cnt in sorted(services.items()):
    print(f'  {svc}: {cnt} traces')
"

kill %1 2>/dev/null
```

#### Step 5: Verify logs in Loki (with trace correlation)

```bash
kubectl port-forward -n monitoring --context k3d-sck-monitoring svc/loki 3100:3100 &
sleep 2

# Should show logs with trace_id and span_id labels
START=$(date -v-30M +%s)
END=$(date +%s)
curl -s "http://localhost:3100/loki/api/v1/query_range?query=%7Bservice_name%3D%22store-service%22%7D&limit=3&start=${START}000000000&end=${END}000000000" \
  | python3 -c "
import sys, json
results = json.load(sys.stdin).get('data',{}).get('result',[])
for r in results[:3]:
    labels = r.get('stream',{})
    print(f'trace_id={\"trace_id\" in labels} span_id={\"span_id\" in labels} log={r.get(\"values\",[[\"\"]])[0][1][:80]}')
"

kill %1 2>/dev/null
```

#### Step 6: Verify spanmetrics in Prometheus

```bash
kubectl port-forward -n monitoring --context k3d-sck-monitoring svc/prometheus-server 9090:80 &
sleep 2

# Should show spanmetric series for all 3 services
curl -s 'http://localhost:9090/api/v1/query?query=duration_milliseconds_count' | python3 -c "
import sys, json
results = json.load(sys.stdin).get('data',{}).get('result',[])
services = {}
for r in results:
    svc = r.get('metric',{}).get('service_name','?')
    services[svc] = services.get(svc, 0) + 1
print(f'Spanmetric series: {len(results)}')
for svc, cnt in sorted(services.items()):
    print(f'  {svc}: {cnt} series')
"

kill %1 2>/dev/null
```

#### Step 7: Verify profile correlation

```bash
kubectl port-forward -n monitoring --context k3d-sck-monitoring svc/tempo 3100:3100 &
sleep 2

# Get a store-service trace and check for pyroscope.profile.id
TRACE_ID=$(curl -s 'http://localhost:3100/api/search?q=%7Bresource.service.name%3D%22store-service%22%7D&limit=1' \
  | python3 -c "import sys,json; t=json.load(sys.stdin).get('traces',[]); print(t[0]['traceID'] if t else '')")

curl -s "http://localhost:3100/api/traces/$TRACE_ID" | python3 -c "
import sys, json
for batch in json.load(sys.stdin).get('batches',[]):
    for ss in batch.get('scopeSpans',[]):
        for span in ss.get('spans',[]):
            attrs = {a['key']: a.get('value',{}).get('stringValue','') for a in span.get('attributes',[])}
            pid = attrs.get('pyroscope.profile.id','NOT PRESENT')
            if pid != 'NOT PRESENT':
                print(f'span={span[\"name\"]} pyroscope.profile.id={pid}')
"

kill %1 2>/dev/null
```

#### Step 8: Verify via Grafana UI

Open **Grafana** at http://localhost:3000 (admin / workshop):
- **Explore > Tempo** — see traces from store-service, point-service, thirdparty-gateway
- **Explore > Loki** — query `{service_name="store-service"}` — click trace_id to jump to Tempo
- **Explore > Prometheus** — query `duration_milliseconds_count` for spanmetrics (RED metrics)
- **Explore > Pyroscope** — select `store-service` → `process_cpu` → see CPU flame graph
- On a trace span: click **"Logs for this span"** → should show correlated logs
- On a CPU-heavy span (e.g. login/bcrypt): click **"Profiles for this span"** → should show profile

Test users: `user_1` through `user_50`, all with password `P@ssw0rd`.

### Step-by-Step

#### Create clusters

```bash
# Create shared Docker network
docker network create k3d-shared

# App cluster (port 80 for the app)
k3d cluster create sck-workshop \
  --network k3d-shared \
  --port "80:80@loadbalancer" \
  --port "443:443@loadbalancer" \
  --k3s-arg "--disable=traefik@server:0" \
  --agents 1 --wait

# Monitoring cluster (port 3000=Grafana, 4317/4318=OTel, 4040=Pyroscope)
k3d cluster create sck-monitoring \
  --network k3d-shared \
  --port "3000:80@loadbalancer" \
  --port "4317:4317@loadbalancer" \
  --port "4318:4318@loadbalancer" \
  --port "4040:4040@loadbalancer" \
  --k3s-arg "--disable=traefik@server:0" \
  --agents 0 --wait
```

#### Install nginx-ingress on app cluster

```bash
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace \
  --kube-context k3d-sck-workshop \
  --set controller.ingressClassResource.name=public \
  --set controller.service.type=LoadBalancer \
  --wait
```

#### Deploy monitoring stack

```bash
# Installs Loki, Tempo, Prometheus, Pyroscope, OTel Collector, Grafana
./deploy/k8s/k3d-monitoring.sh
```

This script also installs nginx-ingress on the monitoring cluster with TCP passthrough for OTel (4317/4318) and Pyroscope (4040).

#### Deploy app with monitoring (agent-gateway pattern)

```bash
# Deploy app + OTel Gateway + monitoring agents
# OTEL_ENDPOINT tells the gateway where to forward (monitoring cluster)
# Services send to the local gateway, NOT directly to the monitoring cluster
OTEL_ENDPOINT=http://k3d-sck-monitoring-serverlb:4317 \
PYROSCOPE_ENDPOINT=http://k3d-sck-monitoring-serverlb:4040 \
./deploy/k8s/k3d-deploy.sh
```

Or use the Makefile shortcut: `make k3d_connect`

This deploys:
- All app services (store-service, point-service, store-web, thirdparty, mysql)
- **OTel Gateway** (Helm release `otel-gateway`) — receives OTLP from all local services + scrapes node-exporter/kube-state-metrics, forwards everything to monitoring cluster with buffering/retry
- Beyla sidecar on thirdparty, OTel Collector sidecar on MySQL — both send to the local gateway
- monitoring-endpoints ConfigMap pointing services to the local gateway (not remote collector)

#### Deploy app without monitoring

```bash
# Uses placeholder endpoints — OTel SDK fails silently
./deploy/k8s/k3d-deploy.sh
```

### Access Points (k3d)

| Service | URL |
|---------|-----|
| App | http://localhost |
| Grafana | http://localhost:3000 (admin/workshop) |
| OTel Collector (gRPC) | localhost:4317 |
| OTel Collector (HTTP) | localhost:4318 |
| Pyroscope | localhost:4040 |

### Teardown

```bash
make k3d_delete_all
# Deletes both clusters + shared Docker network
```

### Makefile Targets

| Target | Description |
|--------|-------------|
| `k3d_create_all` | Create shared network + both clusters + nginx-ingress on app cluster |
| `k3d_deploy_monitoring` | Deploy LGTM stack to monitoring cluster |
| `k3d_deploy` | Deploy app only (no monitoring endpoints) |
| `k3d_connect` | Deploy app with monitoring endpoints auto-configured |
| `k3d_delete_all` | Delete both clusters + shared network |

---

## EKS (Production)

### Prerequisites

- AWS CLI configured with appropriate credentials
- Terraform >= 1.5
- kubectl
- Helm v3+

### Infrastructure Setup (Terraform)

All infrastructure is managed in `deploy/terraform/`:

```bash
cd deploy/terraform

# Initialize Terraform
terraform init

# Plan — review what will be created
terraform plan

# Apply — creates both EKS clusters, VPC peering, LGTM stack
terraform apply
```

#### What Terraform Creates

**App cluster** (`sck-workshop`):
- VPC with public/private subnets (10.0.0.0/16)
- EKS cluster (Kubernetes 1.33) with managed node groups (t3.medium)
- Ingress-nginx controller with internet-facing NLB
- EBS CSI driver for persistent volumes

**Monitoring cluster** (`sck-monitoring`):
- VPC with public/private subnets (10.1.0.0/16)
- EKS cluster with managed node groups (t3.medium)
- Full LGTM stack via Helm:
  - **Loki** (6.16.0) — log aggregation
  - **Tempo** (1.10.3) — distributed tracing with metrics generator
  - **Prometheus** (25.27.0) — metrics with remote-write receiver + exemplar storage
  - **Pyroscope** (1.7.1) — continuous profiling (internal NLB)
  - **OTel Collector** (0.97.1) — telemetry pipeline with spanmetrics + servicegraph connectors (internal NLB)
  - **Grafana** (10.5.15) — dashboards + pre-configured datasources (internet-facing NLB)

**Cross-cluster networking**:
- VPC Peering connection between app VPC and monitoring VPC
- Route table entries for bidirectional traffic
- Security group rules allowing OTLP (4317-4318) and all TCP from app VPC

#### Terraform Files

| File | Purpose |
|------|---------|
| `main.tf` | Providers (AWS, Helm, Kubernetes) for both clusters |
| `variables.tf` | Cluster names, instance types, node counts |
| `app-vpc.tf` | App cluster VPC (10.0.0.0/16) |
| `app-eks.tf` | App EKS cluster + node groups |
| `app-ingress.tf` | Ingress-nginx controller for app cluster |
| `monitoring-vpc.tf` | Monitoring VPC (10.1.0.0/16) |
| `monitoring-eks.tf` | Monitoring EKS cluster + node groups |
| `monitoring-lgtm-stack.tf` | All Helm releases: Loki, Tempo, Prometheus, Pyroscope, OTel Collector, Grafana |
| `monitoring-workshop-agents.tf` | Node-exporter, kube-state-metrics, OTel scraper on app cluster → pushes to monitoring |
| `vpc-peering.tf` | VPC peering + routes + security group rules |
| `outputs.tf` | Cluster endpoints, kubectl config commands |

### Configure kubectl

```bash
# App cluster
aws eks update-kubeconfig --region ap-southeast-7 --name sck-workshop

# Monitoring cluster
aws eks update-kubeconfig --region ap-southeast-7 --name sck-monitoring
```

### Deploy App to EKS

```bash
# Build + push + deploy (auto-generates unique image tag)
make eks_deploy_all

# Or deploy individual services
make eks_deploy_store
make eks_deploy_point

# Deploy all app manifests (without rebuilding images)
make eks_deploy_app

# Add monitoring overlays (OTel sidecar, Beyla)
make eks_deploy_monitoring

# Both app + monitoring overlays
make eks_deploy_full
```

### Get Endpoints (EKS)

```bash
# App ingress URL
kubectl -n ingress-nginx get svc ingress-nginx-controller \
  -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'

# Grafana URL
kubectl --context arn:aws:eks:ap-southeast-7:ACCOUNT_ID:cluster/sck-monitoring \
  -n monitoring get svc grafana \
  -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'

# OTel Collector internal endpoint (for monitoring-endpoints ConfigMap)
kubectl --context arn:aws:eks:ap-southeast-7:ACCOUNT_ID:cluster/sck-monitoring \
  -n monitoring get svc otel-collector-opentelemetry-collector \
  -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
```

### Cluster Contexts

| Cluster | Context |
|---------|---------|
| App | `arn:aws:eks:ap-southeast-7:517425940836:cluster/sck-workshop` |
| Monitoring | `arn:aws:eks:ap-southeast-7:517425940836:cluster/sck-monitoring` |

---

## App Deployment (K8s Manifests)

2-step deployment: **app** (standalone) and **monitoring** (optional overlay).

### Step 1: App Only

```bash
make eks_deploy_app

# Or manually
kubectl apply -f deploy/k8s/app/monitoring-endpoints.yml
kubectl apply -f deploy/k8s/app/store-database/service.yml
kubectl apply -f deploy/k8s/app/store-service/service.yml
kubectl apply -f deploy/k8s/app/point-service/service.yml
kubectl apply -f deploy/k8s/app/store-web/service.yml
kubectl apply -f deploy/k8s/app/liquibase/job.yml
kubectl apply -f deploy/k8s/app/ingress.yml
```

The `monitoring-endpoints.yml` ConfigMap is **required** even without monitoring. Services reference it via `envFrom`. Placeholder values are fine — OTel SDK fails silently when the endpoint is unreachable.

### Step 2: Monitoring Overlays

```bash
make eks_deploy_monitoring

# Or manually
kubectl apply -f deploy/k8s/monitoring/store-database-with-otel.yml
kubectl apply -f deploy/k8s/monitoring/thirdparty-with-beyla.yml
```

Before applying, replace `REPLACE_OTEL_ENDPOINT` placeholders in the monitoring manifests with your actual OTel collector endpoint.

---

## Directory Structure

```
deploy/
├── k8s/
│   ├── app/                           # App manifests (both k3d and EKS)
│   │   ├── monitoring-endpoints.yml   # ConfigMap (OTEL_EXPORTER_OTLP_ENDPOINT, PYROSCOPE_URL)
│   │   ├── ingress.yml
│   │   ├── store-service/service.yml
│   │   ├── point-service/service.yml
│   │   ├── store-web/service.yml
│   │   ├── store-database/service.yml
│   │   ├── thirdparty/service.yml
│   │   └── liquibase/job.yml
│   ├── monitoring/                    # EKS monitoring overlays (OTel sidecar, Beyla)
│   │   ├── store-database-with-otel.yml
│   │   └── thirdparty-with-beyla.yml
│   ├── k3d-deploy.sh                 # Deploy app to k3d app cluster
│   └── k3d-monitoring.sh             # Deploy LGTM stack to k3d monitoring cluster
└── terraform/                         # EKS infrastructure (both clusters)
    ├── main.tf, variables.tf, outputs.tf
    ├── app-vpc.tf, app-eks.tf, app-ingress.tf
    ├── monitoring-vpc.tf, monitoring-eks.tf
    ├── monitoring-lgtm-stack.tf
    ├── monitoring-workshop-agents.tf
    └── vpc-peering.tf
```

## Jenkins

`Jenkinsfile-Deploy` has a `DEPLOY_MONITORING` boolean parameter (default: false). When enabled, it applies the monitoring overlay manifests after the app deployment.
