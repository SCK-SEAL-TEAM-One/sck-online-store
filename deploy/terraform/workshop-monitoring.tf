# Workshop Cluster Monitoring: node-exporter + kube-state-metrics + OTel Collector scraper
# Deployed to the workshop cluster (default provider) — pushes metrics to monitoring cluster

# Look up the OTel Collector NLB hostname in the monitoring cluster
data "kubernetes_service" "otel_collector" {
  provider = kubernetes.monitoring

  metadata {
    name      = "otel-collector-opentelemetry-collector"
    namespace = "monitoring"
  }

  depends_on = [helm_release.otel_collector]
}

locals {
  otel_collector_endpoint = "http://${data.kubernetes_service.otel_collector.status[0].load_balancer[0].ingress[0].hostname}:4318"
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

resource "helm_release" "workshop_otel_scraper" {
  name             = "otel-scraper"
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
    service = {
      enabled = false
    }
    ports = {
      otlp = {
        enabled = false
      }
      otlp-http = {
        enabled = false
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
      exporters = {
        otlphttp = {
          endpoint = local.otel_collector_endpoint
        }
      }
      service = {
        pipelines = {
          metrics = {
            receivers = ["prometheus"]
            exporters = ["otlphttp"]
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
