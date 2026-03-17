# LGTM Stack: Loki + Grafana + Tempo + Mimir + OpenTelemetry Collector
# Deployed to the monitoring cluster

resource "helm_release" "loki" {
  provider = helm.monitoring

  name             = "loki"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "loki"
  version          = "6.16.0"
  namespace        = "monitoring"
  create_namespace = true

  values = [yamlencode({
    deploymentMode = "SingleBinary"
    loki = {
      auth_enabled = false
      commonConfig = {
        replication_factor = 1
      }
      storage = {
        type = "filesystem"
      }
      schemaConfig = {
        configs = [
          {
            from         = "2024-01-01"
            store        = "tsdb"
            object_store = "filesystem"
            schema       = "v13"
            index = {
              prefix = "loki_index_"
              period = "24h"
            }
          }
        ]
      }
      limits_config = {
        allow_structured_metadata = true
      }
    }
    singleBinary = {
      replicas = 1
      persistence = {
        storageClass = "gp2"
        size         = "10Gi"
      }
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          memory = "512Mi"
        }
      }
    }
    gateway = {
      enabled = false
    }
    chunksCache = {
      enabled = false
    }
    resultsCache = {
      enabled = false
    }
    backend = {
      replicas = 0
    }
    read = {
      replicas = 0
    }
    write = {
      replicas = 0
    }
  })]

  depends_on = [module.monitoring_eks]
}

resource "helm_release" "tempo" {
  provider = helm.monitoring

  name             = "tempo"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "tempo"
  version          = "1.10.3"
  namespace        = "monitoring"
  create_namespace = true

  values = [yamlencode({
    tempo = {
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          memory = "512Mi"
        }
      }
    }
  })]

  depends_on = [module.monitoring_eks]
}

resource "helm_release" "prometheus" {
  provider = helm.monitoring

  name             = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "prometheus"
  version          = "25.27.0"
  namespace        = "monitoring"
  create_namespace = true

  values = [yamlencode({
    server = {
      extraFlags = [
        "web.enable-remote-write-receiver"
      ]
      persistentVolume = {
        storageClass = "gp2"
        size         = "10Gi"
      }
      resources = {
        requests = {
          cpu    = "100m"
          memory = "256Mi"
        }
        limits = {
          memory = "512Mi"
        }
      }
    }
    alertmanager = {
      enabled = false
    }
    kube-state-metrics = {
      enabled = false
    }
    prometheus-node-exporter = {
      enabled = false
    }
    prometheus-pushgateway = {
      enabled = false
    }
  })]

  depends_on = [module.monitoring_eks]
}

resource "helm_release" "grafana" {
  provider = helm.monitoring

  name             = "grafana"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "grafana"
  version          = "8.5.12"
  namespace        = "monitoring"
  create_namespace = true

  values = [yamlencode({
    adminUser     = "admin"
    adminPassword = "workshop"
    service = {
      type = "LoadBalancer"
      annotations = {
        "service.beta.kubernetes.io/aws-load-balancer-type"   = "nlb"
        "service.beta.kubernetes.io/aws-load-balancer-scheme" = "internet-facing"
      }
    }
    datasources = {
      "datasources.yaml" = {
        apiVersion = 1
        datasources = [
          {
            name      = "Loki"
            type      = "loki"
            access    = "proxy"
            url       = "http://loki:3100"
            isDefault = false
          },
          {
            name   = "Tempo"
            type   = "tempo"
            access = "proxy"
            url    = "http://tempo:3100"
            jsonData = {
              tracesToLogsV2 = {
                datasourceUid = "loki"
              }
            }
            isDefault = false
          },
          {
            name      = "Prometheus"
            type      = "prometheus"
            access    = "proxy"
            url       = "http://prometheus-server:80"
            isDefault = true
          }
        ]
      }
    }
    dashboardProviders = {
      "dashboardproviders.yaml" = {
        apiVersion = 1
        providers = [
          {
            name            = "default"
            orgId           = 1
            folder          = ""
            type            = "file"
            disableDeletion = false
            editable        = true
            options = {
              path = "/var/lib/grafana/dashboards/default"
            }
          }
        ]
      }
    }
    dashboards = {
      default = {
        store-business-metrics = {
          json = file("${path.module}/../../monitoring/grafana/dashboards/store-business-metrics.json")
        }
        store-workflow-monitor = {
          json = file("${path.module}/../../monitoring/grafana/dashboards/store-workflow-monitor.json")
        }
      }
    }
  })]

  depends_on = [module.monitoring_eks, helm_release.loki, helm_release.tempo, helm_release.prometheus]
}

resource "helm_release" "otel_collector" {
  provider = helm.monitoring

  name             = "otel-collector"
  repository       = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart            = "opentelemetry-collector"
  version          = "0.97.1"
  namespace        = "monitoring"
  create_namespace = true

  values = [yamlencode({
    mode = "deployment"
    replicaCount = 1
    image = {
      repository = "otel/opentelemetry-collector-contrib"
    }
    service = {
      type = "LoadBalancer"
      annotations = {
        "service.beta.kubernetes.io/aws-load-balancer-type"   = "nlb"
        "service.beta.kubernetes.io/aws-load-balancer-scheme" = "internal"
      }
    }
    ports = {
      otlp = {
        enabled     = true
        containerPort = 4317
        servicePort = 4317
        protocol    = "TCP"
      }
      otlp-http = {
        enabled     = true
        containerPort = 4318
        servicePort = 4318
        protocol    = "TCP"
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
      }
      exporters = {
        otlphttp_tempo = {
          endpoint = "http://tempo:4318"
        }
        otlphttp_loki = {
          endpoint = "http://loki:3100/otlp"
        }
        prometheusremotewrite = {
          endpoint = "http://prometheus-server:80/api/v1/write"
        }
      }
      service = {
        pipelines = {
          traces = {
            receivers = ["otlp"]
            exporters = ["otlphttp_tempo"]
          }
          logs = {
            receivers = ["otlp"]
            exporters = ["otlphttp_loki"]
          }
          metrics = {
            receivers = ["otlp"]
            exporters = ["prometheusremotewrite"]
          }
        }
      }
    }
  })]

  depends_on = [module.monitoring_eks, helm_release.loki, helm_release.tempo, helm_release.prometheus]
}
