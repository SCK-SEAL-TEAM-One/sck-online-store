import './trace'; // Must be first — starts OTEL SDK before http/express are loaded
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  app.setGlobalPrefix('api/v1');
  await app.listen(8001);
}
bootstrap();
