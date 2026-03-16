package shipping

import (
	"bytes"
	"context"
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

func (gateway ShippingGateway) GetTrackingNumber(ctx context.Context, shippingGatewaySubmit ShippingGatewaySubmit) (string, error) {
	data, _ := json.Marshal(shippingGatewaySubmit)
	endPoint := gateway.ShippingEndpoint + "/shipping"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endPoint, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
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
