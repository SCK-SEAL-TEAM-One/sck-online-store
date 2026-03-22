# Workshop Cluster Monitoring: node-exporter + kube-state-metrics + OTel Gateway
# Deployed to the workshop cluster (default provider)
# The OTel Gateway receives all local telemetry and forwards to the monitoring cluster's OTel Collector

# Look up the OTel Collector NLB hostname in the monitoring cluster
# The gateway forwards to this endpoint (gRPC)
data "kubernetes_service" "otel_collector" {
  provider = kubernetes.monitoring

  metadata {
    name      = "otel-collector-opentelemetry-collector"
    namespace = "monitoring"
  }

  depends_on = [helm_release.otel_collector]
}

locals {
  # Remote monitoring cluster endpoint — the gateway forwards all telemetry here
  otel_collector_grpc_endpoint = "${data.kubernetes_service.otel_collector.status[0].load_balancer[0].ingress[0].hostname}:4317"
}

resource "helm_release" "workshop_node_exporter" {
  name             = "prometheus-node-exporter"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "prometheus-node-exporter"
  version          = "4.37.3"
  namespace        = "monitoring"
  create_namespace = true

  values = [yamlencode({
    resources = {
      requests = {
        cpu    = "50m"
        memory = "32Mi"
      }
      limits = {
        memory = "64Mi"
      }
    }
  })]

  depends_on = [module.eks]
}

resource "helm_release" "workshop_kube_state_metrics" {
  name             = "kube-state-metrics"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-state-metrics"
  version          = "5.25.1"
  namespace        = "monitoring"
  create_namespace = true

  values = [yamlencode({
    replicas = 1
    resources = {
      requests = {
        cpu    = "50m"
        memory = "64Mi"
      }
      limits = {
        memory = "128Mi"
      }
    }
  })]

  depends_on = [module.eks]
}

resource "helm_release" "workshop_otel_gateway" {
  name             = "otel-gateway"
  repository       = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart            = "opentelemetry-collector"
  version          = "0.97.1"
  namespace        = "monitoring"
  create_namespace = true

  values = [yamlencode({
    mode         = "deployment"
    replicaCount = 1
    image = {
      repository = "otel/opentelemetry-collector-contrib"
    }
    ports = {
      otlp = {
        enabled       = true
        containerPort = 4317
        servicePort   = 4317
        protocol      = "TCP"
      }
      otlp-http = {
        enabled       = true
        containerPort = 4318
        servicePort   = 4318
        protocol      = "TCP"
      }
      jaeger-compact = {
        enabled = false
      }
      jaeger-thrift = {
        enabled = false
      }
      jaeger-grpc = {
        enabled = false
      }
      zipkin = {
        enabled = false
      }
    }
    config = {
      receivers = {
        otlp = {
          protocols = {
            grpc = {
              endpoint = "0.0.0.0:4317"
            }
            http = {
              endpoint = "0.0.0.0:4318"
            }
          }
        }
        prometheus = {
          config = {
            scrape_configs = [
              {
                job_name        = "node-exporter"
                scrape_interval = "30s"
                dns_sd_configs = [
                  {
                    names = ["prometheus-node-exporter.monitoring.svc.cluster.local"]
                    type  = "A"
                    port  = 9100
                  }
                ]
              },
              {
                job_name        = "kube-state-metrics"
                scrape_interval = "30s"
                static_configs = [
                  {
                    targets = ["kube-state-metrics.monitoring.svc.cluster.local:8080"]
                  }
                ]
              }
            ]
          }
        }
      }
      processors = {
        batch = {
          send_batch_size = 1024
          timeout         = "5s"
        }
        memory_limiter = {
          check_interval  = "5s"
          limit_mib       = 256
          spike_limit_mib = 64
        }
      }
      exporters = {
        otlp = {
          endpoint = local.otel_collector_grpc_endpoint
          tls = {
            insecure = true
          }
          sending_queue = {
            enabled       = true
            num_consumers = 4
            queue_size    = 256
          }
          retry_on_failure = {
            enabled          = true
            initial_interval = "5s"
            max_interval     = "30s"
            max_elapsed_time = "300s"
          }
        }
      }
      service = {
        pipelines = {
          traces = {
            receivers  = ["otlp"]
            processors = ["memory_limiter", "batch"]
            exporters  = ["otlp"]
          }
          metrics = {
            receivers  = ["otlp", "prometheus"]
            processors = ["memory_limiter", "batch"]
            exporters  = ["otlp"]
          }
          logs = {
            receivers  = ["otlp"]
            processors = ["memory_limiter", "batch"]
            exporters  = ["otlp"]
          }
        }
      }
    }
  })]

  depends_on = [
    module.eks,
    helm_release.workshop_node_exporter,
    helm_release.workshop_kube_state_metrics,
    helm_release.otel_collector
  ]
}
