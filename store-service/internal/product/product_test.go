package product_test

import (
	"errors"
	"testing"

	"store-service/internal/product"

	"github.com/stretchr/testify/assert"
)

func Test_GetProducts_Should_be_Return_Total_10031_and_Products_include_PriceTHB(t *testing.T) {
	expected := product.ProductResult{
		Total: 10031,
		Products: []product.Product{
			{
				ID:           1,
				Name:         "Balance Training Bicycle",
				Price:        119.95,
				PriceTHB:     4314.6,
				PriceFullTHB: 4314.597182,
				Image:        "/Balance_Training_Bicycle.png",
			},
		},
	}
	keyword := ""
	limit := ""
	offset := ""

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("GetProducts", keyword, limit, offset).Return(product.ProductResult{
		Total: 10031,
		Products: []product.Product{
			{
				ID:           1,
				Name:         "Balance Training Bicycle",
				Price:        119.95,
				PriceTHB:     0,
				PriceFullTHB: 0,
				Image:        "/Balance_Training_Bicycle.png",
			},
		},
	}, nil)

	productService := product.ProductService{
		ProductRepository: mockProductRepository,
	}
	actual, err := productService.GetProducts(keyword, limit, offset)

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_GetProducts_Should_be_Return_GetProducts_Error(t *testing.T) {
	expected := product.ProductResult{}
	keyword := ""
	limit := ""
	offset := ""

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("GetProducts", keyword, limit, offset).Return(product.ProductResult{}, errors.New("GetProducts Error"))

	productService := product.ProductService{
		ProductRepository: mockProductRepository,
	}
	actual, err := productService.GetProducts(keyword, limit, offset)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_GetProductByID_Should_be_Return_ProductDetail_ID_1_include_PriceTHB(t *testing.T) {
	expected := product.ProductDetail{
		ID:           1,
		Name:         "Balance Training Bicycle",
		Price:        119.95,
		PriceTHB:     4314.6,
		PriceFullTHB: 4314.597182,
		Image:        "/Balance_Training_Bicycle.png",
		Stock:        100,
		Brand:        "SportsFun",
	}
	pid := 1

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("GetProductByID", pid).Return(product.ProductDetail{
		ID:           1,
		Name:         "Balance Training Bicycle",
		Price:        119.95,
		PriceTHB:     0,
		PriceFullTHB: 0,
		Image:        "/Balance_Training_Bicycle.png",
		Stock:        100,
		Brand:        "SportsFun",
	}, nil)

	productService := product.ProductService{
		ProductRepository: mockProductRepository,
	}
	actual, err := productService.GetProductByID(pid)

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_GetProductByID_Should_be_Return_GetProductByID_Error(t *testing.T) {
	expected := product.ProductDetail{}
	pid := 1

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("GetProductByID", pid).Return(product.ProductDetail{}, errors.New("GetProductByID Error"))

	productService := product.ProductService{
		ProductRepository: mockProductRepository,
	}
	actual, err := productService.GetProductByID(pid)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}
