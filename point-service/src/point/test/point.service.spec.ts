
import { Test, TestingModule } from '@nestjs/testing';
import { getRepositoryToken } from '@nestjs/typeorm';
import { CreatePointDto } from '../point.dto';
import { Point } from '../point.entity';
import { PointService } from '../point.service';

describe('PointService', () => {
  let service: PointService;

  const mockPointRepository = {
    save: jest.fn(),
    find: jest.fn(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        PointService,
        {
          provide: getRepositoryToken(Point),
          useValue: mockPointRepository,
        },
      ],
    }).compile();

    service = module.get<PointService>(PointService);
  });

  it('Should be defined', () => {
    expect(service).toBeDefined();
  });

  it('Create => Should create a new point and return its data', async () => {
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

    jest.spyOn(mockPointRepository, 'save').mockReturnValue(createPointResponse);

    // act
    const result = await service.deductPoint(createPointInput);

    // assert
    expect(mockPointRepository.save).toBeCalled();
    expect(mockPointRepository.save).toBeCalledWith(createPointInput);
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

    jest.spyOn(mockPointRepository, 'find').mockReturnValue(points);

    //act
    const result = await service.getPoint();

    // assert
    expect(result).toEqual(points);
    expect(mockPointRepository.find).toBeCalled();
  });

});