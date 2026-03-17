package order

import "time"

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
	SubTotalPrice        float64        `json:"sub_total_price"`
	DiscountPrice        float64        `json:"discount_price"`
	TotalPrice           float64        `json:"total_price"`
	BurnPoint            int            `json:"burn_point"`
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

type OrderProductWithPrice struct {
	ProductBrand string  `json:"product_brand" db:"product_brand"`
	ProductName  string  `json:"product_name" db:"product_name"`
	Quantity     int     `json:"quantity" db:"quantity"`
	Price        float64 `json:"price" db:"product_price"`
}

type OrderSummaryProduct struct {
	ProductBrand  string  `json:"product_brand"`
	ProductName   string  `json:"product_name"`
	Quantity      int     `json:"quantity"`
	PriceTHB      float64 `json:"price_thb"`
	TotalPriceTHB float64 `json:"total_price_thb"`
}

type Order struct {
	OrderNumber int64
}

type OrderDetail struct {
	ID               int     `json:"id"  db:"id"`
	OrderNumber      int64   `json:"order_number" db:"order_number"`
	UserID           int     `json:"user_id"  db:"user_id"`
	ShippingMethodID int     `json:"shipping_method_id"  db:"shipping_method_id"`
	PaymentMethodID  int     `json:"payment_method_id"  db:"payment_method_id"`
	SubTotalPrice    float64 `json:"sub_total_price" db:"sub_total_price"`
	DiscountPrice    float64 `json:"discount_price" db:"discount_price"`
	TotalPrice       float64 `json:"total_price" db:"total_price"`
	ShippingFee      float64 `json:"shipping_fee" db:"shipping_fee"`
	BurnPoint        int     `json:"burn_point" db:"burn_point"`
	EarnPoint        int     `json:"earn_point" db:"earn_point"`
	TransactionID    string  `json:"transaction_id" db:"transaction_id"`
	Status           string  `json:"status" db:"status"`
}

type OrderDetailWithTrackingNumber struct {
	ID               int       `json:"id"  db:"id"`
	OrderNumber      int64     `json:"order_number" db:"order_number"`
	UserID           int       `json:"user_id"  db:"user_id"`
	ShippingMethodID int       `json:"shipping_method_id"  db:"shipping_method_id"`
	PaymentMethodID  int       `json:"payment_method_id"  db:"payment_method_id"`
	SubTotalPrice    float64   `json:"sub_total_price" db:"sub_total_price"`
	DiscountPrice    float64   `json:"discount_price" db:"discount_price"`
	TotalPrice       float64   `json:"total_price" db:"total_price"`
	ShippingFee      float64   `json:"shipping_fee" db:"shipping_fee"`
	BurnPoint        int       `json:"burn_point" db:"burn_point"`
	EarnPoint        int       `json:"earn_point" db:"earn_point"`
	TransactionID    string    `json:"transaction_id" db:"transaction_id"`
	Status           string    `json:"status" db:"status"`
	TrackingNumber   string    `json:"tracking_no" db:"tracking_no"`
	Updated          time.Time `json:"updated" db:"updated"`
}

type OrderSummary struct {
	OrderNumber      int64                 `json:"order_number"`
	FirstName        string                `json:"first_name"`
	LastName         string                `json:"last_name"`
	TrackingNumber   string                `json:"tracking_no"`
	ShippingMethod   string                `json:"shipping_method"`
	PaymentMethod    string                `json:"payment_method"`
	OrderProductList []OrderSummaryProduct `json:"products"`
	SubTotalPrice    float64               `json:"subtotal_price"`
	DiscountPrice    float64               `json:"discount_price"`
	TotalPrice       float64               `json:"total_price"`
	ShippingFee      float64               `json:"shipping_fee"`
	BurnPoint        int                   `json:"burn_point"`
	ReceivingPoint   int                   `json:"receiving_point"`
	IssuedDate       string                `json:"issued_date"`
}
