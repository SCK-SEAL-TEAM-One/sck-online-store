package product

type ProductResult struct {
	Total    int       `json:"total"`
	Products []Product `json:"products"`
}

type Product struct {
	ID           int     `json:"id" db:"id"`
	Name         string  `json:"product_name" db:"product_name"`
	Price        float64 `json:"product_price" db:"product_price"`
	PriceTHB     float64 `json:"product_price_thb"`
	PriceFullTHB float64 `json:"product_price_full_thb"`
	Image        string  `json:"product_image" db:"image_url"`
}

type ProductDetail struct {
	ID           int     `json:"id" db:"id"`
	Name         string  `json:"product_name" db:"product_name"`
	Price        float64 `json:"product_price" db:"product_price"`
	PriceTHB     float64 `json:"product_price_thb"`
	PriceFullTHB float64 `json:"product_price_full_thb"`
	Image        string  `json:"product_image" db:"image_url"`
	Stock        int     `json:"stock" db:"stock"`
	Brand        string  `json:"product_brand" db:"product_brand"`
}
