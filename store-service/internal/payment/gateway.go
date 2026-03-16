package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BankGateway struct {
	BankEndpoint string
}

type BankGatewayResponse struct {
	Status        string `json:"status"`
	PaymentDate   string `json:"payment_date"`
	TransactionID string `json:"transaction_id"`
}

func (gateway BankGateway) Payment(ctx context.Context, paymentDetail PaymentDetail) (string, error) {
	data, _ := json.Marshal(paymentDetail)
	endPoint := gateway.BankEndpoint + "/payment/visa"
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

	var BankGatewayResponse BankGatewayResponse
	err = json.Unmarshal(responseData, &BankGatewayResponse)
	if err != nil {
		return "", err
	}

	return BankGatewayResponse.TransactionID, nil
}

func (gateway BankGateway) GetCardDetail(ctx context.Context, orgID int, userID int) (CardDetail, error) {
	endPoint := gateway.BankEndpoint + fmt.Sprintf("/card/information?oid=%d&uid=%d", orgID, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endPoint, nil)
	if err != nil {
		return CardDetail{}, err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return CardDetail{}, err
	}
	if response.StatusCode != 200 {
		return CardDetail{}, fmt.Errorf("response is not ok but it's %d", response.StatusCode)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return CardDetail{}, err
	}

	var CardDetailResponse CardDetail
	err = json.Unmarshal(responseData, &CardDetailResponse)
	if err != nil {
		return CardDetail{}, err
	}

	return CardDetailResponse, nil
}
