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

type PointGatewayResponseItem struct {
	OrgId  int `json:"orgId"`
	UserId int `json:"userId"`
	Amount int `json:"amount"`
}

func (gateway PointGateway) GetPoints(uid int) ([]PointGatewayResponseItem, error) {
	endPoint := gateway.PointEndpoint + "/api/v1/point"
	response, err := http.Get(endPoint)
	if err != nil {
		return []PointGatewayResponseItem{}, err
	}
	if response.StatusCode != 200 {
		return []PointGatewayResponseItem{}, fmt.Errorf("response is not ok but it's %d", response.StatusCode)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []PointGatewayResponseItem{}, err
	}

	var PointGatewayResponse []PointGatewayResponseItem
	err = json.Unmarshal(responseData, &PointGatewayResponse)
	if err != nil {
		return []PointGatewayResponseItem{}, err
	}

	return PointGatewayResponse, nil
}

func (gateway PointGateway) CreatePoint(uid int, body Point) (PointGatewayResponseItem, error) {
	data, _ := json.Marshal(body)
	endPoint := gateway.PointEndpoint + "/api/v1/point"
	response, err := http.Post(endPoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return PointGatewayResponseItem{}, err
	}
	if response.StatusCode != 200 && response.StatusCode != 201 {
		return PointGatewayResponseItem{}, fmt.Errorf("response is not ok but it's %d", response.StatusCode)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return PointGatewayResponseItem{}, err
	}

	var PointGatewayResponse PointGatewayResponseItem
	err = json.Unmarshal(responseData, &PointGatewayResponse)
	if err != nil {
		return PointGatewayResponseItem{}, err
	}

	return PointGatewayResponse, nil
}
