package point

import (
	"context"
	"fmt"
	"log"
)

type PointInterface interface {
	TotalPoint(ctx context.Context, uid int) (TotalPoint, error)
	DeductPoint(ctx context.Context, uid int, submitedPoint SubmitedPoint) (TotalPoint, error)
	CheckBurnPoint(ctx context.Context, uid int, amount int) (bool, error)
}

type PointService struct {
	PointGateway PointGatewayInterface
}

type PointGatewayInterface interface {
	GetPoints(ctx context.Context, uid int) ([]Point, error)
	CreatePoint(ctx context.Context, uid int, body Point) (Point, error)
}

func (pointService PointService) TotalPoint(ctx context.Context, uid int) (TotalPoint, error) {
	points, err := pointService.PointGateway.GetPoints(ctx, uid)
	if err != nil {
		log.Printf("pointService.PointGateway.GetPoints internal error %s", err.Error())
	}

	total := 0
	for _, point := range points {
		total += point.Amount
	}
	return TotalPoint{
		Point: total,
	}, err
}

func (pointService PointService) DeductPoint(ctx context.Context, uid int, submitedPoint SubmitedPoint) (TotalPoint, error) {
	_, err := pointService.CheckBurnPoint(ctx, uid, submitedPoint.Amount)
	if err != nil {
		return TotalPoint{}, err
	}

	point := Point{
		OrgID:  1,
		UserID: 1,
		Amount: submitedPoint.Amount,
	}
	_, err_ := pointService.PointGateway.CreatePoint(ctx, uid, point)
	if err_ != nil {
		log.Printf("pointService.PointGateway.CreatePoint internal error %s", err_.Error())
		return TotalPoint{}, err_
	}
	return pointService.TotalPoint(ctx, uid)
}

func (pointService PointService) CheckBurnPoint(ctx context.Context, uid int, amount int) (bool, error) {
	total, err := pointService.TotalPoint(ctx, uid)
	if err != nil {
		log.Printf("pointService.TotalPoint internal error %s", err.Error())
		return false, err
	}
	if amount+total.Point < 0 {
		return false, fmt.Errorf("points are not enough, please try again")
	}
	return true, nil
}
