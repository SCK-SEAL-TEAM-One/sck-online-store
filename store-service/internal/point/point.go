package point

import (
	"fmt"
	"log"
)

type PointGatewayInterface interface {
	GetPoints(uid int) ([]PointGatewayResponseItem, error)
	CreatePoint(uid int, body Point) (PointGatewayResponseItem, error)
}

type PointService struct {
	PointGateway PointGatewayInterface
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
	total, err := pointService.TotalPoint(uid)
	if err != nil {
		log.Printf("pointService.TotalPoint internal error %s", err.Error())
		return TotalPoint{}, err
	}

	if submitedPoint.Amount+total.Point < 0 {
		return TotalPoint{}, fmt.Errorf("points are not enough, please try again")
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
