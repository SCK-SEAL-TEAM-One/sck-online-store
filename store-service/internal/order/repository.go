package order

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, userID int, orderDetail OrderDetail) (int, error)
	GetOrderByOrderNumber(ctx context.Context, orderNumber int64) (OrderDetail, error)
	GetNextSequence(ctx context.Context, datePrefix string, userID int) (int, error)
	GetOrderWithTrackingNumberByOrderNumber(ctx context.Context, orderNumber int64) (OrderDetailWithTrackingNumber, error)
	CreateOrderProduct(ctx context.Context, orderID, productID, quantity int, productPrice float64) error
	UpdateOrderTransaction(ctx context.Context, orderID int, transactionID string) error
	UpdateOrderTrackingNumber(ctx context.Context, orderID int, trackingNumber string) error
	GetOrderProduct(ctx context.Context, orderID int) ([]OrderProduct, error)
	GetOrderProductWithPrice(ctx context.Context, orderID int) ([]OrderProductWithPrice, error)
	CreateShipping(ctx context.Context, userID int, orderID int, shippingInfo ShippingInfo) (int, error)
}

type OrderRepositoryMySQL struct {
	DBConnection *sqlx.DB
}

func (orderRepository OrderRepositoryMySQL) CreateOrder(ctx context.Context, userID int, orderDetail OrderDetail) (int, error) {
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

	sqlResult, err := orderRepository.DBConnection.ExecContext(ctx, query, userID, orderDetail.OrderNumber, orderDetail.ShippingMethodID, orderDetail.PaymentMethodID, orderDetail.SubTotalPrice, orderDetail.DiscountPrice, orderDetail.TotalPrice, orderDetail.ShippingFee, orderDetail.BurnPoint, orderDetail.EarnPoint)
	if err != nil {
		return 0, err
	}
	insertedId, err := sqlResult.LastInsertId()
	return int(insertedId), err
}

func (orderRepository OrderRepositoryMySQL) GetOrderByOrderNumber(ctx context.Context, orderNumber int64) (OrderDetail, error) {
	result := OrderDetail{}
	err := orderRepository.DBConnection.GetContext(ctx, &result, `
		SELECT id, order_number, user_id, shipping_method_id, payment_method_id, sub_total_price, discount_price, total_price, shipping_fee, burn_point, earn_point, transaction_id, status
		FROM orders WHERE order_number = ?
	`, orderNumber)
	return result, err
}

func (orderRepository OrderRepositoryMySQL) GetNextSequence(ctx context.Context, datePrefix string, userID int) (int, error) {
	query := `
		INSERT INTO order_sequences (date_prefix, user_id, current_seq)
		VALUES (?, ?, LAST_INSERT_ID(1))
		ON DUPLICATE KEY UPDATE current_seq = LAST_INSERT_ID(current_seq + 1)`
	result, err := orderRepository.DBConnection.ExecContext(ctx, query, datePrefix, userID)
	if err != nil {
		return 0, err
	}
	seq, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(seq), nil
}

func (orderRepository OrderRepositoryMySQL) GetOrderWithTrackingNumberByOrderNumber(ctx context.Context, orderNumber int64) (OrderDetailWithTrackingNumber, error) {
	result := OrderDetailWithTrackingNumber{}
	err := orderRepository.DBConnection.GetContext(ctx, &result, `
		SELECT id, order_number, user_id, shipping_method_id, payment_method_id, sub_total_price, discount_price, total_price, shipping_fee, burn_point, earn_point, transaction_id, status, tracking_no, updated
		FROM orders WHERE order_number = ?
	`, orderNumber)
	return result, err
}

func (orderRepository OrderRepositoryMySQL) CreateOrderProduct(ctx context.Context, orderID int, productID, quantity int, productPrice float64) error {
	sqlResult, err := orderRepository.DBConnection.ExecContext(ctx, "INSERT INTO order_product (order_id, product_id, quantity, product_price) VALUE (?,?,?,?)", orderID, productID, quantity, productPrice)
	if err != nil {
		return err
	}
	_, err = sqlResult.RowsAffected()
	return err
}

func (orderRepository OrderRepositoryMySQL) UpdateOrderTransaction(ctx context.Context, orderID int, transactionID string) error {
	status := "paid"
	sqlResult, err := orderRepository.DBConnection.ExecContext(ctx, "UPDATE orders SET transaction_id=? , status=? WHERE id = ?", transactionID, status, orderID)
	if err != nil {
		return err
	}
	rowAffected, err := sqlResult.RowsAffected()
	if rowAffected == 0 {
		return fmt.Errorf("no any row affected , update not completed")
	}
	return err
}

func (orderRepository OrderRepositoryMySQL) UpdateOrderTrackingNumber(ctx context.Context, orderID int, trackingNumber string) error {
	sqlResult, err := orderRepository.DBConnection.ExecContext(ctx, "UPDATE orders SET tracking_no=? WHERE id = ?", trackingNumber, orderID)
	if err != nil {
		return err
	}
	rowAffected, err := sqlResult.RowsAffected()
	if rowAffected == 0 {
		return fmt.Errorf("no any row affected , update not completed")
	}
	return err
}

func (repository OrderRepositoryMySQL) GetOrderProduct(ctx context.Context, orderID int) ([]OrderProduct, error) {
	var orderProducts []OrderProduct
	err := repository.DBConnection.SelectContext(ctx, &orderProducts, "SELECT product_id, quantity FROM order_product WHERE order_id = ?", orderID)
	return orderProducts, err
}

func (repository OrderRepositoryMySQL) GetOrderProductWithPrice(ctx context.Context, orderID int) ([]OrderProductWithPrice, error) {
	var orderProducts []OrderProductWithPrice
	err := repository.DBConnection.SelectContext(ctx, &orderProducts, "SELECT p.product_brand, p.product_name, op.quantity, op.product_price FROM products p JOIN order_product op ON p.id = op.product_id WHERE order_id = ?", orderID)
	return orderProducts, err
}

func (orderRepository OrderRepositoryMySQL) CreateShipping(ctx context.Context, userID int, orderID int, shippingInfo ShippingInfo) (int, error) {
	result, err := orderRepository.
		DBConnection.
		ExecContext(ctx, `INSERT INTO order_shipping (order_id, user_id, method_id, address, sub_district, district, province, zip_code, recipient_first_name, recipient_last_name, phone_number) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
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
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}
