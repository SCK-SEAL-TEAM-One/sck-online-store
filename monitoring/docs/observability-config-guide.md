# Observability Configuration Guide

> **Last updated:** 2026-03-22
> This document explains every observability configuration in the project line by line, what each setting controls, and what breaks when it changes or is removed.

## Overview: How Signals Flow

### Docker Compose (local dev)

In Docker Compose, all services send directly to `lgtm` (grafana/otel-lgtm all-in-one container):

```
  store-service (Go)         point-service (Node.js)       thirdparty (Node.js)
  ┌──────────────────┐       ┌──────────────────┐          ┌──────────────────┐
  │ OTel SDK (Go)    │       │ OTel SDK (Node)  │          │ (no OTel SDK)    │
  │  traces ─────────┤       │  traces ─────────┤          │                  │
  │  metrics ────────┤       │  logs ───────────┤          └───────┬──────────┘
  │  logs ───────────┤       │                  │                  │
  │                  │       │ Pyroscope SDK    │           Beyla (eBPF)
  │ Pyroscope SDK    │       │  wall profiles ──┤          │  traces ─────────┐
  │  cpu profiles ───┤       └────────┬─────────┘          └────────┬─────────┘
  │                  │                │                              │
  │ otel-profiling-go│                │                              │
  │ (adds profile.id)│                │                              │
  └────────┬─────────┘                │                              │
           │                          │                              │
           ▼                          ▼                              ▼
    ┌──────────────────────────────────────────────────────────────────┐
    │                     lgtm (grafana/otel-lgtm)                    │
    │  :4317 gRPC ──→ OTel Collector ──→ Tempo (traces)              │
    │  :4318 HTTP                    ──→ Prometheus (metrics)         │
    │                                ──→ Loki (logs)                  │
    │  :4040 ──→ Pyroscope (profiles, direct push, NOT via OTel)     │
    │  :3000 ──→ Grafana (dashboards)                                │
    └──────────────────────────────────────────────────────────────────┘
           ▲
           │
    otel-collector (contrib)
    │  MySQL metrics ──→ Prometheus (via lgtm:4317)
    │  MySQL slow logs ──→ Loki (via lgtm:4317)
```

### k3d / EKS (Agent-Gateway Pattern)

In the two-cluster setup, all telemetry flows through a local **OTel Gateway** on the app cluster before being forwarded to the monitoring cluster. This provides buffering, retry, and a single cross-cluster connection.

```
APP CLUSTER (sck-workshop)                    MONITORING CLUSTER (sck-monitoring)
┌───────────────────────────────────┐         ┌─────────────────────────────────────┐
│ store-service ──┐                 │         │                                     │
│ point-service ──┤                 │  OTLP   │ OTel Collector                      │
│ Beyla sidecar ──┼► OTel Gateway ─┼──gRPC──►│ ► spanmetrics/servicegraph          │
│ MySQL sidecar ──┤  (forwarder)   │         │ ► export to Tempo, Loki, Prometheus │
│ node-exporter ──┤  batch/retry   │         │                                     │
│ kube-state-m  ──┘                │         │ Pyroscope ◄── direct from services  │
└───────────────────────────────────┘         │ Grafana                             │
                                              └─────────────────────────────────────┘
```

**Key design:** No processing duplication. The app cluster gateway is a lightweight forwarder (OTLP recv + batch + retry). All connectors (spanmetrics, servicegraph) and backend routing stay in the monitoring cluster's OTel Collector.

**Pyroscope** is accessed directly (not through the gateway) because it uses its own push API, not OTLP.

### Signal Correlation Across the Gateway

All signal correlations are preserved through the gateway — it forwards full OTLP payloads without stripping attributes:

| Correlation | How it works | Key attribute |
|---|---|---|
| Trace → Log | Logs carry `trace_id` + `span_id` via OTLP. Grafana Tempo `tracesToLogsV2` links to Loki | `trace_id`, `span_id` in log labels |
| Trace → Profile | `otel-profiling-go` injects `pyroscope.profile.id` into spans. Grafana `tracesToProfiles` links to Pyroscope | `pyroscope.profile.id` span attribute |
| Metric → Trace | `spanmetrics` connector generates exemplars with `trace_id`. Grafana Prometheus `exemplarTraceIdDestinations` links to Tempo | `trace_id` in exemplar |
| Log → Trace | Loki logs have `trace_id` label, enabling reverse lookup | `trace_id` in log labels |

---

## 1. docker-compose.yml — Service Environment Variables

### store-service

```yaml
environment:
  - OTEL_EXPORTER_OTLP_ENDPOINT=http://lgtm:4317    # [1]
  - OTEL_EXPORTER_OTLP_PROTOCOL=grpc                 # [2]
  - OTEL_SERVICE_NAME=store-service                   # [3]
  - PYROSCOPE_URL=http://lgtm:4040                    # [4]
```

| # | What it does | If removed/changed |
|---|---|---|
| [1] | Tells OTel SDK where to send traces, metrics, logs via gRPC | No telemetry exported. Traces, metrics, logs all stop. OTel SDK init skipped entirely (checked in `cmd/main.go:55`) |
| [2] | Selects gRPC protocol (vs `http/protobuf`). Read by `internal/otel/otel.go:33` | Defaults to `grpc` if empty. Change to `http/protobuf` for HTTP transport (uses port 4318 instead) |
| [3] | Sets `service.name` resource attribute on all OTel signals | All telemetry shows as "unknown_service". Grafana service maps, filters, and dashboard queries break |
| [4] | Tells Pyroscope SDK where to push profiles | No profiling. `cmd/main.go:64` skips Pyroscope init. `otel-profiling-go` wrapper is also skipped, so `pyroscope.profile.id` disappears from spans and "Profiles for this span" buttons vanish |

### point-service

```yaml
environment:
  - DB_HOST=db                                        # [5]
  - DB_PORT=3306                                      # [6]
  - DB_USERNAME=user                                  # [7]
  - DB_PASSWORD=password                              # [8]
  - OTEL_EXPORTER_OTLP_ENDPOINT=http://lgtm:4317     # [9]
  - OTEL_EXPORTER_OTLP_PROTOCOL=grpc                 # [10]
  - NODE_DEBUG=otel                                   # [11]
  - PYROSCOPE_URL=http://lgtm:4040                    # [12]
```

| # | What it does | If removed/changed |
|---|---|---|
| [5-8] | TypeORM MySQL connection. Used by `app.module.ts` via `ConfigModule` | Service crashes on startup — cannot connect to database |
| [9] | OTel SDK sends traces + logs to lgtm via gRPC | No traces or logs exported. SDK still starts but exports nowhere |
| [10] | Selects gRPC exporter in `trace.ts` (hardcoded to grpc imports regardless) | No effect currently — `trace.ts` always uses grpc importers. But keep for consistency |
| [11] | Node.js debug output for OTel internals | Removed = no debug logs. No impact on functionality |
| [12] | Pyroscope SDK sends wall profiles to lgtm:4040 | `trace.ts` skips Pyroscope init. No profiling. Flame graphs in Grafana Drilldown > Profiles disappear for point-service |

### beyla-thirdparty

```yaml
environment:
  - BEYLA_OPEN_PORT=8882,8883                          # [13]
  - OTEL_EXPORTER_OTLP_ENDPOINT=http://lgtm:4317      # [14]
  - OTEL_EXPORTER_OTLP_PROTOCOL=grpc                   # [15]
  - OTEL_SERVICE_NAME=thirdparty-gateway                # [16]
  - BEYLA_BPF_TRACK_REQUEST_HEADERS=true                # [17]
```

| # | What it does | If removed/changed |
|---|---|---|
| [13] | Beyla eBPF monitors these ports for HTTP traffic | No spans generated for thirdparty service. It disappears from service maps and traces |
| [14-15] | Where Beyla sends auto-generated traces | No telemetry from thirdparty |
| [16] | Service name in traces | Shows as "unknown_service" in Grafana |
| [17] | Captures HTTP request headers in spans (e.g., traceparent) | Trace context propagation may break. Beyla won't join distributed traces |

Also note:
```yaml
privileged: true       # Required for eBPF. Without it, Beyla cannot attach to kernel
pid: "service:thirdparty"  # Shares PID namespace with thirdparty container — required to see its processes
```

### otel-collector

```yaml
volumes:
  - ./monitoring/otel-collector/config.yaml:/etc/otelcol-contrib/config.yaml:ro  # [18]
  - mysql-logs:/var/log/mysql:ro                                                  # [19]
```

| # | What it does | If removed/changed |
|---|---|---|
| [18] | Collector config (receivers, processors, exporters). See Section 3 below | Collector starts with no pipelines. No MySQL metrics or slow query logs |
| [19] | Mounts MySQL log volume so filelog receiver can read slow query logs | Slow query log pipeline breaks. MySQL metrics still work |

### lgtm container

```yaml
ports:
  - "4040:4040"    # Pyroscope HTTP API (profile push + query)
  - "4317:4317"    # OTLP gRPC receiver (traces, metrics, logs)
  - "4318:4318"    # OTLP HTTP receiver
  - "3001:3000"    # Grafana UI (mapped to 3001 because store-web uses 3000)
volumes:
  - ./monitoring/grafana/provisioning/dashboards.yml:/.../custom-dashboards.yaml   # [20]
  - ./monitoring/grafana/provisioning/datasources.yml:/.../grafana-datasources.yaml # [21]
  - ./monitoring/grafana/dashboards:/var/lib/grafana/dashboards                     # [22]
```

| # | What it does | If removed/changed |
|---|---|---|
| [20] | Tells Grafana where to find dashboard JSON files | Custom dashboards disappear from Grafana |
| [21] | **Critical** — Provisions all datasource connections (Prometheus, Tempo, Loki, Pyroscope) with cross-linking config. See Section 2 below | All datasource connections revert to otel-lgtm defaults (no trace-to-log, trace-to-profile linking). Must re-configure manually in Grafana UI |
| [22] | Mounts dashboard JSON files | Dashboards not available |

---

## 2. datasources.yml — Grafana Datasource Provisioning (Line by Line)

File: `monitoring/grafana/provisioning/datasources.yml`
Mount target: `/otel-lgtm/grafana/conf/provisioning/datasources/grafana-datasources.yaml`

This file **overrides** the default datasource config from the `grafana/otel-lgtm` image. It defines 4 datasources and their cross-linking configuration.

### Default vs Custom: What We Changed from the lgtm Image

The `grafana/otel-lgtm` image ships with its own `grafana-datasources.yaml`. Our `datasources.yml` volume mount replaces it entirely. Below is a comparison of every section showing what the image provides by default and what this project customizes.

#### Prometheus — No changes

```
Default (lgtm image)                    │ Custom (this project)
────────────────────────────────────────│────────────────────────────────────
(identical)                             │ (identical)
```

The Prometheus datasource config is identical. No customization needed. The exemplar-to-trace linking works out of the box.

#### Tempo — 2 additions

```
Default (lgtm image)                    │ Custom (this project)
────────────────────────────────────────│────────────────────────────────────
tracesToLogsV2:                         │ tracesToLogsV2:
  customQuery: true                     │   customQuery: true
  datasourceUid: loki                   │   datasourceUid: loki
  query: '{$$...} | trace_id = "$$..."' │   query: '{$$...} | trace_id = "$$..."'
  tags:                                 │   spanStartTimeShift: "-1m"        ← ADDED
    - key: service.name                 │   spanEndTimeShift: "1m"           ← ADDED
      value: service_name               │   tags:
                                        │     - key: service.name
                                        │       value: service_name
                                        │
(not present)                           │ tracesToProfiles:                   ← ADDED (entire block)
                                        │   datasourceUid: pyroscope
                                        │   customQuery: false
                                        │   profileTypeId: "process_cpu:..."
                                        │   tags:
                                        │     - key: service.name
                                        │       value: service_name
                                        │
serviceMap:                             │ serviceMap:
  datasourceUid: prometheus             │   datasourceUid: prometheus
search:                                 │ search:
  hide: false                           │   hide: false
nodeGraph:                              │ nodeGraph:
  enabled: true                         │   enabled: true
lokiSearch:                             │ lokiSearch:
  datasourceUid: loki                   │   datasourceUid: loki
```

**What we added and why:**

| Addition | Why needed | What happens without it |
|---|---|---|
| `spanStartTimeShift: "-1m"` | Widens log search window to 1 min before span start | Default = exact span time range only. Logs that arrive slightly before the span are missed |
| `spanEndTimeShift: "1m"` | Widens log search window to 1 min after span end | Same — logs arriving slightly after are missed |
| `tracesToProfiles` (entire block) | **Not in default image at all.** Links traces → Pyroscope profiles. Enables "Profiles for this span" button | Without it: no profile buttons on any span, even if `pyroscope.profile.id` exists. This was the original issue we discovered |

**Important:** The default `grafana/otel-lgtm` image does NOT include `tracesToProfiles`. You must add it manually to enable trace-to-profile linking.

#### Loki — No changes

```
Default (lgtm image)                    │ Custom (this project)
────────────────────────────────────────│────────────────────────────────────
(identical)                             │ (identical)
```

The Loki datasource with `derivedFields` for log-to-trace linking is identical. Works out of the box.

#### Pyroscope — No changes

```
Default (lgtm image)                    │ Custom (this project)
────────────────────────────────────────│────────────────────────────────────
(identical)                             │ (identical)
```

The Pyroscope datasource is identical. Just name, type, uid, url. No extra config needed.

#### Summary of all differences

| Datasource | Default in lgtm image | Customized in this project |
|---|---|---|
| Prometheus | exemplar → trace linking | **Same** (no change) |
| Tempo / tracesToLogsV2 | Basic log linking (no time shift) | **Added** `spanStartTimeShift`, `spanEndTimeShift` |
| Tempo / tracesToProfiles | **Not present** | **Added** entire block (profile linking) |
| Tempo / serviceMap, search, nodeGraph, lokiSearch | Present | **Same** (no change) |
| Loki | derivedFields for log → trace | **Same** (no change) |
| Pyroscope | Basic datasource definition | **Same** (no change) |

**Bottom line:** Only 3 lines were added to the lgtm default. But the `tracesToProfiles` block is the critical one — without it, the entire trace-to-profile feature doesn't work.

---

### Prometheus Datasource

```yaml
- name: Prometheus                          # Display name in Grafana
  type: prometheus                          # Datasource plugin type
  uid: prometheus                           # Unique ID — referenced by other datasources
  url: http://127.0.0.1:9090               # Inside lgtm container, Prometheus runs on 9090
  editable: true                            # Users can modify in Grafana UI
  jsonData:
    timeInterval: 60s                       # [P1] Default scrape/query interval
    exemplarTraceIdDestinations:            # [P2] Links metric exemplars → traces
      - name: trace_id                      # [P3] Exemplar label containing trace ID
        datasourceUid: tempo                # [P4] Which datasource to open trace in
        urlDisplayLabel: "Trace: $${__value.raw}"  # [P5] Button text shown on exemplars
```

| # | What it controls | If removed/changed |
|---|---|---|
| P1 | Metric query resolution. `60s` means 1 data point per minute | Higher value = less granular metrics. Lower = more load |
| P2-P5 | **Metric → Trace linking.** When a metric has an exemplar with `trace_id`, a link appears to jump to that trace in Tempo | Removed = no "Trace" links on metric exemplars. Metrics and traces become disconnected |
| P4 | Must match Tempo's `uid: tempo` | Wrong UID = exemplar links go to wrong/missing datasource |

### Tempo Datasource

```yaml
- name: Tempo
  type: tempo
  uid: tempo                                # Referenced by Prometheus exemplars and Loki derived fields
  url: http://127.0.0.1:3200               # Inside lgtm, Tempo runs on 3200
  editable: true
  jsonData:
```

#### tracesToLogsV2 — Trace → Log linking

```yaml
    tracesToLogsV2:
      customQuery: true                     # [T1] Use custom LogQL query instead of default
      datasourceUid: loki                   # [T2] Which Loki datasource to query
      query: '{$${__tags}} | trace_id = "$${__trace.traceId}"'  # [T3] LogQL query template
      spanStartTimeShift: "-1m"             # [T4] Search logs from 1 min before span start
      spanEndTimeShift: "1m"                # [T5] Search logs until 1 min after span end
      tags:                                 # [T6] Map span attributes to Loki labels
        - key: service.name                 # [T7] OTel resource attribute name
          value: service_name               # [T8] Corresponding Loki label name
```

| # | What it controls | If removed/changed |
|---|---|---|
| T1 | Enables custom LogQL query. `true` = use `query` field. `false` = Grafana auto-generates query from tags | `false` = simpler but less precise log matching. `true` is needed for trace_id correlation |
| T2 | Must match Loki's `uid: loki` | Wrong UID = "Logs for this span" button navigates to wrong/missing datasource |
| T3 | The LogQL query. `$${__tags}` expands to mapped tags (e.g., `service_name="store-service"`). `$${__trace.traceId}` expands to the current trace ID | Removed = no log correlation. Wrong query = empty results or errors. The `$$` is YAML escaping to produce literal `${...}` in Grafana |
| T4-T5 | Time window padding. Logs may arrive slightly before/after the span | Too narrow = miss logs. Too wide = too many unrelated logs |
| T6-T8 | Maps OTel `service.name` resource attribute to Loki's `service_name` label | Removed = query doesn't filter by service. Wrong mapping = empty results |

**Produces:** "Logs for this span" button on every span in trace view.

#### tracesToProfiles — Trace → Profile linking

```yaml
    tracesToProfiles:
      datasourceUid: pyroscope              # [R1] Which Pyroscope datasource
      customQuery: false                    # [R2] DO NOT set to true — see warning below
      profileTypeId: "process_cpu:cpu:nanoseconds:cpu:nanoseconds"  # [R3] Default profile type
      tags:                                 # [R4] Map span attributes to Pyroscope labels
        - key: service.name                 # [R5] OTel resource attribute
          value: service_name               # [R6] Pyroscope label
```

| # | What it controls | If removed/changed |
|---|---|---|
| R1 | Must match Pyroscope's `uid: pyroscope` | Wrong UID = profile buttons broken |
| R2 | **MUST be `false`.** Setting to `true` hides "Profiles for this span" and "Open in Profiles Drilldown" buttons on ALL services (tested on Grafana 12.4.1). See `monitoring/docs/profiling-constraints.md` | `true` = ALL profile buttons disappear, even for store-service which has `pyroscope.profile.id` |
| R3 | Which profile type to query. `process_cpu` works for store-service (Go). Point-service sends `wall` profiles but this doesn't matter because point-service has no `pyroscope.profile.id` anyway | Changed = wrong profile type shown in flame graph. This only affects store-service since point-service can't link regardless |
| R4-R6 | Maps OTel `service.name` to Pyroscope's `service_name` label | Removed = profile queries don't filter by service |

**Produces:** "Profiles for this span" button and embedded flame graph on spans that have `pyroscope.profile.id` attribute. "Open in Profiles Drilldown" link to navigate to Pyroscope UI.

**Important:** These buttons ONLY appear if the span has `pyroscope.profile.id` attribute. Only store-service (Go) has this, via `otel-profiling-go`. See `monitoring/docs/profiling-constraints.md`.

#### serviceMap — Trace → Metrics linking for service map

```yaml
    serviceMap:
      datasourceUid: prometheus             # [S1] Prometheus datasource for service map metrics
```

| # | What it controls | If removed/changed |
|---|---|---|
| S1 | Service map visualization uses Prometheus metrics (generated from traces by Tempo's metrics-generator) | Removed = Service Graph tab in Tempo shows no data |

#### Other Tempo settings

```yaml
    search:
      hide: false                           # [O1] Show search tab in Tempo datasource
    nodeGraph:
      enabled: true                         # [O2] Enable node graph visualization
    lokiSearch:
      datasourceUid: loki                   # [O3] Enable Loki-based trace search
```

| # | What it controls | If removed/changed |
|---|---|---|
| O1 | `true` = hide search. `false` = show search tab | `true` = users can't search traces by attributes |
| O2 | Enables visual node graph (dependency map) in trace view | `false` = Node graph tab disappears |
| O3 | Allows searching traces via Loki logs (alternative to Tempo's native search) | Removed = can't find traces through log correlation |

### Loki Datasource

```yaml
- name: Loki
  type: loki
  uid: loki                                 # Referenced by Tempo's tracesToLogsV2 and lokiSearch
  url: http://127.0.0.1:3100               # Inside lgtm, Loki runs on 3100
  editable: true
  jsonData:
    derivedFields:                          # [L1] Extract fields from log lines → create links
      - name: trace_id                      # [L2] Derived field name
        matcherType: label                  # [L3] Match by label (not regex on log body)
        matcherRegex: trace_id              # [L4] Label name to extract from
        url: "$${__value.raw}"              # [L5] URL value (the trace ID itself)
        datasourceUid: tempo                # [L6] Which datasource to link to
        urlDisplayLabel: "Trace: $${__value.raw}"  # [L7] Link text shown on log lines
```

| # | What it controls | If removed/changed |
|---|---|---|
| L1-L7 | **Log → Trace linking.** When a log line has a `trace_id` label, a clickable link appears to jump to that trace in Tempo | Removed = no "Trace: ..." links on log lines. Logs and traces become disconnected |
| L3-L4 | `label` type means it reads from Loki labels (not regex on log body). OTel auto-instrumentation adds `trace_id` as a label | Changed to `regex` = needs regex pattern on log body instead |
| L6 | Must match Tempo's `uid: tempo` | Wrong UID = trace links go to wrong datasource |

**Produces:** "Trace: abc123..." clickable links on log lines in Grafana Explore > Loki.

### Pyroscope Datasource

```yaml
- name: Pyroscope
  type: grafana-pyroscope-datasource
  uid: pyroscope                            # Referenced by Tempo's tracesToProfiles
  url: http://127.0.0.1:4040               # Inside lgtm, Pyroscope runs on 4040
```

| Field | If removed/changed |
|---|---|
| `uid: pyroscope` | Must match `tracesToProfiles.datasourceUid`. Wrong = profile buttons link to nothing |
| `url` | Wrong URL = Pyroscope queries fail. Flame graphs show "No data" |

---

## 3. OTel Collector Config — `monitoring/otel-collector/config.yaml`

The OTel Collector runs as a separate container, NOT inside lgtm. It collects MySQL-specific telemetry.

```yaml
receivers:
  mysql:                                     # [C1] MySQL metrics receiver
    endpoint: db:3306                        # [C2] MySQL host:port
    username: user                           # [C3] MySQL credentials
    password: password
    collection_interval: 15s                 # [C4] How often to scrape metrics
    statement_events:
      digest_text_limit: 120                 # [C5] Max length of SQL statement text
      time_limit: 24h                        # [C6] Look back period for statement events
      limit: 250                             # [C7] Max number of statement events
    metrics:                                 # [C8] Which MySQL metrics to enable
      mysql.commands:          { enabled: true }
      mysql.query.count:       { enabled: true }
      mysql.query.slow.count:  { enabled: true }
      mysql.connection.count:  { enabled: true }
      mysql.connection.errors: { enabled: true }
      mysql.joins:             { enabled: true }
      mysql.table_open_cache:  { enabled: true }
      mysql.statement_event.count:      { enabled: true }
      mysql.statement_event.wait.time:  { enabled: true }

  filelog/slow_query:                        # [C9] File-based log receiver
    include:
      - /var/log/mysql/slow.log              # [C10] Path to MySQL slow query log
    multiline:
      line_start_pattern: '^# Time:'         # [C11] Each slow query entry starts with "# Time:"
```

| # | What it controls | If removed/changed |
|---|---|---|
| C1-C4 | MySQL metrics collection. Appears in Grafana as `mysql.*` metrics | Removed = no MySQL database metrics. Only application-level metrics remain |
| C5-C7 | Statement event tracking (top SQL queries) | Removed = no per-query metrics |
| C8 | Each metric can be toggled. `false` = that metric stops being collected | Specific metrics disappear from dashboards |
| C9-C11 | Reads MySQL slow query log file and sends to Loki | Removed = no slow query logs in Grafana. The `mysql-logs` volume mount must also exist |

```yaml
processors:
  resource/slow_query:                       # [C12] Adds resource attributes to slow query logs
    attributes:
      - key: service.name                    # [C13] Tags logs as coming from "mysql" service
        value: mysql
        action: upsert
      - key: source                          # [C14] Tags logs with source identifier
        value: mysql-slow-query
        action: upsert

exporters:
  otlp:                                      # [C15] OTLP exporter to lgtm
    endpoint: lgtm:4317                      # [C16] Same endpoint as services use
    tls:
      insecure: true                         # [C17] No TLS (internal network)

service:
  pipelines:
    metrics:                                 # [C18] Metrics pipeline: mysql → otlp
      receivers: [mysql]
      exporters: [otlp]
    logs:                                    # [C19] Logs pipeline: filelog → resource processor → otlp
      receivers: [filelog/slow_query]
      processors: [resource/slow_query]
      exporters: [otlp]
```

| # | What it controls | If removed/changed |
|---|---|---|
| C12-C14 | Without `service.name: mysql`, slow query logs appear with no service name in Loki | Loki queries like `{service_name="mysql"}` won't match |
| C15-C17 | Where to send data. `insecure: true` because it's Docker internal networking | Removed = collector receives data but doesn't forward it. `insecure: false` = TLS handshake fails |
| C18-C19 | Pipeline definitions. A receiver not listed in any pipeline is inactive | Removed pipeline = that signal stops flowing |

---

## 4. store-service OTel Setup — `store-service/internal/otel/otel.go`

```go
func InitOtel(ctx context.Context) (func(), error) {
    // [G1] Create resource with env var attributes + SDK info
    res, err := resource.New(ctx,
        resource.WithFromEnv(),          // reads OTEL_SERVICE_NAME, OTEL_RESOURCE_ATTRIBUTES
        resource.WithTelemetrySDK(),     // adds telemetry.sdk.name, telemetry.sdk.language, telemetry.sdk.version
    )

    // [G2] Protocol selection from env var
    protocol := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")

    // [G3] Create exporters based on protocol (grpc or http)
    // gRPC: otlptracegrpc, otlpmetricgrpc, otlploggrpc
    // HTTP: otlptracehttp, otlpmetrichttp, otlploghttp

    // [G4] TracerProvider — batches and exports spans
    tp := trace.NewTracerProvider(
        trace.WithBatcher(traceExporter),
        trace.WithResource(res),
    )
    otel.SetTracerProvider(tp)

    // [G5] MeterProvider — exports metrics every 30s
    mp := metric.NewMeterProvider(
        metric.WithReader(metric.NewPeriodicReader(metricExporter,
            metric.WithInterval(30*time.Second))),
        metric.WithResource(res),
    )
    otel.SetMeterProvider(mp)

    // [G6] Context propagation — W3C TraceContext + Baggage
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))

    // [G7] LoggerProvider — batches and exports logs via OTel
    lp := otellog.NewLoggerProvider(...)
    handler := otelslog.NewHandler("store-service", ...)
    slog.SetDefault(slog.New(handler))       // [G8] All slog calls go through OTel
}
```

| # | What it controls | If removed/changed |
|---|---|---|
| G1 | Resource identity. All signals share this. `service.name` comes from `OTEL_SERVICE_NAME` env var | Removed `WithFromEnv()` = service name not set. All telemetry is "unknown_service" |
| G2-G3 | Protocol switch. Allows changing between gRPC (:4317) and HTTP (:4318) without code changes | Defaults to grpc if env var missing |
| G4 | Trace export. `WithBatcher` = async batching for performance | Removed = no traces from store-service |
| G5 | Metric export. 30s interval = metrics sent every 30 seconds | Removed = no metrics. Longer interval = less granular but lower overhead |
| G6 | **Critical for distributed tracing.** `TraceContext` propagates trace/span IDs via HTTP headers (`traceparent`). Without this, each service creates independent traces | Removed = traces break across service boundaries. store-service and point-service traces are separate, not connected |
| G7-G8 | Routes Go's `slog` logger through OTel to Loki. Every `slog.Info()`, `slog.Error()` call becomes an OTel log record | Removed = logs go to stdout only, not to Loki |

### store-service Pyroscope + OTel Bridge — `cmd/main.go`

```go
// [G9] Condition: only if PYROSCOPE_URL is set
if pyroscopeURL != "" {
    // [G10] Start Pyroscope profiler — sends profiles to :4040
    profiler, err := profiling.StartPyroscope(pyroscopeURL, serviceName)
    // [G11] Wrap TracerProvider with otel-profiling-go
    otel.SetTracerProvider(otelpyroscope.NewTracerProvider(otel.GetTracerProvider()))
}
```

| # | What it controls | If removed/changed |
|---|---|---|
| G9 | Guards profiling init. No env var = no profiling | Always present when running in Docker Compose |
| G10 | Starts CPU + memory + goroutine + mutex + block profiling. Pushes to Pyroscope every 10s | Removed = no profile data for store-service in Pyroscope |
| G11 | **This is the key line for span-profile linking.** Wraps the existing TracerProvider. On every span, it: (1) tags Pyroscope samples with the span ID, (2) adds `pyroscope.profile.id` attribute to the span. This attribute is what Grafana uses to show "Profiles for this span" button | Removed = `pyroscope.profile.id` disappears from all spans. Profile buttons vanish. Profiling still works (flame graphs by service) but no span-level linking |

---

## 5. point-service OTel Setup — `point-service/src/trace.ts`

```typescript
// [N1] Pyroscope init — must happen before OTel SDK starts
let Pyroscope: any = null
const pyroscopeUrl = process.env.PYROSCOPE_URL
if (pyroscopeUrl) {
    Pyroscope = require('@pyroscope/nodejs')     // [N2] CJS require (needs p-limit override)
    Pyroscope.init({
        serverAddress: pyroscopeUrl,
        appName: 'point-service',
        wall: { collectCpuTime: true },          // [N3] Wall-time profiling (not CPU sampling)
    })
    Pyroscope.startWallProfiling()                // [N4] Start wall profiler
    Pyroscope.startHeapProfiling()                // [N5] Start heap profiler
}
```

| # | What it controls | If removed/changed |
|---|---|---|
| N1 | Pyroscope must init before OTel SDK to avoid instrumentation conflicts | Swapping order may cause Pyroscope HTTP calls to be instrumented |
| N2 | `require()` = CJS. Needs `p-limit` override to v3.1.0 in package.json (v4+ is ESM-only) | Without override: `ERR_REQUIRE_ESM` crash. See package.json `overrides` section |
| N3 | Wall-time profiling captures all time (CPU + I/O wait). Different from Go's CPU profiling | Profile type in Pyroscope is `wall`, not `process_cpu`. Need to select "wall" in Grafana dropdown |
| N4-N5 | Start continuous profiling | Removed = no profile data for point-service |

```typescript
// [N6] OTel SDK setup
export const otelSDK = new NodeSDK({
    // [N7] Trace exporter with batching
    spanProcessors: [new BatchSpanProcessor(grpcExporter, {
        maxQueueSize: 2048,
        maxExportBatchSize: 512,
        scheduledDelayMillis: 5000,
        exportTimeoutMillis: 30000,
    }) as any],                              // [N8] `as any` needed due to OTel version drift

    // [N9] Log exporter
    logRecordProcessors: [new BatchLogRecordProcessor(logExporter)],

    // [N10] Auto-instrumentations
    instrumentations: [
        getNodeAutoInstrumentations({
            '@opentelemetry/instrumentation-http': {
                enabled: true,
                ignoreOutgoingRequestHook: (request) => {
                    const host = request.hostname || request.host || '';
                    return host.includes('lgtm');        // [N11] Ignore HTTP calls to lgtm
                },
            },
            '@opentelemetry/instrumentation-dns': { enabled: false },   // [N12]
            '@opentelemetry/instrumentation-net': { enabled: false },   // [N13]
            '@opentelemetry/instrumentation-undici': {
                enabled: true,
                ignoreRequestHook: (request) => {
                    const origin = request.origin || '';
                    return origin.includes('lgtm');      // [N14] Ignore undici calls to lgtm
                },
            },
            '@opentelemetry/instrumentation-mysql2': { enabled: true },
            '@opentelemetry/instrumentation-pg': { enabled: true },
        }),
    ],
    serviceName: 'point-service',                        // [N15]
});

otelSDK.start();                                          // [N16] Must start before NestJS imports
```

| # | What it controls | If removed/changed |
|---|---|---|
| N7 | Batch span processor with custom queue size. `5000ms` delay = sends every 5s | Removed = no traces. Smaller batch = more network calls |
| N8 | TypeScript type cast. Without it: `BatchSpanProcessor` not assignable to `SpanProcessor` due to OTel package version drift | Remove cast = TypeScript compilation error |
| N11 | Prevents Pyroscope's profile push HTTP calls from being captured as spans | Removed = fragment traces appear (dns.lookup, tcp.connect, POST to `/ingest`) that pollute Tempo |
| N12-N13 | DNS and Net instrumentations disabled to prevent noise from Pyroscope SDK's internal connections | Enabled = dns.lookup and tcp.connect spans appear as separate traces |
| N14 | **Critical.** Node.js 20 uses `undici` internally for HTTP. Pyroscope SDK calls go through undici, not `http` module. This hook filters them out | Removed = POST spans to `lgtm:4040/ingest` appear as separate traces in Tempo |
| N15 | Service name. Overrides any env var | Removed = service appears as "unknown_service:node" |
| N16 | **Must happen at module load time** (before NestJS imports Express). `trace.ts` is imported first in `main.ts` | If SDK starts after NestJS, HTTP/Express instrumentation doesn't work — no automatic spans for incoming requests |

---

## 6. Cross-Linking Summary

This shows how all the pieces connect. Each arrow is a navigable link in Grafana:

```
Prometheus metrics ──(exemplar trace_id)──→ Tempo traces
                                              │
Tempo traces ──(tracesToLogsV2)──────────────→ Loki logs
Tempo traces ──(tracesToProfiles)────────────→ Pyroscope profiles
Tempo traces ──(serviceMap)──────────────────→ Prometheus (service graph)
                                              │
Loki logs ──(derivedFields trace_id)─────────→ Tempo traces
```

**If you remove any single datasource UID or cross-link config, the corresponding navigation link breaks.**

### Quick reference: what controls each Grafana button

| Button/Link | Where it appears | Config that controls it | Data requirement |
|---|---|---|---|
| "Logs for this span" | Trace span detail | `tracesToLogsV2` in Tempo DS | Logs must exist in Loki with matching `trace_id` label |
| "Profiles for this span" | Trace span detail | `tracesToProfiles` in Tempo DS | Span must have `pyroscope.profile.id` attribute |
| "Open in Profiles Drilldown" | Trace span detail | `tracesToProfiles` in Tempo DS | Same as above |
| "Trace: abc123..." | Log line in Loki | `derivedFields` in Loki DS | Log must have `trace_id` label |
| Exemplar trace link | Metric graph | `exemplarTraceIdDestinations` in Prometheus DS | Metric must have exemplar with `trace_id` |
| Service Graph | Tempo Service Graph tab | `serviceMap` in Tempo DS | Tempo metrics-generator must produce span metrics |

---

## References

- [Grafana Tempo data source configuration](https://grafana.com/docs/grafana/latest/datasources/tempo/configure-tempo-data-source/)
- [Configure traces to profiles](https://grafana.com/docs/grafana/latest/datasources/pyroscope/configure-traces-to-profiles/)
- [OpenTelemetry Collector MySQL receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/mysqlreceiver)
- [Grafana Beyla eBPF](https://grafana.com/docs/beyla/latest/)
- [otel-profiling-go](https://github.com/grafana/otel-profiling-go) — Go bridge for span-profile linking
- [Profiling constraints](./profiling-constraints.md) — Why point-service cannot link traces to profiles
