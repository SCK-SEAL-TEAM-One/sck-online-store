package payment_test

import (
	"context"
	"store-service/internal/order"
	"store-service/internal/payment"
	"store-service/internal/product"
	"store-service/internal/shipping"

	"github.com/stretchr/testify/mock"
)

type mockBankGateway struct {
	mock.Mock
}

func (gateway *mockBankGateway) Payment(ctx context.Context, paymentDetail payment.PaymentDetail) (string, error) {
	argument := gateway.Called(ctx, paymentDetail)
	return argument.String(0), argument.Error(1)
}

func (gateway *mockBankGateway) GetCardDetail(ctx context.Context, orgID int, userID int) (payment.CardDetail, error) {
	argument := gateway.Called(ctx, orgID, userID)
	return argument.Get(0).(payment.CardDetail), argument.Error(1)
}

type mockShippingGateway struct {
	mock.Mock
}

func (gateway *mockShippingGateway) GetTrackingNumber(ctx context.Context, shippingGatewaySubmit shipping.ShippingGatewaySubmit) (string, error) {
	argument := gateway.Called(ctx, shippingGatewaySubmit)
	return argument.String(0), argument.Error(1)
}

type mockOrderRepository struct {
	mock.Mock
}

func (repo *mockOrderRepository) CreateOrder(ctx context.Context, userID int, orderDetail order.OrderDetail) (int, error) {
	argument := repo.Called(ctx, userID, orderDetail)
	return argument.Int(0), argument.Error(1)
}

func (repo *mockOrderRepository) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (order.OrderDetail, error) {
	argument := repo.Called(ctx, orderNumber)
	return argument.Get(0).(order.OrderDetail), argument.Error(1)
}

func (repo *mockOrderRepository) GetLastOrderNumber(ctx context.Context, yearPrefix string) (string, error) {
	argument := repo.Called(ctx, yearPrefix)
	return argument.Get(0).(string), argument.Error(1)
}

func (repo *mockOrderRepository) GetOrderWithTrackingNumberByOrderNumber(ctx context.Context, orderNumber string) (order.OrderDetailWithTrackingNumber, error) {
	argument := repo.Called(ctx, orderNumber)
	return argument.Get(0).(order.OrderDetailWithTrackingNumber), argument.Error(1)
}

func (repo *mockOrderRepository) CreateOrderProduct(ctx context.Context, orderID, productID, quantity int, productPrice float64) error {
	argument := repo.Called(ctx, orderID, productID, quantity, productPrice)
	return argument.Error(0)
}

func (repo *mockOrderRepository) GetOrderProduct(ctx context.Context, orderID int) ([]order.OrderProduct, error) {
	argument := repo.Called(ctx, orderID)
	return argument.Get(0).([]order.OrderProduct), argument.Error(1)
}

func (repo *mockOrderRepository) GetOrderProductWithPrice(ctx context.Context, orderID int) ([]order.OrderProductWithPrice, error) {
	argument := repo.Called(ctx, orderID)
	return argument.Get(0).([]order.OrderProductWithPrice), argument.Error(1)
}

func (repo *mockOrderRepository) CreateShipping(ctx context.Context, userID int, orderID int, shippingInfo order.ShippingInfo) (int, error) {
	argument := repo.Called(ctx, userID, orderID, shippingInfo)
	return argument.Int(0), argument.Error(1)
}

func (repo *mockOrderRepository) UpdateOrderTransaction(ctx context.Context, orderID int, transactionID string) error {
	argument := repo.Called(ctx, orderID, transactionID)
	return argument.Error(0)
}

func (repo *mockOrderRepository) UpdateOrderTrackingNumber(ctx context.Context, orderID int, trackingNumber string) error {
	argument := repo.Called(ctx, orderID, trackingNumber)
	return argument.Error(0)
}

type mockProductRepository struct {
	mock.Mock
}

func (repo *mockProductRepository) GetProducts(ctx context.Context, keyword string, limit string, offset string) (product.ProductResult, error) {
	argument := repo.Called(ctx, keyword)
	return argument.Get(0).(product.ProductResult), argument.Error(1)
}

func (repository *mockProductRepository) GetProductByID(ctx context.Context, id int) (product.ProductDetail, error) {
	argument := repository.Called(ctx, id)
	return argument.Get(0).(product.ProductDetail), argument.Error(1)
}

func (repository *mockProductRepository) UpdateStock(ctx context.Context, productId int, quantity int) error {
	argument := repository.Called(ctx, productId, quantity)
	return argument.Error(0)
}
