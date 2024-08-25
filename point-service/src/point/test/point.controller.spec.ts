
import { Test, TestingModule } from '@nestjs/testing';
import { CreatePointDto } from '../point.dto';
import { PointController } from '../point.controller';
import { PointService } from '../point.service';

describe('PointController', () => {
  let controller: PointController;

  const mockPointService = {
    getPoint: jest.fn(),
    deductPoint: jest.fn(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [PointController],
      providers: [
        {
          provide: PointService,
          useValue: mockPointService,
        },
      ],
    }).compile();

    controller = module.get<PointController>(PointController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });

  it('Create => should create a new point by a given data', async () => {
    // arrange
    const createPointInput = {
      orgId: 1,
      userId: 1,
      amount: 200,
    } as CreatePointDto;

    const createPointResponse = {
      id: 1,
      orgId: 1,
      userId: 1,
      amount: 200,
      created: '2024-08-25T09:06:58',
      updated: '2024-08-25T09:06:58',
    } as CreatePointDto;


    jest.spyOn(mockPointService, 'deductPoint').mockReturnValue(createPointResponse);

    // act
    const result = await controller.createPoint(createPointInput);

    // assert
    expect(mockPointService.deductPoint).toBeCalled();
    expect(mockPointService.deductPoint).toBeCalledWith(createPointInput);
    expect(result).toEqual(createPointResponse);
  });

  it('Find => should return an array of point', async () => {
    //arrange
    const point = {
      id: 2,
      orgId: 1,
      userId: 1,
      amount: 300,
      created: '2024-08-25T09:06:58',
      updated: '2024-08-25T09:06:58',
    };
    const points = [point];
    jest.spyOn(mockPointService, 'getPoint').mockReturnValue(points);

    //act
    const result = await controller.getPoint();

    // assert
    expect(result).toEqual(points);
    expect(mockPointService.getPoint).toBeCalled();
  });

});