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

  values = [<<-EOT
tempo:
  resources:
    requests:
      cpu: "100m"
      memory: "512Mi"
    limits:
      memory: "2Gi"
config: |
  multitenancy_enabled: false
  usage_report:
    reporting_enabled: true
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
EOT
  ]

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
        "web.enable-remote-write-receiver",
        "enable-feature=exemplar-storage"
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

resource "helm_release" "pyroscope" {
  provider = helm.monitoring

  name             = "pyroscope"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "pyroscope"
  version          = "1.7.1"
  namespace        = "monitoring"
  create_namespace = true

  values = [yamlencode({
    pyroscope = {
      extraArgs = {
        "store-gateway.sharding-ring.replication-factor" = "1"
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
      persistence = {
        enabled          = true
        storageClassName = "gp2"
        size             = "10Gi"
      }
    }
    alloy = {
      enabled = false
    }
    minio = {
      enabled = false
    }
    service = {
      type = "LoadBalancer"
      annotations = {
        "service.beta.kubernetes.io/aws-load-balancer-type"   = "nlb"
        "service.beta.kubernetes.io/aws-load-balancer-scheme" = "internal"
      }
    }
  })]

  depends_on = [module.monitoring_eks]
}

resource "helm_release" "grafana" {
  provider = helm.monitoring

  name             = "grafana"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "grafana"
  version          = "10.5.15"
  namespace        = "monitoring"
  create_namespace = true

  values = [yamlencode({
    adminUser     = "admin"
    adminPassword = "workshop"
    extraInitContainers = [
      {
        name  = "install-plugins"
        image = "grafana/grafana:12.3.1"
        command = ["sh", "-c", <<-EOT
          grafana cli --pluginsDir /var/lib/grafana/plugins plugins install grafana-lokiexplore-app &&
          grafana cli --pluginsDir /var/lib/grafana/plugins plugins install grafana-exploretraces-app &&
          grafana cli --pluginsDir /var/lib/grafana/plugins plugins install grafana-metricsdrilldown-app &&
          grafana cli --pluginsDir /var/lib/grafana/plugins plugins install grafana-pyroscope-app
        EOT
        ]
        volumeMounts = [
          {
            name      = "storage"
            mountPath = "/var/lib/grafana"
          }
        ]
      }
    ]
    "grafana.ini" = {
      "plugin.grafana-lokiexplore-app" = {
        enabled = true
      }
      "plugin.grafana-exploretraces-app" = {
        enabled = true
      }
      "plugin.grafana-metricsdrilldown-app" = {
        enabled = true
      }
      "plugin.grafana-pyroscope-app" = {
        enabled = true
      }
    }
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
            uid       = "loki"
            access    = "proxy"
            url       = "http://loki:3100"
            isDefault = false
          },
          {
            name   = "Tempo"
            type   = "tempo"
            uid    = "tempo"
            access = "proxy"
            url    = "http://tempo:3100"
            jsonData = {
              tracesToLogsV2 = {
                datasourceUid          = "loki"
                spanStartTimeShift     = "-1m"
                spanEndTimeShift       = "1m"
              }
              tracesToProfiles = {
                datasourceUid = "pyroscope"
                customQuery   = false
                profileTypeId = "process_cpu:cpu:nanoseconds:cpu:nanoseconds"
                tags = [
                  {
                    key   = "service.name"
                    value = "service_name"
                  }
                ]
              }
            }
            isDefault = false
          },
          {
            name      = "Pyroscope"
            type      = "grafana-pyroscope-datasource"
            uid       = "pyroscope"
            access    = "proxy"
            url       = "http://pyroscope:4040"
            isDefault = false
          },
          {
            name      = "Prometheus"
            type      = "prometheus"
            uid       = "prometheus"
            access    = "proxy"
            url       = "http://prometheus-server:80"
            isDefault = true
            jsonData = {
              exemplarTraceIdDestinations = [
                {
                  name          = "trace_id"
                  datasourceUid = "tempo"
                }
              ]
            }
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
        mysql-overview = {
          json = file("${path.module}/../../monitoring/grafana/dashboards/mysql-overview.json")
        }
        node-exporter = {
          json = file("${path.module}/../../monitoring/grafana/dashboards/node-exporter.json")
        }
        node-exporter-full = {
          gnetId     = 1860
          revision   = 37
          datasource = "Prometheus"
        }
        node-exporter-for-prometheus = {
          gnetId     = 11074
          revision   = 9
          datasource = "Prometheus"
        }
        k8s-cluster-resources = {
          gnetId     = 7249
          revision   = 1
          datasource = "Prometheus"
        }
        k8s-node-pods-resources = {
          gnetId     = 15760
          revision   = 2
          datasource = "Prometheus"
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
    mode         = "deployment"
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
      }
      connectors = {
        spanmetrics = {
          histogram = {
            explicit = {
              buckets = ["5ms", "10ms", "25ms", "50ms", "100ms", "250ms", "500ms", "1s", "2.5s", "5s", "10s"]
            }
          }
          dimensions = [
            { name = "http.method" },
            { name = "http.route" },
            { name = "http.status_code" }
          ]
          exemplars = {
            enabled = true
          }
          dimensions_cache_size   = 1000
          aggregation_temporality = "AGGREGATION_TEMPORALITY_CUMULATIVE"
          metrics_flush_interval  = "15s"
        }
        servicegraph = {
          latency_histogram_buckets = ["5ms", "10ms", "25ms", "50ms", "100ms", "250ms", "500ms", "1s", "2.5s", "5s"]
          dimensions                = ["http.method", "http.route"]
          store = {
            ttl       = "10s"
            max_items = 1000
          }
        }
      }
      exporters = {
        "otlphttp/tempo" = {
          endpoint = "http://tempo:4318"
        }
        "otlphttp/loki" = {
          endpoint = "http://loki:3100/otlp"
        }
        prometheusremotewrite = {
          endpoint = "http://prometheus-server:80/api/v1/write"
          resource_to_telemetry_conversion = {
            enabled = true
          }
        }
      }
      service = {
        pipelines = {
          traces = {
            receivers = ["otlp"]
            exporters = ["spanmetrics", "servicegraph", "otlphttp/tempo"]
          }
          "metrics/spanmetrics" = {
            receivers = ["spanmetrics", "servicegraph"]
            exporters = ["prometheusremotewrite"]
          }
          metrics = {
            receivers = ["otlp"]
            exporters = ["prometheusremotewrite"]
          }
          logs = {
            receivers = ["otlp"]
            exporters = ["otlphttp/loki"]
          }
        }
      }
    }
  })]

  depends_on = [module.monitoring_eks, helm_release.loki, helm_release.tempo, helm_release.prometheus]
}
