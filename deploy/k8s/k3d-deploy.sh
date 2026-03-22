#!/bin/bash
# Deploy app to local k3d cluster
# Usage: ./deploy/k8s/k3d-deploy.sh
#
# Environment variables:
#   K8S_CONTEXT          kubectl context (default: k3d-sck-workshop)
#   K8S_NAMESPACE        target namespace (default: public)
#   OTEL_ENDPOINT        Remote OTel collector (monitoring cluster), e.g. http://k3d-sck-monitoring-serverlb:4317
#                        When set, deploys a local otel-gateway that receives all telemetry
#                        from services and forwards to this remote endpoint
#   PYROSCOPE_ENDPOINT   Pyroscope endpoint, e.g. http://k3d-sck-monitoring-serverlb:4040
#
# Prerequisites:
#   k3d cluster running with images imported (sck/*:k3d)
#   namespace 'public' created with secrets (db-credentials, app-secrets)

set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
NAMESPACE="${K8S_NAMESPACE:-public}"
CONTEXT="${K8S_CONTEXT:-k3d-sck-workshop}"

# Optional: set these to connect to a monitoring cluster
OTEL_ENDPOINT="${OTEL_ENDPOINT:-}"
PYROSCOPE_ENDPOINT="${PYROSCOPE_ENDPOINT:-}"

# When monitoring is enabled, services send to local otel-gateway (not directly to remote collector)
GATEWAY_SVC="otel-gateway-opentelemetry-collector.$NAMESPACE"

echo "=== Deploying to k3d (namespace: $NAMESPACE, context: $CONTEXT) ==="

# Apply monitoring-endpoints ConfigMap
# When monitoring is enabled, point services to the local otel-gateway (not the remote collector)
if [ -n "$OTEL_ENDPOINT" ] || [ -n "$PYROSCOPE_ENDPOINT" ]; then
  echo "Patching monitoring-endpoints to use local otel-gateway..."
  OTEL_VALUE="http://$GATEWAY_SVC:4317"
  PYRO_VALUE="${PYROSCOPE_ENDPOINT:-http://REPLACE_PYROSCOPE_ENDPOINT:4040}"
  sed "s|http://REPLACE_OTEL_ENDPOINT:4317|$OTEL_VALUE|;s|http://REPLACE_PYROSCOPE_ENDPOINT:4040|$PYRO_VALUE|" \
    "$SCRIPT_DIR/app/monitoring-endpoints.yml" | kubectl apply --context "$CONTEXT" -n "$NAMESPACE" -f -
else
  kubectl apply --context "$CONTEXT" -f "$SCRIPT_DIR/app/monitoring-endpoints.yml" -n "$NAMESPACE"
fi

# Deploy store-database with local image
# When OTEL_ENDPOINT is set, use the OTel Collector sidecar for MySQL metrics + slow query logs
if [ -n "$OTEL_ENDPOINT" ]; then
  echo "Deploying store-database with OTel Collector sidecar..."
  sed "s|image: siamchamnankit/store-database:.*|image: sck/store-database:k3d|" \
    "$SCRIPT_DIR/monitoring/store-database-with-otel.yml" \
    | sed "s|http://REPLACE_OTEL_ENDPOINT:4318|http://$GATEWAY_SVC:4318|" \
    | kubectl apply --context "$CONTEXT" -n "$NAMESPACE" -f -
else
  sed 's|image: siamchamnankit/store-database:.*|image: sck/store-database:k3d|;s|imagePullPolicy: Always|imagePullPolicy: Never|' \
    "$SCRIPT_DIR/app/store-database/service.yml" | kubectl apply --context "$CONTEXT" -n "$NAMESPACE" -f -
fi

echo "Waiting for MySQL to be ready..."
kubectl wait --for=condition=ready pod -l app=mysql -n "$NAMESPACE" --context "$CONTEXT" --timeout=180s || true
sleep 5  # extra wait for MySQL to accept connections

# Run seed job (creates databases + tables + data)
echo "Running seed job..."
kubectl delete job mysql-seed -n "$NAMESPACE" --context "$CONTEXT" 2>/dev/null || true
kubectl apply --context "$CONTEXT" -f "$SCRIPT_DIR/app/store-database/seed-job.yml" -n "$NAMESPACE"
echo "Waiting for seed to complete..."
kubectl wait --for=condition=complete job/mysql-seed -n "$NAMESPACE" --context "$CONTEXT" --timeout=180s || {
  echo "Seed job logs:"
  kubectl logs job/mysql-seed -n "$NAMESPACE" --context "$CONTEXT" --tail=30 || true
}

# Deploy liquibase with local image
echo "Running Liquibase migration..."
kubectl delete job liquibase-migration-latest -n "$NAMESPACE" --context "$CONTEXT" 2>/dev/null || true
sed 's|image: siamchamnankit/liquibase:.*|image: sck/liquibase:k3d|;s|imagePullPolicy: Always|imagePullPolicy: Never|' \
  "$SCRIPT_DIR/app/liquibase/job.yml" | kubectl apply --context "$CONTEXT" -n "$NAMESPACE" -f -

# Wait for migration
echo "Waiting for Liquibase migration..."
kubectl wait --for=condition=complete job/liquibase-migration-latest -n "$NAMESPACE" --context "$CONTEXT" --timeout=120s || {
  echo "Liquibase logs:"
  kubectl logs job/liquibase-migration-latest -n "$NAMESPACE" --context "$CONTEXT" --tail=30 || true
}

# Deploy thirdparty with local image
# When OTEL_ENDPOINT is set, use the Beyla eBPF sidecar for auto-instrumentation
if [ -n "$OTEL_ENDPOINT" ]; then
  echo "Deploying thirdparty with Beyla sidecar..."
  sed "s|image: siamchamnankit/store-shipping-gateway:.*|image: sck/store-shipping-gateway:k3d|" \
    "$SCRIPT_DIR/monitoring/thirdparty-with-beyla.yml" \
    | sed "s|http://REPLACE_OTEL_ENDPOINT:4317|http://$GATEWAY_SVC:4317|" \
    | kubectl apply --context "$CONTEXT" -n "$NAMESPACE" -f -
else
  sed 's|image: siamchamnankit/store-shipping-gateway:.*|image: sck/store-shipping-gateway:k3d|;s|imagePullPolicy: Always|imagePullPolicy: Never|' \
    "$SCRIPT_DIR/app/thirdparty/service.yml" | kubectl apply --context "$CONTEXT" -n "$NAMESPACE" -f -
fi

# Deploy store-service with local image + point to localhost store-web
sed 's|image: siamchamnankit/store-service:.*|image: sck/store-service:k3d|;s|imagePullPolicy: Always|imagePullPolicy: Never|' \
  "$SCRIPT_DIR/app/store-service/service.yml" | kubectl apply --context "$CONTEXT" -n "$NAMESPACE" -f -

# Deploy point-service with local image
sed 's|image: siamchamnankit/point-service:.*|image: sck/point-service:k3d|;s|imagePullPolicy: Always|imagePullPolicy: Never|' \
  "$SCRIPT_DIR/app/point-service/service.yml" | kubectl apply --context "$CONTEXT" -n "$NAMESPACE" -f -

# Deploy store-web with local image
sed 's|image: siamchamnankit/store-web:.*|image: sck/store-web:k3d|;s|imagePullPolicy: Always|imagePullPolicy: Never|' \
  "$SCRIPT_DIR/app/store-web/service.yml" | kubectl apply --context "$CONTEXT" -n "$NAMESPACE" -f -

# Deploy ingress
kubectl apply --context "$CONTEXT" -f "$SCRIPT_DIR/app/ingress.yml" -n "$NAMESPACE"

# --- Deploy monitoring agents on app cluster (node-exporter, kube-state-metrics) ---
# These run on the app cluster nodes and forward metrics to monitoring cluster via OTLP
if [ -n "$OTEL_ENDPOINT" ]; then
  echo ""
  echo "--- Deploying monitoring agents (node-exporter, kube-state-metrics) ---"

  helm repo add prometheus-community https://prometheus-community.github.io/helm-charts 2>/dev/null || true
  helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts 2>/dev/null || true
  helm repo update

  # Node exporter — DaemonSet with hostNetwork, exposes host metrics on :9100
  echo "Installing node-exporter..."
  helm upgrade --install node-exporter prometheus-community/prometheus-node-exporter \
    --namespace "$NAMESPACE" --kube-context "$CONTEXT" \
    --set fullnameOverride=node-exporter

  # Kube state metrics — k8s object metrics (pods, deployments, etc.)
  echo "Installing kube-state-metrics..."
  helm upgrade --install kube-state-metrics prometheus-community/kube-state-metrics \
    --namespace "$NAMESPACE" --kube-context "$CONTEXT" \
    --set fullnameOverride=kube-state-metrics

  # OTel Gateway — receives all local telemetry + scrapes prometheus, forwards to monitoring cluster
  OTEL_GRPC_ENDPOINT=$(echo "$OTEL_ENDPOINT" | sed 's|^http://||;s|^https://||')
  echo "Installing otel-gateway (local telemetry gateway, forwards to $OTEL_GRPC_ENDPOINT)..."
  helm upgrade --install otel-gateway open-telemetry/opentelemetry-collector \
    --version 0.97.1 \
    --namespace "$NAMESPACE" --kube-context "$CONTEXT" \
    --values - <<EOF
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
    prometheus:
      config:
        scrape_configs:
          - job_name: node-exporter
            scrape_interval: 15s
            static_configs:
              - targets: ['node-exporter:9100']
                labels:
                  cluster: sck-workshop
          - job_name: kube-state-metrics
            scrape_interval: 15s
            static_configs:
              - targets: ['kube-state-metrics:8080']
                labels:
                  cluster: sck-workshop
  processors:
    batch:
      send_batch_size: 1024
      timeout: 5s
    memory_limiter:
      check_interval: 5s
      limit_mib: 256
      spike_limit_mib: 64
  exporters:
    otlp:
      endpoint: ${OTEL_GRPC_ENDPOINT}
      tls:
        insecure: true
      sending_queue:
        enabled: true
        num_consumers: 4
        queue_size: 256
      retry_on_failure:
        enabled: true
        initial_interval: 5s
        max_interval: 30s
        max_elapsed_time: 300s
  service:
    pipelines:
      traces:
        receivers: [otlp]
        processors: [memory_limiter, batch]
        exporters: [otlp]
      metrics:
        receivers: [otlp, prometheus]
        processors: [memory_limiter, batch]
        exporters: [otlp]
      logs:
        receivers: [otlp]
        processors: [memory_limiter, batch]
        exporters: [otlp]
EOF
fi

echo ""
echo "=== Deployment complete ==="
echo "Waiting for all pods to be ready..."
kubectl wait --for=condition=ready pod --all -n "$NAMESPACE" --context "$CONTEXT" --timeout=180s || true
kubectl get pods -n "$NAMESPACE" --context "$CONTEXT" -o wide
echo ""
echo "App should be accessible at http://localhost"
