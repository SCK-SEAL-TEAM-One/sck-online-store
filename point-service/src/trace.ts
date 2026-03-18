import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { BatchLogRecordProcessor } from '@opentelemetry/sdk-logs';
import { NodeSDK } from '@opentelemetry/sdk-node';
import * as process from 'process';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';

import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { OTLPLogExporter } from '@opentelemetry/exporter-logs-otlp-grpc';

// Pyroscope profiler init — must be before OTel SDK start
let Pyroscope: any = null
const pyroscopeUrl = process.env.PYROSCOPE_URL
if (pyroscopeUrl) {
  try {
    Pyroscope = require('@pyroscope/nodejs')
    if (Pyroscope.default) Pyroscope = Pyroscope.default
    Pyroscope.init({
      serverAddress: pyroscopeUrl,
      appName: 'point-service',
      wall: {
        collectCpuTime: true,
      },
    })
    Pyroscope.startWallProfiling()
    Pyroscope.startHeapProfiling()
    console.log(`[Pyroscope] Profiler started, sending to ${pyroscopeUrl}`)
  } catch (err) {
    console.warn('[Pyroscope] Failed to initialize profiler:', err.message)
    Pyroscope = null
  }
}

const grpcExporter = new OTLPTraceExporter({
  headers: {},
  concurrencyLimit: 10,
  timeoutMillis: 5000,
});

const logExporter = new OTLPLogExporter({});

console.log(
  '[OTEL] Initializing OpenTelemetry SDK for point-service with auto-instrumentations',
);

export const otelSDK = new NodeSDK({
  spanProcessor: new BatchSpanProcessor(grpcExporter, {
    maxQueueSize: 2048,
    maxExportBatchSize: 512,
    scheduledDelayMillis: 5000,
    exportTimeoutMillis: 30000,
  }),

  logRecordProcessors: [new BatchLogRecordProcessor(logExporter)],

  instrumentations: [
    getNodeAutoInstrumentations({
      '@opentelemetry/instrumentation-http': {
        enabled: true,
      },
      '@opentelemetry/instrumentation-mysql2': {
        enabled: true,
      },
      '@opentelemetry/instrumentation-pg': {
        enabled: true,
      },
    }),
  ],
  serviceName: 'point-service',
});

console.log(
  '[OTEL] SDK configured with auto-instrumentations (includes HTTP, MySQL, PostgreSQL, etc.)',
);

// Start SDK immediately at module load time so http/express are patched
// BEFORE NestJS imports them
otelSDK.start();
console.log('[OTEL] SDK started');

// gracefully shut down the SDK on process exit
process.on('SIGTERM', () => {
  console.log('[OTEL] Shutting down OpenTelemetry SDK');
  if (Pyroscope) {
    Pyroscope.stopWallProfiling();
    Pyroscope.stopHeapProfiling();
    console.log('[Pyroscope] Profiler stopped');
  }
  otelSDK
    .shutdown()
    .then(
      () => console.log('[OTEL] OpenTelemetry SDK shut down successfully'),
      (err) => console.log('[OTEL] Error shutting down OpenTelemetry SDK', err),
    )
    .finally(() => process.exit(0));
});
