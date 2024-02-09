package order

type SubmitedOrder struct {
	Cart                 []OrderProduct `json:"cart"`
	ShippingMethodID     int            `json:"shipping_method_id"`
	ShippingAddress      string         `json:"shipping_address"`
	ShippingSubDistrict  string         `json:"shipping_sub_disterict"`
	ShippingDistrict     string         `json:"shipping_district"`
	ShippingProvince     string         `json:"shipping_province"`
	ShippingZipCode      string         `json:"shipping_zip_code"`
	RecipientFirstName   string         `json:"recipient_first_name"`
	RecipientLastName    string         `json:"recipient_last_name"`
	RecipientPhoneNumber string         `json:"recipient_phone_number"`
	PaymentMethodID      int            `json:"payment_method_id"`
	BurnPoint            int            `json:"burn_point"`
	SubTotalPrice        float64        `json:"sub_total_price"`
	DiscountPrice        float64        `json:"discount_price"`
	TotalPrice           float64        `json:"total_price"`
}

type ShippingInfo struct {
	UserID               int    `db:"user_id"`
	ShippingMethodID     int    `db:"method_id"`
	ShippingAddress      string `db:"address"`
	ShippingSubDistrict  string `db:"sub_district"`
	ShippingDistrict     string `db:"district"`
	ShippingProvince     string `db:"province"`
	ShippingZipCode      string `db:"zip_code"`
	RecipientFirstName   string `db:"recipient_first_name"`
	RecipientLastName    string `db:"recipient_last_name"`
	RecipientPhoneNumber string `db:"phone_number"`
}

type OrderProduct struct {
	ProductID int `json:"product_id" db:"product_id"`
	Quantity  int `json:"quantity" db:"quantity"`
}

type Order struct {
	OrderID int
}

// func (s SubmitedOrder) GetShippingFee() float64 {
// 	return 2.00
// }

// func (s SubmitedOrder) GetShippingMethodProvider() string {
// 	return "Kerry"
// }
