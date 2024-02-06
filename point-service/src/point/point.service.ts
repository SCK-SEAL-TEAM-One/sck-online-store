import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Point } from './point.entity';
import { CreatePointDto } from './point.dto';

@Injectable()
export class PointService {
  constructor(
    @InjectRepository(Point)
    private pointRepository: Repository<Point>,
  ) {}

  async getPoint(): Promise<Point[]> {
    return await this.pointRepository.find();
  }

  async deductPoint(point: CreatePointDto): Promise<Point> {
    return await this.pointRepository.save(point);
  }
}
