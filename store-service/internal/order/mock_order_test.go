package order_test

import (
	"context"
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

func (service *mockPointInterface) TotalPoint(ctx context.Context, uid int) (point.TotalPoint, error) {
	argument := service.Called(ctx, uid)
	return argument.Get(0).(point.TotalPoint), argument.Error(1)
}

func (service *mockPointInterface) DeductPoint(ctx context.Context, uid int, submitedPoint point.SubmitedPoint) (point.TotalPoint, error) {
	argument := service.Called(ctx, uid, submitedPoint)
	return argument.Get(0).(point.TotalPoint), argument.Error(1)
}

func (service *mockPointInterface) CheckBurnPoint(ctx context.Context, uid int, amount int) (bool, error) {
	argument := service.Called(ctx, uid, amount)
	return argument.Bool(0), argument.Error(1)
}

type mockOrderHelper struct {
	mock.Mock
}

func (m *mockOrderHelper) GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq int, now time.Time) (int64, error) {
	args := m.Called(paymentMethodID, shippingMethodID, userID, seq, now)
	return args.Get(0).(int64), args.Error(1)
}

type mockOrderRepository struct {
	mock.Mock
}

func (repo *mockOrderRepository) CreateOrder(ctx context.Context, userID int, orderDetail order.OrderDetail) (int, error) {
	argument := repo.Called(ctx, userID, orderDetail)
	return argument.Int(0), argument.Error(1)
}

func (repo *mockOrderRepository) GetOrderByOrderNumber(ctx context.Context, orderNumber int64) (order.OrderDetail, error) {
	argument := repo.Called(ctx, orderNumber)
	return argument.Get(0).(order.OrderDetail), argument.Error(1)
}

func (repo *mockOrderRepository) GetNextSequence(ctx context.Context, yearPrefix string, userID int) (int, error) {
	argument := repo.Called(ctx, yearPrefix, userID)
	return argument.Int(0), argument.Error(1)
}

func (repo *mockOrderRepository) GetOrderWithTrackingNumberByOrderNumber(ctx context.Context, orderNumber int64) (order.OrderDetailWithTrackingNumber, error) {
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
	return argument.Error(1)
}

func (repo *mockOrderRepository) UpdateOrderTrackingNumber(ctx context.Context, orderID int, trackingNumber string) error {
	argument := repo.Called(ctx, orderID, trackingNumber)
	return argument.Error(1)
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

type mockCartRepository struct {
	mock.Mock
}

func (repo *mockCartRepository) GetCartDetail(ctx context.Context, userID int) ([]cart.CartDetail, error) {
	argument := repo.Called(ctx, userID)
	return argument.Get(0).([]cart.CartDetail), argument.Error(1)
}

func (repo *mockCartRepository) GetCartByProductID(ctx context.Context, userID int, productID int) (cart.Cart, error) {
	argument := repo.Called(ctx, userID, productID)
	return argument.Get(0).(cart.Cart), argument.Error(1)
}

func (repo *mockCartRepository) CreateCart(ctx context.Context, userID int, productID int, quantity int) (int, error) {
	argument := repo.Called(ctx, userID, productID, quantity)
	return argument.Int(0), argument.Error(1)
}

func (repo *mockCartRepository) UpdateCart(ctx context.Context, userID int, productID int, quantity int) error {
	argument := repo.Called(ctx, userID, productID, quantity)
	return argument.Error(0)
}

func (repo *mockCartRepository) DeleteCart(ctx context.Context, userID int, productID int) error {
	argument := repo.Called(ctx, userID, productID)
	return argument.Error(0)
}

type mockShippingRepository struct {
	mock.Mock
}

func (repo *mockShippingRepository) GetShippingMethodByID(ctx context.Context, ID int) (shipping.ShippingMethodDetail, error) {
	argument := repo.Called(ctx, ID)
	return argument.Get(0).(shipping.ShippingMethodDetail), argument.Error(1)
}

type mockUserRepository struct {
	mock.Mock
}

func (repo *mockUserRepository) FindByID(ctx context.Context, uid int) (auth.UserPayload, error) {
	args := repo.Called(ctx, uid)
	return args.Get(0).(auth.UserPayload), args.Error(1)
}

func (repo *mockUserRepository) FindByUsername(ctx context.Context, username string) (auth.User, error) {
	args := repo.Called(ctx, username)
	return args.Get(0).(auth.User), args.Error(1)
}
