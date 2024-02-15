package product_test

import (
	"store-service/internal/product"

	"github.com/stretchr/testify/mock"
)

type mockProductRepository struct {
	mock.Mock
}

func (repo *mockProductRepository) GetProducts(keyword string, limit string, offset string) (product.ProductResult, error) {
	argument := repo.Called(keyword, limit, offset)
	return argument.Get(0).(product.ProductResult), argument.Error(1)
}

func (repository *mockProductRepository) GetProductByID(id int) (product.ProductDetail, error) {
	argument := repository.Called(id)
	return argument.Get(0).(product.ProductDetail), argument.Error(1)
}

func (repository *mockProductRepository) UpdateStock(productID int, stock int) error {
	argument := repository.Called(productID, stock)
	return argument.Error(0)
}
