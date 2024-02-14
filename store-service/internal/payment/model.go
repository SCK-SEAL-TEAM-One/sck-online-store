package payment

import "time"

type SubmitedPayment struct {
	OrderID int    `json:"order_id"`
	OTP     int    `json:"otp"`
	RefOTP  string `json:"ref_otp"`
}

type SubmitedPaymentResponse struct {
	OrderID          int       `json:"order_id"`
	PaymentDate      time.Time `json:"payment_date"`
	ShippingMethodID int       `json:"shipping_method_id"`
	TrackingNumber   string    `json:"tracking_number"`
}

type PaymentDetail struct {
	CardNumber   string  `json:"card_number"`
	CVV          int     `json:"cvv"`
	ExpiredMonth int     `json:"expired_month"`
	ExpiredYear  int     `json:"expired_year"`
	CardName     string  `json:"card_name"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	MerchantID   int     `json:"merchant_id"`
}

type CardDetail struct {
	CardNumber   string `json:"card_number"`
	CVV          int    `json:"cvv"`
	ExpiredMonth int    `json:"expired_month"`
	ExpiredYear  int    `json:"expired_year"`
	CardName     string `json:"card_name"`
}
