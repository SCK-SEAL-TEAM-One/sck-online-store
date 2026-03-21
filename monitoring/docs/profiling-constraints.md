# Profiling Constraints & Known Limitations

> **Last verified:** 2026-03-21
> **Status:** Active constraint — revisit when OpenTelemetry Profiling signal reaches stable or Pyroscope releases Node.js OTel bridge.

## Summary

Span-level profile linking ("Profiles for this span" and "Open in Profiles Drilldown" buttons in Grafana Tempo) works for **store-service (Go)** but **NOT for point-service (Node.js)**. This is a known upstream limitation, not a bug in our code.

## Root Cause

The feature requires a `pyroscope.profile.id` span attribute injected into every OTel span. This attribute links a specific trace span to a specific profile sample in Pyroscope.

| Service | Language | Profiling SDK | OTel Bridge | `pyroscope.profile.id` | Span-Profile Link |
|---------|----------|---------------|-------------|------------------------|-------------------|
| store-service | Go | `grafana/pyroscope-go` | `grafana/otel-profiling-go` | Yes (auto-injected) | Works |
| point-service | Node.js | `@pyroscope/nodejs` | **Does not exist** | No | **Does not work** |

### How store-service gets it (Go)

In `store-service/cmd/main.go`:
```go
import otelpyroscope "github.com/grafana/otel-profiling-go"

// Wraps TracerProvider to auto-inject pyroscope.profile.id into all spans
otel.SetTracerProvider(otelpyroscope.NewTracerProvider(otel.GetTracerProvider()))
```

The `otel-profiling-go` library wraps the OTel TracerProvider. On every span start/end, it:
1. Tags the Pyroscope profile sample with the span ID
2. Adds `pyroscope.profile.id` attribute to the span
3. This allows Grafana to correlate a specific span with its profile data

### Why point-service doesn't have it (Node.js)

Grafana/Pyroscope has released OTel bridge libraries for Go, Java, Ruby, .NET, and Python — but **NOT for Node.js**. The `@pyroscope/nodejs` SDK collects wall-time profiles and sends them to Pyroscope, but does not integrate with OpenTelemetry spans.

## What works for point-service today

- **Continuous profiling** works — `@pyroscope/nodejs` sends wall-time profiles to Pyroscope
- **Flame graphs** are visible in Grafana Drilldown > Profiles > select `point-service` > profile type `wall`
- The CPU-intensive `calculatePointsReward()` in `point.service.ts` is clearly visible in flame graphs
- To view point-service profiles: navigate manually to Drilldown > Profiles > select `point-service` > profile type `wall`

## What does NOT work for point-service

- **"Profiles for this span" button** — does not appear (requires `pyroscope.profile.id`)
- **"Open in Profiles Drilldown" button** — does not appear on point-service spans
- **Span-level profile filtering** — cannot isolate the profile of a single request/span

## OpenTelemetry Profiling Signal Status

Profiling is the 4th signal in OpenTelemetry (after Traces, Metrics, Logs) but is still in **development/experimental** phase:

| Milestone | Date | Status |
|-----------|------|--------|
| Profiling SIG started | May 2022 | Done |
| Data model proposed (v1, v2) | Sep-Nov 2023 | Done |
| Profiling announced as 4th signal | Jun 2024 | Done |
| Added to OTLP protocol (v1.3.0) | 2024 | Done (unstable) |
| Collector support (v0.112.0) | 2024 | Experimental (feature gate required) |
| Stable specification | TBD | Not yet |
| Language SDK support | TBD | Not yet (no SDK has stable profiling) |

When the OTel SDK natively supports profiling, vendor bridges like `otel-profiling-go` will become unnecessary — span-profile correlation will be built into the SDK itself.

## When to revisit this constraint

Check these periodically and implement if available:

1. **Pyroscope releases Node.js OTel bridge** — Look for an npm package like `@pyroscope/otel-nodejs` or similar. If released:
   - Install the package in point-service
   - Wrap the OTel TracerProvider (similar to Go's `otelpyroscope.NewTracerProvider()`)
   - `pyroscope.profile.id` will auto-appear in spans
   - All Grafana buttons will work

2. **OpenTelemetry Profiling signal reaches stable** — Check https://opentelemetry.io/docs/specs/status/
   - When profiling API/SDK is stable, upgrade `@opentelemetry/sdk-node`
   - Native span-profile correlation may replace vendor bridges entirely

3. **Grafana changes Tempo UI** — The `tracesToProfiles` datasource config may evolve to support service-level linking without requiring `pyroscope.profile.id`

## Grafana tracesToProfiles config constraint

**DO NOT set `customQuery: true`** in `tracesToProfiles`. Setting `customQuery: true` hides the "Profiles for this span" and "Open in Profiles Drilldown" buttons for ALL services (including store-service which has `pyroscope.profile.id`). This was tested on Grafana 12.4.1.

Keep `customQuery: false` with `profileTypeId: "process_cpu:cpu:nanoseconds:cpu:nanoseconds"`. This works for store-service (Go) which sends `process_cpu` profiles. Point-service (Node.js) sends `wall` profiles so the buttons won't appear regardless (no `pyroscope.profile.id`).

## Related files

- `store-service/cmd/main.go` — Go OTel bridge setup (`otelpyroscope.NewTracerProvider`)
- `store-service/internal/profiling/Profiling.go` — Pyroscope Go SDK init
- `point-service/src/trace.ts` — Node.js OTel + Pyroscope setup (no bridge)
- `point-service/package.json` — `@pyroscope/nodejs` dependency + `p-limit` override
- `monitoring/grafana/provisioning/datasources.yml` — `tracesToProfiles` config with `customQuery: false`
- `deploy/terraform/lgtm-stack.tf` — Same config for EKS deployment

## References

Sources used to confirm this constraint:

### OpenTelemetry Profiling Signal Status
- [OpenTelemetry announces support for profiling (Jun 2024)](https://opentelemetry.io/blog/2024/profiling/) — Official announcement of profiling as the 4th signal
- [The State of Profiling (Oct 2024)](https://opentelemetry.io/blog/2024/state-profiling/) — Progress update: OTLP v1.3.0, Collector v0.112.0, still unstable/experimental
- [OpenTelemetry Specification Status Summary](https://opentelemetry.io/docs/specs/status/) — Profiles protocol listed as "Development" phase
- [OpenTelemetry Project Roadmap](https://opentelemetry.io/community/roadmap/) — Profiling listed as P2 priority, started May 2022, no stable timeline

### Pyroscope Span-Profile Integration (per language)
- [otel-profiling-go (Go)](https://github.com/grafana/otel-profiling-go) — Go OTel bridge that injects `pyroscope.profile.id` into spans
- [Span profiles for Java](https://grafana.com/docs/pyroscope/latest/configure-client/trace-span-profiles/java-span-profiles/) — Java OTel bridge docs
- [pyroscope-otel (Python) on PyPI](https://pypi.org/project/pyroscope-otel/) — Python OTel bridge (v1.0.0 released Feb 2026)
- [Pyroscope.OpenTelemetry (.NET) on NuGet](https://www.nuget.org/packages/Pyroscope.OpenTelemetry) — .NET OTel bridge

### Node.js — No OTel Bridge Exists
- [Pyroscope Node.js SDK docs](https://grafana.com/docs/pyroscope/latest/configure-client/language-sdks/nodejs/) — Documents `@pyroscope/nodejs` SDK but has NO mention of OTel span integration; the "Link traces and profiles" section lists Go, Java, Ruby, .NET — **Node.js is absent**
- [@pyroscope/nodejs on npm](https://www.npmjs.com/package/@pyroscope/nodejs) — SDK package; no OTel TracerProvider wrapper functionality

### Grafana Tempo tracesToProfiles Config
- [Grafana Tempo data source docs](https://grafana.com/docs/grafana/latest/datasources/tempo/configure-tempo-data-source/) — Documents `tracesToProfiles` configuration including `customQuery`, `profileTypeId`, and tag-based linking

## Related investigation

See also: `monitoring/docs/span-profiles-investigation-en.md` for the detailed investigation of the "0 samples" issue for Go spans (a different but related profiling issue).
