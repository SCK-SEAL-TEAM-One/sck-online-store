package order

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	CreateOrder(userID int, orderDetail OrderDetail) (int, error)
	GetOrderByOrderNumber(orderNumber string) (OrderDetail, error)
	GetLastOrderNumber(yearPrefix string) (string, error)
	GetOrderWithTrackingNumberByOrderNumber(orderNumber string) (OrderDetailWithTrackingNumber, error)
	CreateOrderProduct(orderID, productID, quantity int, productPrice float64) error
	UpdateOrderTransaction(orderID int, transactionID string) error
	UpdateOrderTrackingNumber(orderID int, trackingNumber string) error
	GetOrderProduct(orderID int) ([]OrderProduct, error)
	GetOrderProductWithPrice(orderID int) ([]OrderProductWithPrice, error)
	CreateShipping(userID int, orderID int, shippingInfo ShippingInfo) (int, error)
}

type OrderRepositoryMySQL struct {
	DBConnection *sqlx.DB
}

func (orderRepository OrderRepositoryMySQL) CreateOrder(userID int, orderDetail OrderDetail) (int, error) {
	query := `
		INSERT INTO orders (
			user_id,
			order_number,
			shipping_method_id,
			payment_method_id,
			sub_total_price,
			discount_price,
			total_price,
			shipping_fee,
			burn_point,
			earn_point
		) VALUES (?,?,?,?,?,?,?,?,?,?)`

	sqlResult := orderRepository.DBConnection.MustExec(query, userID, orderDetail.OrderNumber, orderDetail.ShippingMethodID, orderDetail.PaymentMethodID, orderDetail.SubTotalPrice, orderDetail.DiscountPrice, orderDetail.TotalPrice, orderDetail.ShippingFee, orderDetail.BurnPoint, orderDetail.EarnPoint)
	insertedId, err := sqlResult.LastInsertId()
	return int(insertedId), err
}

func (orderRepository OrderRepositoryMySQL) GetOrderByOrderNumber(orderNumber string) (OrderDetail, error) {
	result := OrderDetail{}
	err := orderRepository.DBConnection.Get(&result, `
		SELECT id, order_number, user_id, shipping_method_id, payment_method_id, sub_total_price, discount_price, total_price, shipping_fee, burn_point, earn_point, transaction_id, status
		FROM orders WHERE order_number = ?
	`, orderNumber)
	return result, err
}

func (orderRepository OrderRepositoryMySQL) GetLastOrderNumber(yearPrefix string) (string, error) {
	lastOrderNumber := ""
	pattern := yearPrefix + "%"
	query := `
		SELECT order_number
		FROM orders
		WHERE order_number LIKE ?
		ORDER BY updated DESC
		LIMIT 1`
	err := orderRepository.DBConnection.Get(&lastOrderNumber, query, pattern)
	return lastOrderNumber, err
}

func (orderRepository OrderRepositoryMySQL) GetOrderWithTrackingNumberByOrderNumber(orderNumber string) (OrderDetailWithTrackingNumber, error) {
	result := OrderDetailWithTrackingNumber{}
	err := orderRepository.DBConnection.Get(&result, `
		SELECT id, order_number, user_id, shipping_method_id, payment_method_id, sub_total_price, discount_price, total_price, shipping_fee, burn_point, earn_point, transaction_id, status, tracking_no, updated
		FROM orders WHERE order_number = ?
	`, orderNumber)
	return result, err
}

func (orderRepository OrderRepositoryMySQL) CreateOrderProduct(orderID int, productID, quantity int, productPrice float64) error {
	sqlResult := orderRepository.DBConnection.MustExec("INSERT INTO order_product (order_id, product_id, quantity, product_price) VALUE (?,?,?,?)", orderID, productID, quantity, productPrice)
	_, err := sqlResult.RowsAffected()
	return err
}

func (orderRepository OrderRepositoryMySQL) UpdateOrderTransaction(orderID int, transactionID string) error {
	status := "paid"
	sqlResult := orderRepository.DBConnection.MustExec("UPDATE orders SET transaction_id=? , status=? WHERE id = ?", transactionID, status, orderID)
	rowAffected, err := sqlResult.RowsAffected()
	if rowAffected == 0 {
		return fmt.Errorf("no any row affected , update not completed")
	}
	return err
}

func (orderRepository OrderRepositoryMySQL) UpdateOrderTrackingNumber(orderID int, trackingNumber string) error {
	sqlResult := orderRepository.DBConnection.MustExec("UPDATE orders SET tracking_no=? WHERE id = ?", trackingNumber, orderID)
	rowAffected, err := sqlResult.RowsAffected()
	if rowAffected == 0 {
		return fmt.Errorf("no any row affected , update not completed")
	}
	return err
}

func (repository OrderRepositoryMySQL) GetOrderProduct(orderID int) ([]OrderProduct, error) {
	var orderProducts []OrderProduct
	err := repository.DBConnection.Select(&orderProducts, "SELECT product_id, quantity FROM order_product WHERE order_id = ?", orderID)
	return orderProducts, err
}

func (repository OrderRepositoryMySQL) GetOrderProductWithPrice(orderID int) ([]OrderProductWithPrice, error) {
	var orderProducts []OrderProductWithPrice
	err := repository.DBConnection.Select(&orderProducts, "SELECT p.product_brand, p.product_name, op.quantity, op.product_price FROM products p JOIN order_product op ON p.id = op.product_id WHERE order_id = ?", orderID)
	return orderProducts, err
}

func (orderRepository OrderRepositoryMySQL) CreateShipping(userID int, orderID int, shippingInfo ShippingInfo) (int, error) {
	result := orderRepository.
		DBConnection.
		MustExec(`INSERT INTO order_shipping (order_id, user_id, method_id, address, sub_district, district, province, zip_code, recipient_first_name, recipient_last_name, phone_number) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			orderID,
			userID,
			shippingInfo.ShippingMethodID,
			shippingInfo.ShippingAddress,
			shippingInfo.ShippingSubDistrict,
			shippingInfo.ShippingDistrict,
			shippingInfo.ShippingProvince,
			shippingInfo.ShippingZipCode,
			shippingInfo.RecipientFirstName,
			shippingInfo.RecipientLastName,
			shippingInfo.RecipientPhoneNumber,
		)
	id, err := result.LastInsertId()
	return int(id), err
}
