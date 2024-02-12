package payment

import (
	"bytes"
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

func (gateway BankGateway) Payment(paymentDetail PaymentDetail) (string, error) {
	data, _ := json.Marshal(paymentDetail)
	endPoint := gateway.BankEndpoint + "/payment/visa"
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

	var BankGatewayResponse BankGatewayResponse
	err = json.Unmarshal(responseData, &BankGatewayResponse)
	if err != nil {
		return "", err
	}

	return BankGatewayResponse.TransactionID, nil
}

func (gateway BankGateway) GetCardDetail(orgID int, userID int) (CardDetail, error) {
	endPoint := gateway.BankEndpoint + fmt.Sprintf("/card/information?oid=%d&uid=%d", orgID, userID)
	response, err := http.Get(endPoint)
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
