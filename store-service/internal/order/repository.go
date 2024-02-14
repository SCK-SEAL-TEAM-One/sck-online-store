package order

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	CreateOrder(userID int, orderDetail OrderDetail) (int, error)
	GetOrderByID(ID int) (OrderDetail, error)
	CreateOrderProduct(orderID, productID, quantity int, productPrice float64) error
	UpdateOrderTransaction(orderID int, transactionID string) error
	GetOrderProduct(orderID int) ([]OrderProduct, error)
	CreateShipping(userID int, orderID int, shippingInfo ShippingInfo) (int, error)
}

type OrderRepositoryMySQL struct {
	DBConnection *sqlx.DB
}

func (orderRepository OrderRepositoryMySQL) CreateOrder(userID int, orderDetail OrderDetail) (int, error) {
	sqlResult := orderRepository.DBConnection.MustExec("INSERT INTO orders (user_id, shipping_method_id, payment_method_id, sub_total_price, discount_price, total_price, shipping_fee, burn_point, earn_point) VALUE (?,?,?,?,?,?,?,?,?)", userID, orderDetail.ShippingMethodID, orderDetail.PaymentMethodID, orderDetail.SubTotalPrice, orderDetail.DiscountPrice, orderDetail.TotalPrice, orderDetail.ShippingFee, orderDetail.BurnPoint, orderDetail.EarnPoint)
	insertedId, err := sqlResult.LastInsertId()
	return int(insertedId), err
}

func (orderRepository OrderRepositoryMySQL) GetOrderByID(ID int) (OrderDetail, error) {
	result := OrderDetail{}
	err := orderRepository.DBConnection.Get(&result, `
		SELECT id,user_id,shipping_method_id,payment_method_id,sub_total_price,discount_price,total_price,shipping_fee,burn_point,earn_point,transaction_id,status
		FROM orders WHERE id=?
	`, ID)
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

func (repository OrderRepositoryMySQL) GetOrderProduct(orderID int) ([]OrderProduct, error) {
	var orderProducts []OrderProduct
	err := repository.DBConnection.Select(&orderProducts, "SELECT product_id, quantity FROM order_product WHERE order_id = ?", orderID)
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
