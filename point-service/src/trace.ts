import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { BatchLogRecordProcessor } from '@opentelemetry/sdk-logs';
import { NodeSDK } from '@opentelemetry/sdk-node';
import * as process from 'process';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';

import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { OTLPLogExporter } from '@opentelemetry/exporter-logs-otlp-grpc';

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
  otelSDK
    .shutdown()
    .then(
      () => console.log('[OTEL] OpenTelemetry SDK shut down successfully'),
      (err) => console.log('[OTEL] Error shutting down OpenTelemetry SDK', err),
    )
    .finally(() => process.exit(0));
});
