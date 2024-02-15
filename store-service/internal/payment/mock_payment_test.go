package payment_test

import (
	"store-service/internal/order"
	"store-service/internal/payment"
	"store-service/internal/product"
	"store-service/internal/shipping"

	"github.com/stretchr/testify/mock"
)

type mockBankGateway struct {
	mock.Mock
}

func (gateway *mockBankGateway) Payment(paymentDetail payment.PaymentDetail) (string, error) {
	argument := gateway.Called(paymentDetail)
	return argument.String(0), argument.Error(1)
}

func (gateway *mockBankGateway) GetCardDetail(orgID int, userID int) (payment.CardDetail, error) {
	argument := gateway.Called(orgID, userID)
	return argument.Get(0).(payment.CardDetail), argument.Error(1)
}

type mockShippingGateway struct {
	mock.Mock
}

func (gateway *mockShippingGateway) GetTrackingNumber(shippingGatewaySubmit shipping.ShippingGatewaySubmit) (string, error) {
	argument := gateway.Called(shippingGatewaySubmit)
	return argument.String(0), argument.Error(1)
}

type mockOrderRepository struct {
	mock.Mock
}

func (repo *mockOrderRepository) CreateOrder(userID int, orderDetail order.OrderDetail) (int, error) {
	argument := repo.Called(userID, orderDetail)
	return argument.Int(0), argument.Error(1)
}

func (repo *mockOrderRepository) GetOrderByID(ID int) (order.OrderDetail, error) {
	argument := repo.Called(ID)
	return argument.Get(0).(order.OrderDetail), argument.Error(1)
}

func (repo *mockOrderRepository) CreateOrderProduct(orderID, productID, quantity int, productPrice float64) error {
	argument := repo.Called(orderID, productID, quantity, productPrice)
	return argument.Error(0)
}

func (repo *mockOrderRepository) GetOrderProduct(orderID int) ([]order.OrderProduct, error) {
	argument := repo.Called(orderID)
	return argument.Get(0).([]order.OrderProduct), argument.Error(1)
}

func (repo *mockOrderRepository) CreateShipping(userID int, orderID int, shippingInfo order.ShippingInfo) (int, error) {
	argument := repo.Called(userID, orderID, shippingInfo)
	return argument.Int(0), argument.Error(1)
}

func (repo *mockOrderRepository) UpdateOrderTransaction(orderID int, transactionID string) error {
	argument := repo.Called(orderID, transactionID)
	return argument.Error(0)
}

type mockProductRepository struct {
	mock.Mock
}

func (repo *mockProductRepository) GetProducts(keyword string, limit string, offset string) (product.ProductResult, error) {
	argument := repo.Called(keyword)
	return argument.Get(0).(product.ProductResult), argument.Error(1)
}

func (repository *mockProductRepository) GetProductByID(id int) (product.ProductDetail, error) {
	argument := repository.Called(id)
	return argument.Get(0).(product.ProductDetail), argument.Error(1)
}

func (repository *mockProductRepository) UpdateStock(productId int, quantity int) error {
	argument := repository.Called(productId, quantity)
	return argument.Error(0)
}
