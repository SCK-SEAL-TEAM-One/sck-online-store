package product

import (
	"log"
	"store-service/internal/common"
)

type ProductInterface interface {
	GetProducts(keyword string, limit string, offset string) (ProductResult, error)
	GetProductByID(ID int) (ProductDetail, error)
}

type ProductService struct {
	ProductRepository ProductRepository
}

func (productService ProductService) GetProducts(keyword string, limit string, offset string) (ProductResult, error) {
	res, err := productService.ProductRepository.GetProducts(keyword, limit, offset)
	if err != nil {
		log.Printf("ProductRepository.GetProducts internal error %s", err.Error())
		return ProductResult{}, err
	}

	for i := range res.Products {
		p := &res.Products[i]
		digit := common.ConvertToThb(p.Price)

		p.PriceTHB = digit.ShortDecimal
		p.PriceFullTHB = digit.LongDecimal
	}
	return res, err
}

func (productService ProductService) GetProductByID(ID int) (ProductDetail, error) {
	productDetail, err := productService.ProductRepository.GetProductByID(ID)
	if err != nil {
		log.Printf("ProductRepository.GetProductByID internal error %s", err.Error())
		return ProductDetail{}, err
	}

	p := &productDetail
	digit := common.ConvertToThb(p.Price)

	p.PriceTHB = digit.ShortDecimal
	p.PriceFullTHB = digit.LongDecimal

	return productDetail, err
}
