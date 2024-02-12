//go:build integration
// +build integration

package shipping_test

import (
	"store-service/internal/shipping"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetTrackingNumber_Input_ShippingMethodID_Should_Be_Tracking_Number_No_Error(t *testing.T) {

	service := shipping.ShippingGateway{
		ShippingEndpoint: "http://localhost:8883",
	}
	actualTrackingNumber, err := service.GetTrackingNumber(shipping.ShippingGatewaySubmit{
		ShippingMethodID: 1,
	})

	assert.NotEmpty(t, actualTrackingNumber)
	assert.Equal(t, nil, err)
}
