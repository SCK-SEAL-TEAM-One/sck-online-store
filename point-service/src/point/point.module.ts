import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { PointService } from './point.service';
import { PointController } from './point.controller';
import { Point } from './point.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Point])],
  providers: [PointService],
  controllers: [PointController],
})
export class PointModule {}
