package shipping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ShippingGateway struct {
	ShippingEndpoint string
}

type ShippingGatewayResponse struct {
	TrackingNumber string `json:"tracking_number"`
}

func (gateway ShippingGateway) GetTrackingNumber(shippingGatewaySubmit ShippingGatewaySubmit) (string, error) {
	data, _ := json.Marshal(shippingGatewaySubmit)
	endPoint := gateway.ShippingEndpoint + "/shipping"
	response, err := http.Post(endPoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("response is not ok but it's %d", response.StatusCode)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var ShippingGatewayResponse ShippingGatewayResponse
	err = json.Unmarshal(responseData, &ShippingGatewayResponse)
	if err != nil {
		return "", err
	}

	return ShippingGatewayResponse.TrackingNumber, nil
}
