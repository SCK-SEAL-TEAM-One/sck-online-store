#!/bin/bash
# Deploy LGTM monitoring stack to k3d-sck-monitoring cluster
# Usage: ./deploy/k8s/k3d-monitoring.sh
#
# Prerequisites:
#   k3d cluster 'sck-monitoring' running on k3d-shared network
#   Helm 3 installed
#
# Deploys: Loki, Tempo, Prometheus, Pyroscope, OTel Collector, Grafana
# Matches EKS Terraform config in deploy/terraform/monitoring-lgtm-stack.tf

set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CONTEXT="k3d-sck-monitoring"
NAMESPACE="monitoring"

echo "=== Deploying LGTM stack to $CONTEXT (namespace: $NAMESPACE) ==="

# Add Helm repos
echo "Adding Helm repositories..."
helm repo add grafana https://grafana.github.io/helm-charts 2>/dev/null || true
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts 2>/dev/null || true
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts 2>/dev/null || true
helm repo update

# Create namespace
kubectl create namespace "$NAMESPACE" --context "$CONTEXT" --dry-run=client -o yaml | \
  kubectl apply --context "$CONTEXT" -f -

# --- Install nginx-ingress with TCP passthrough for OTel + Pyroscope ---
echo ""
echo "--- Installing nginx-ingress (class: monitoring) ---"
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace \
  --kube-context "$CONTEXT" \
  --set controller.ingressClassResource.name=monitoring \
  --set controller.ingressClassResource.controllerValue=k8s.io/ingress-nginx-monitoring \
  --set controller.service.type=LoadBalancer \
  --set tcp.4317="$NAMESPACE/otel-collector-opentelemetry-collector:4317" \
  --set tcp.4318="$NAMESPACE/otel-collector-opentelemetry-collector:4318" \
  --set tcp.4040="$NAMESPACE/pyroscope:4040" \
  --wait --timeout 120s

# --- Loki (log storage) ---
echo ""
echo "--- Installing Loki ---"
helm upgrade --install loki grafana/loki \
  --version 6.16.0 \
  --namespace "$NAMESPACE" \
  --kube-context "$CONTEXT" \
  --values - <<'EOF'
deploymentMode: SingleBinary
loki:
  auth_enabled: false
  commonConfig:
    replication_factor: 1
  storage:
    type: filesystem
  schemaConfig:
    configs:
      - from: "2024-01-01"
        store: tsdb
        object_store: filesystem
        schema: v13
        index:
          prefix: loki_index_
          period: 24h
  limits_config:
    allow_structured_metadata: true
singleBinary:
  replicas: 1
  persistence:
    size: 2Gi
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      memory: 512Mi
gateway:
  enabled: false
chunksCache:
  enabled: false
resultsCache:
  enabled: false
backend:
  replicas: 0
read:
  replicas: 0
write:
  replicas: 0
EOF

# --- Tempo (trace storage) ---
echo ""
echo "--- Installing Tempo ---"
helm upgrade --install tempo grafana/tempo \
  --version 1.10.3 \
  --namespace "$NAMESPACE" \
  --kube-context "$CONTEXT" \
  --values - <<'EOF'
tempo:
  memBallastSizeMbs: 256
  resources:
    requests:
      cpu: 100m
      memory: 512Mi
    limits:
      memory: 2Gi
config: |
  multitenancy_enabled: false
  usage_report:
    reporting_enabled: false
  compactor:
    compaction:
      block_retention: 24h
  distributor:
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317
          http:
            endpoint: 0.0.0.0:4318
  ingester: {}
  server:
    http_listen_port: 3100
  storage:
    trace:
      backend: local
      local:
        path: /var/tempo/traces
      wal:
        path: /var/tempo/wal
  querier: {}
  query_frontend: {}
  overrides:
    per_tenant_override_config: /conf/overrides.yaml
    metrics_generator_processors:
      - local-blocks
      - service-graphs
      - span-metrics
  metrics_generator:
    storage:
      path: /tmp/tempo
      remote_write:
        - url: http://prometheus-server:80/api/v1/write
          send_exemplars: true
    traces_storage:
      path: /var/tempo/generator/traces
    processor:
      local_blocks:
        filter_server_spans: false
EOF

# --- Prometheus (metrics) ---
echo ""
echo "--- Installing Prometheus ---"
helm upgrade --install prometheus prometheus-community/prometheus \
  --version 25.27.0 \
  --namespace "$NAMESPACE" \
  --kube-context "$CONTEXT" \
  --values - <<'EOF'
server:
  extraFlags:
    - web.enable-remote-write-receiver
    - enable-feature=exemplar-storage
  persistentVolume:
    size: 2Gi
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      memory: 512Mi
alertmanager:
  enabled: false
kube-state-metrics:
  enabled: false
prometheus-node-exporter:
  enabled: false
prometheus-pushgateway:
  enabled: false
EOF

# --- Pyroscope (continuous profiling) ---
echo ""
echo "--- Installing Pyroscope ---"
helm upgrade --install pyroscope grafana/pyroscope \
  --version 1.7.1 \
  --namespace "$NAMESPACE" \
  --kube-context "$CONTEXT" \
  --values - <<'EOF'
pyroscope:
  extraArgs:
    store-gateway.sharding-ring.replication-factor: "1"
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      memory: 512Mi
  persistence:
    enabled: true
    size: 2Gi
alloy:
  enabled: false
minio:
  enabled: false
EOF

# --- OTel Collector (telemetry pipeline) ---
echo ""
echo "--- Installing OTel Collector ---"
helm upgrade --install otel-collector open-telemetry/opentelemetry-collector \
  --version 0.97.1 \
  --namespace "$NAMESPACE" \
  --kube-context "$CONTEXT" \
  --values - <<'EOF'
mode: deployment
replicaCount: 1
image:
  repository: otel/opentelemetry-collector-contrib
ports:
  otlp:
    enabled: true
    containerPort: 4317
    servicePort: 4317
    protocol: TCP
  otlp-http:
    enabled: true
    containerPort: 4318
    servicePort: 4318
    protocol: TCP
  jaeger-compact:
    enabled: false
  jaeger-thrift:
    enabled: false
  jaeger-grpc:
    enabled: false
  zipkin:
    enabled: false
config:
  receivers:
    otlp:
      protocols:
        grpc:
          endpoint: 0.0.0.0:4317
        http:
          endpoint: 0.0.0.0:4318
  connectors:
    spanmetrics:
      histogram:
        explicit:
          buckets:
            - 5ms
            - 10ms
            - 25ms
            - 50ms
            - 100ms
            - 250ms
            - 500ms
            - 1s
            - 2.5s
            - 5s
            - 10s
      dimensions:
        - name: http.method
        - name: http.route
        - name: http.status_code
      exemplars:
        enabled: true
      dimensions_cache_size: 1000
      aggregation_temporality: AGGREGATION_TEMPORALITY_CUMULATIVE
      metrics_flush_interval: 15s
    servicegraph:
      latency_histogram_buckets:
        - 5ms
        - 10ms
        - 25ms
        - 50ms
        - 100ms
        - 250ms
        - 500ms
        - 1s
        - 2.5s
        - 5s
      dimensions:
        - http.method
        - http.route
      store:
        ttl: 10s
        max_items: 1000
  exporters:
    otlphttp/tempo:
      endpoint: http://tempo:4318
    otlphttp/loki:
      endpoint: http://loki:3100/otlp
    prometheusremotewrite:
      endpoint: http://prometheus-server:80/api/v1/write
      resource_to_telemetry_conversion:
        enabled: true
  service:
    pipelines:
      traces:
        receivers: [otlp]
        exporters: [spanmetrics, servicegraph, otlphttp/tempo]
      metrics/spanmetrics:
        receivers: [spanmetrics, servicegraph]
        exporters: [prometheusremotewrite]
      metrics:
        receivers: [otlp]
        exporters: [prometheusremotewrite]
      logs:
        receivers: [otlp]
        exporters: [otlphttp/loki]
EOF

# --- Create dashboard ConfigMap from local JSON files ---
echo ""
echo "--- Creating dashboard ConfigMaps ---"
DASHBOARD_DIR="$PROJECT_ROOT/monitoring/grafana/dashboards"
if [ -d "$DASHBOARD_DIR" ]; then
  kubectl create configmap grafana-dashboards \
    --from-file="$DASHBOARD_DIR/" \
    --namespace "$NAMESPACE" --context "$CONTEXT" \
    --dry-run=client -o yaml | kubectl apply --context "$CONTEXT" -f -
  # Label the ConfigMap so Grafana sidecar can pick it up (not needed for our approach but good practice)
  kubectl label configmap grafana-dashboards grafana_dashboard=1 \
    --namespace "$NAMESPACE" --context "$CONTEXT" --overwrite
else
  echo "Warning: Dashboard directory not found at $DASHBOARD_DIR, skipping"
fi

# --- Grafana (visualization) ---
echo ""
echo "--- Installing Grafana ---"
helm upgrade --install grafana grafana/grafana \
  --version 10.5.15 \
  --namespace "$NAMESPACE" \
  --kube-context "$CONTEXT" \
  --values - <<'EOF'
adminUser: admin
adminPassword: workshop
extraInitContainers:
  - name: install-plugins
    image: grafana/grafana:12.3.1
    command:
      - sh
      - -c
      - |
        grafana cli --pluginsDir /var/lib/grafana/plugins plugins install grafana-lokiexplore-app &&
        grafana cli --pluginsDir /var/lib/grafana/plugins plugins install grafana-exploretraces-app &&
        grafana cli --pluginsDir /var/lib/grafana/plugins plugins install grafana-metricsdrilldown-app &&
        grafana cli --pluginsDir /var/lib/grafana/plugins plugins install grafana-pyroscope-app
    volumeMounts:
      - name: storage
        mountPath: /var/lib/grafana
grafana.ini:
  plugin.grafana-lokiexplore-app:
    enabled: true
  plugin.grafana-exploretraces-app:
    enabled: true
  plugin.grafana-metricsdrilldown-app:
    enabled: true
  plugin.grafana-pyroscope-app:
    enabled: true
datasources:
  datasources.yaml:
    apiVersion: 1
    datasources:
      - name: Loki
        type: loki
        uid: loki
        access: proxy
        url: http://loki:3100
        isDefault: false
      - name: Tempo
        type: tempo
        uid: tempo
        access: proxy
        url: http://tempo:3100
        isDefault: false
        jsonData:
          tracesToLogsV2:
            datasourceUid: loki
            spanStartTimeShift: "-1m"
            spanEndTimeShift: "1m"
            filterByTraceID: true
            filterBySpanID: false
          tracesToProfiles:
            datasourceUid: pyroscope
            customQuery: false
            profileTypeId: "process_cpu:cpu:nanoseconds:cpu:nanoseconds"
            tags:
              - key: service.name
                value: service_name
      - name: Pyroscope
        type: grafana-pyroscope-datasource
        uid: pyroscope
        access: proxy
        url: http://pyroscope:4040
        isDefault: false
      - name: Prometheus
        type: prometheus
        uid: prometheus
        access: proxy
        url: http://prometheus-server:80
        isDefault: true
        jsonData:
          exemplarTraceIdDestinations:
            - name: trace_id
              datasourceUid: tempo
dashboardProviders:
  dashboardproviders.yaml:
    apiVersion: 1
    providers:
      - name: default
        orgId: 1
        folder: ""
        type: file
        disableDeletion: false
        editable: true
        options:
          path: /var/lib/grafana/dashboards/default
dashboardsConfigMaps:
  default: grafana-dashboards
EOF

# --- Ingress for Grafana ---
echo ""
echo "--- Creating Grafana Ingress ---"
kubectl apply --context "$CONTEXT" -n "$NAMESPACE" -f - <<'EOF'
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: monitoring-ingress
spec:
  ingressClassName: monitoring
  rules:
    - http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: grafana
                port:
                  number: 80
EOF

echo ""
echo "=== LGTM stack deployment complete ==="
echo "Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod --all -n "$NAMESPACE" --context "$CONTEXT" --timeout=300s || true
kubectl get pods -n "$NAMESPACE" --context "$CONTEXT" -o wide

echo ""
echo "Grafana:        http://localhost:3000  (admin/workshop)"
echo "OTel Collector: localhost:4317 (gRPC), localhost:4318 (HTTP)"
echo "Pyroscope:      localhost:4040"
echo ""
echo "Cross-cluster endpoint (from app cluster):"
echo "  OTEL: http://k3d-sck-monitoring-serverlb:4317"
echo "  Pyroscope: http://k3d-sck-monitoring-serverlb:4040"
