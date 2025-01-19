import { SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { NodeSDK } from '@opentelemetry/sdk-node';
import * as process from 'process';
import { HttpInstrumentation } from '@opentelemetry/instrumentation-http';
import { NestInstrumentation } from '@opentelemetry/instrumentation-nestjs-core';
import { TypeormInstrumentation } from 'opentelemetry-instrumentation-typeorm';

import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';

const collectorOptions = {
  url: 'http://lgtm:4317', // url is optional and can be omitted - default is http://localhost:4318/v1/traces
  headers: {}, // an optional object containing custom headers to be sent with each request
  concurrencyLimit: 10, // an optional limit on pending requests
};

const exporter = new OTLPTraceExporter(collectorOptions);

export const otelSDK = new NodeSDK({
  spanProcessor: new SimpleSpanProcessor(exporter),

  instrumentations: [
    new HttpInstrumentation(),
    new NestInstrumentation(),
    new TypeormInstrumentation(),
  ],
  serviceName: 'point-service',
});

// gracefully shut down the SDK on process exit
process.on('SIGTERM', () => {
  otelSDK
    .shutdown()
    .then(
      () => console.log('SDK shut down successfully'),
      (err) => console.log('Error shutting down SDK', err),
    )
    .finally(() => process.exit(0));
});
