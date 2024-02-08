package point_test

import (
	"store-service/internal/point"

	"github.com/stretchr/testify/mock"
)

type mockPointGateway struct {
	mock.Mock
}

func (gateway *mockPointGateway) GetPoints(userID int) ([]point.Point, error) {
	argument := gateway.Called(userID)
	return argument.Get(0).([]point.Point), argument.Error(1)
}

func (gateway *mockPointGateway) CreatePoint(userID int, pointItem point.Point) (point.Point, error) {
	argument := gateway.Called(userID, pointItem)
	return argument.Get(0).(point.Point), argument.Error(1)
}
