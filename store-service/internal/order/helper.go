package order

import (
	"fmt"
	"strconv"
	"time"
)

type OrderHelperInterface interface {
	GenerateOrderNumber(paymentMethodID, shippingMethodID, seq int, now time.Time) (string, error)
	GetNextSequence(yearPrefix string) (int, error)
}

type OrderHelper struct{}

func (helper OrderHelper) GenerateOrderNumber(paymentMethodID, shippingMethodID, seq int, now time.Time) (string, error) {
	if seq > 999 {
		return "", fmt.Errorf("Invalid sequence: %d (limit is 999)", seq)
	}
	if seq < 1 {
		return "", fmt.Errorf("Invalid sequence: %d (must be positive)", seq)
	}

	var PaymentMethod = map[int]string{
		1: "95",
		2: "98",
	}
	var ShippingMethod = map[int]string{
		1: "22",
		2: "33",
		3: "44",
	}
	paymentMethod := PaymentMethod[paymentMethodID]
	shippingMethod := ShippingMethod[shippingMethodID]

	// Format: YYMMDD (Year-Month-Day)
	dateStr := now.Format("060102")

	orderNumber := fmt.Sprintf("%s%s%s%03d", dateStr, paymentMethod, shippingMethod, seq)
	return orderNumber, nil
}

func (helper OrderHelper) GetNextSequence(lastOrderNumber string) (int, error) {
	const (
		orderNumberLength = 13
		sequenceLength    = 3
		maxSequence       = 999
	)

	// Format: YYMMDDPPSMSEQ **** 13 Digits ****
	if len(lastOrderNumber) != orderNumberLength {
		return 0, fmt.Errorf("Invalid length for order %s", lastOrderNumber)
	}

	seqStr := lastOrderNumber[len(lastOrderNumber)-sequenceLength:]
	currentSeq, err := strconv.Atoi(seqStr)
	if err != nil || currentSeq == 0 {
		return 0, fmt.Errorf("Invalid sequence number: %s", seqStr)
	}

	nextSeq := currentSeq + 1
	if nextSeq > 999 {
		return 0, fmt.Errorf("Annual order limit (%d) reached", maxSequence)
	}

	return nextSeq, nil
}
