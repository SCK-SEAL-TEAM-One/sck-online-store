package order_test

import (
	"store-service/internal/auth"
	"store-service/internal/cart"
	"store-service/internal/order"
	"store-service/internal/point"
	"store-service/internal/product"
	"store-service/internal/shipping"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockPointInterface struct {
	mock.Mock
}

func (service *mockPointInterface) TotalPoint(uid int) (point.TotalPoint, error) {
	argument := service.Called(uid)
	return argument.Get(0).(point.TotalPoint), argument.Error(1)
}

func (service *mockPointInterface) DeductPoint(uid int, submitedPoint point.SubmitedPoint) (point.TotalPoint, error) {
	argument := service.Called(uid, submitedPoint)
	return argument.Get(0).(point.TotalPoint), argument.Error(1)
}

func (service *mockPointInterface) CheckBurnPoint(uid int, amount int) (bool, error) {
	argument := service.Called(uid, amount)
	return argument.Bool(0), argument.Error(1)
}

type mockOrderHelper struct {
	mock.Mock
}

func (m *mockOrderHelper) GenerateOrderNumber(paymentMethodID, shippingMethodID, seq int, now time.Time) (string, error) {
	args := m.Called(paymentMethodID, shippingMethodID, seq, now)
	return args.String(0), args.Error(1)
}

func (m *mockOrderHelper) GetNextSequence(yearPrefix string) (int, error) {
	args := m.Called(yearPrefix)
	return args.Int(0), args.Error(1)
}

type mockOrderRepository struct {
	mock.Mock
}

func (repo *mockOrderRepository) CreateOrder(userID int, orderDetail order.OrderDetail) (int, error) {
	argument := repo.Called(userID, orderDetail)
	return argument.Int(0), argument.Error(1)
}

func (repo *mockOrderRepository) GetOrderByOrderNumber(orderNumber string) (order.OrderDetail, error) {
	argument := repo.Called(orderNumber)
	return argument.Get(0).(order.OrderDetail), argument.Error(1)
}

func (repo *mockOrderRepository) GetLastOrderNumber(yearPrefix string) (string, error) {
	argument := repo.Called(yearPrefix)
	return argument.Get(0).(string), argument.Error(1)
}

func (repo *mockOrderRepository) GetOrderWithTrackingNumberByOrderNumber(orderNumber string) (order.OrderDetailWithTrackingNumber, error) {
	argument := repo.Called(orderNumber)
	return argument.Get(0).(order.OrderDetailWithTrackingNumber), argument.Error(1)
}

func (repo *mockOrderRepository) CreateOrderProduct(orderID, productID, quantity int, productPrice float64) error {
	argument := repo.Called(orderID, productID, quantity, productPrice)
	return argument.Error(0)
}

func (repo *mockOrderRepository) GetOrderProduct(orderID int) ([]order.OrderProduct, error) {
	argument := repo.Called(orderID)
	return argument.Get(0).([]order.OrderProduct), argument.Error(1)
}

func (repo *mockOrderRepository) GetOrderProductWithPrice(orderID int) ([]order.OrderProductWithPrice, error) {
	argument := repo.Called(orderID)
	return argument.Get(0).([]order.OrderProductWithPrice), argument.Error(1)
}

func (repo *mockOrderRepository) CreateShipping(userID int, orderID int, shippingInfo order.ShippingInfo) (int, error) {
	argument := repo.Called(userID, orderID, shippingInfo)
	return argument.Int(0), argument.Error(1)
}

func (repo *mockOrderRepository) UpdateOrderTransaction(orderID int, transactionID string) error {
	argument := repo.Called(orderID, transactionID)
	return argument.Error(1)
}

func (repo *mockOrderRepository) UpdateOrderTrackingNumber(orderID int, trackingNumber string) error {
	argument := repo.Called(orderID, trackingNumber)
	return argument.Error(1)
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

type mockCartRepository struct {
	mock.Mock
}

func (repo *mockCartRepository) GetCartDetail(userID int) ([]cart.CartDetail, error) {
	argument := repo.Called(userID)
	return argument.Get(0).([]cart.CartDetail), argument.Error(1)
}

func (repo *mockCartRepository) GetCartByProductID(userID int, productID int) (cart.Cart, error) {
	argument := repo.Called(userID, productID)
	return argument.Get(0).(cart.Cart), argument.Error(1)
}

func (repo *mockCartRepository) CreateCart(userID int, productID int, quantity int) (int, error) {
	argument := repo.Called(userID, productID, quantity)
	return argument.Int(0), argument.Error(1)
}

func (repo *mockCartRepository) UpdateCart(userID int, productID int, quantity int) error {
	argument := repo.Called(userID, productID, quantity)
	return argument.Error(0)
}

func (repo *mockCartRepository) DeleteCart(userID int, productID int) error {
	argument := repo.Called(userID, productID)
	return argument.Error(0)
}

type mockShippingRepository struct {
	mock.Mock
}

func (repo *mockShippingRepository) GetShippingMethodByID(ID int) (shipping.ShippingMethodDetail, error) {
	argument := repo.Called(ID)
	return argument.Get(0).(shipping.ShippingMethodDetail), argument.Error(1)
}

type mockUserRepository struct {
	mock.Mock
}

func (repo *mockUserRepository) FindByID(uid int) (auth.UserPayload, error) {
	args := repo.Called(uid)
	return args.Get(0).(auth.UserPayload), args.Error(1)
}

func (repo *mockUserRepository) FindByUsername(username string) (auth.User, error) {
	args := repo.Called(username)
	return args.Get(0).(auth.User), args.Error(1)
}
