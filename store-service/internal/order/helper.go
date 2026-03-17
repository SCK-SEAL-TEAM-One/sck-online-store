package order

import (
	"fmt"
	"time"
)

type OrderHelperInterface interface {
	GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq int, now time.Time) (int64, error)
}

type OrderHelper struct{}

func (helper OrderHelper) GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq int, now time.Time) (int64, error) {
	if userID < 1 || userID > 999 {
		return 0, fmt.Errorf("Invalid userID: %d (must be 1-999)", userID)
	}
	if seq > 999 {
		return 0, fmt.Errorf("Invalid sequence: %d (limit is 999)", seq)
	}
	if seq < 1 {
		return 0, fmt.Errorf("Invalid sequence: %d (must be positive)", seq)
	}

	var PaymentMethod = map[int]int64{
		1: 95,
		2: 98,
	}
	var ShippingMethod = map[int]int64{
		1: 22,
		2: 33,
		3: 44,
	}
	paymentMethod := PaymentMethod[paymentMethodID]
	shippingMethod := ShippingMethod[shippingMethodID]

	yy := int64(now.Year() % 100)
	mm := int64(now.Month())
	dd := int64(now.Day())
	datePart := yy*10000 + mm*100 + dd

	// Format: YYMMDD PP SM UUU SEQ (16 digits)
	orderNumber := datePart*1e10 + paymentMethod*1e8 + shippingMethod*1e6 + int64(userID)*1e3 + int64(seq)

	return orderNumber, nil
}
