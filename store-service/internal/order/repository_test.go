//go:build integration
// +build integration

package order_test

import (
	"fmt"
	"store-service/internal/order"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_OrderRepository(t *testing.T) {
	connection, err := sqlx.Connect("mysql", "user:password@(localhost:3305)/store")
	if err != nil {
		t.Fatalf("cannot tearup data err %s", err)
	}
	repository := order.OrderRepositoryMySQL{
		DBConnection: connection,
	}

	t.Run("CreateOrder_Input_SubmitedOrder_Should_Be_OrderID_No_Error", func(t *testing.T) {
		uid := 1
		submittedOrder := order.SubmitedOrder{
			Cart: []order.OrderProduct{
				{
					ProductID: 2,
					Quantity:  1,
				},
			},
			ShippingMethodID:     1,
			ShippingAddress:      "405/37 ถ.มหิดล",
			ShippingSubDistrict:  "ท่าศาลา",
			ShippingDistrict:     "เมือง",
			ShippingProvince:     "เชียงใหม่",
			ShippingZipCode:      "50000",
			RecipientFirstName:   "ณัฐญา",
			RecipientLastName:    "ชุติบุตร",
			RecipientPhoneNumber: "0970809292",
			PaymentMethodID:      1,
			BurnPoint:            0,
			SubTotalPrice:        100.00,
			DiscountPrice:        10.00,
			TotalPrice:           90.00,
		}

		actualId, err := repository.CreateOrder(uid, submittedOrder)

		assert.Equal(t, nil, err)
		assert.NotEmpty(t, actualId)
	})

	t.Run("CreateOrderProduct_Input_OrderID_2_And_ProductID_2_Should_Be_No_Error", func(t *testing.T) {
		oid := 1
		pid := 2
		qty := 1
		productPrice := 12.95
		err := repository.CreateOrderProduct(oid, pid, qty, productPrice)

		assert.Equal(t, nil, err)
	})

	t.Run("UpdateOrder_Input_TransactionID_TXN202002021525_OrderID_1_Should_No_Error", func(t *testing.T) {
		txn := "TXN202002021525"
		oid := 1

		err := repository.UpdateOrder(oid, txn)

		assert.Equal(t, nil, err)
	})

	t.Run("UpdateOrder_Input_TransactionID_TXN202002021525_OrderID_11111111119_Should_Get_Error_No_Row_Affected", func(t *testing.T) {
		expectedError := fmt.Errorf("no any row affected , update not completed")
		transactionID := "TXN202002021525"
		orderID := 11111111119

		err := repository.UpdateOrder(orderID, transactionID)

		assert.Equal(t, expectedError, err)
	})

	t.Run("GetOrderProduct_Input_OrderID_2_Should_Be_OrderProducts", func(t *testing.T) {
		expected := []order.OrderProduct{
			{
				ProductID: 2,
				Quantity:  1,
			},
		}

		oid := 1

		actual, err := repository.GetOrderProduct(oid)

		assert.Equal(t, expected, actual)
		assert.Equal(t, nil, err)
	})

	t.Run("CreateShipping_Input_OrderID_1_and_ShippingInfo_Should_Be_ShippingID_No_Error", func(t *testing.T) {
		uid := 1
		oid := 1
		shippingInfo := order.ShippingInfo{
			ShippingMethodID:     1,
			ShippingAddress:      "405/37 ถ.มหิดล",
			ShippingSubDistrict:  "ท่าศาลา",
			ShippingDistrict:     "เมือง",
			ShippingProvince:     "เชียงใหม่",
			ShippingZipCode:      "50000",
			RecipientFirstName:   "ณัฐญา",
			RecipientLastName:    "ชุติบุตร",
			RecipientPhoneNumber: "0970809292",
		}

		actualId, err := repository.CreateShipping(uid, oid, shippingInfo)

		assert.Equal(t, nil, err)
		assert.NotEmpty(t, actualId)
	})
}
