package shipping

type ShippingGatewaySubmit struct {
	ShippingMethodID int `json:"shipping_method_id"`
}

type ShippingMethodDetail struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Fee         float64 `json:"fee" db:"fee"`
}
