package product_test

import (
	"context"
	"store-service/internal/product"

	"github.com/stretchr/testify/mock"
)

type mockProductRepository struct {
	mock.Mock
}

func (repo *mockProductRepository) GetProducts(ctx context.Context, keyword string, limit string, offset string) (product.ProductResult, error) {
	argument := repo.Called(ctx, keyword, limit, offset)
	return argument.Get(0).(product.ProductResult), argument.Error(1)
}

func (repository *mockProductRepository) GetProductByID(ctx context.Context, id int) (product.ProductDetail, error) {
	argument := repository.Called(ctx, id)
	return argument.Get(0).(product.ProductDetail), argument.Error(1)
}

func (repository *mockProductRepository) UpdateStock(ctx context.Context, productID int, stock int) error {
	argument := repository.Called(ctx, productID, stock)
	return argument.Error(0)
}
