import { Body, Controller, Get, Post } from '@nestjs/common';
import { PointService } from './point.service';
import { CreatePointDto } from './point.dto';

@Controller('point')
export class PointController {
  constructor(private readonly pointService: PointService) {}

  @Get()
  getPoint() {
    return this.pointService.getPoint();
  }

  @Post()
  createPoint(@Body() body: CreatePointDto) {
    return this.pointService.deductPoint(body);
  }
}
