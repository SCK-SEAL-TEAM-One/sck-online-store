import { Module } from '@nestjs/common';
import { PointModule } from './point/point.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Point } from './point/point.entity';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [
    ConfigModule.forRoot(),
    TypeOrmModule.forRoot({
      type: 'mysql',
      host: process.env.DB_HOST,
      port: Number(process.env.DB_PORT),
      username: process.env.DB_USERNAME,
      password: process.env.DB_PASSWORD,
      database: 'point',
      entities: [Point],
      synchronize: true,
    }),
    PointModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}
