import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { otelSDK } from './trace';

async function bootstrap() {
    await otelSDK.start();
  const app = await NestFactory.create(AppModule);
  app.setGlobalPrefix('api/v1');
  await app.listen(8001);
}
bootstrap();
