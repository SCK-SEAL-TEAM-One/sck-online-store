package point

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PointGateway struct {
	PointEndpoint string
}

func (gateway PointGateway) GetPoints(uid int) ([]Point, error) {
	endPoint := gateway.PointEndpoint + "/api/v1/point"
	response, err := http.Get(endPoint)
	if err != nil {
		return []Point{}, err
	}
	if response.StatusCode != 200 {
		return []Point{}, fmt.Errorf("response is not ok but it's %d", response.StatusCode)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []Point{}, err
	}

	var PointGatewayResponse []Point
	err = json.Unmarshal(responseData, &PointGatewayResponse)
	if err != nil {
		return []Point{}, err
	}

	return PointGatewayResponse, nil
}

func (gateway PointGateway) CreatePoint(uid int, body Point) (Point, error) {
	data, _ := json.Marshal(body)
	endPoint := gateway.PointEndpoint + "/api/v1/point"
	response, err := http.Post(endPoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return Point{}, err
	}
	if response.StatusCode != 200 && response.StatusCode != 201 {
		return Point{}, fmt.Errorf("response is not ok but it's %d", response.StatusCode)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Point{}, err
	}

	var PointGatewayResponse Point
	err = json.Unmarshal(responseData, &PointGatewayResponse)
	if err != nil {
		return Point{}, err
	}

	return PointGatewayResponse, nil
}
