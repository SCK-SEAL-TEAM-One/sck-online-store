package point

import (
	"fmt"
	"log"
)

type PointInterface interface {
	TotalPoint(uid int) (TotalPoint, error)
	DeductPoint(uid int, submitedPoint SubmitedPoint) (TotalPoint, error)
	CheckBurnPoint(uid int, amount int) (bool, error)
}

type PointService struct {
	PointGateway PointGatewayInterface
}

type PointGatewayInterface interface {
	GetPoints(uid int) ([]Point, error)
	CreatePoint(uid int, body Point) (Point, error)
}

func (pointService PointService) TotalPoint(uid int) (TotalPoint, error) {
	points, err := pointService.PointGateway.GetPoints(uid)
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

func (pointService PointService) DeductPoint(uid int, submitedPoint SubmitedPoint) (TotalPoint, error) {
	_, err := pointService.CheckBurnPoint(uid, submitedPoint.Amount)
	if err != nil {
		return TotalPoint{}, err
	}

	point := Point{
		OrgID:  1,
		UserID: 1,
		Amount: submitedPoint.Amount,
	}
	_, err_ := pointService.PointGateway.CreatePoint(uid, point)
	if err_ != nil {
		log.Printf("pointService.PointGateway.CreatePoint internal error %s", err.Error())
		return TotalPoint{}, err_
	}
	return pointService.TotalPoint(uid)
}

func (pointService PointService) CheckBurnPoint(uid int, amount int) (bool, error) {
	total, err := pointService.TotalPoint(uid)
	if err != nil {
		log.Printf("pointService.TotalPoint internal error %s", err.Error())
		return false, err
	}
	if amount+total.Point < 0 {
		return false, fmt.Errorf("points are not enough, please try again")
	}
	return true, nil
}
