package cart

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CartRepository interface {
	GetCartDetail(ctx context.Context, userID int) ([]CartDetail, error)
	GetCartByProductID(ctx context.Context, userID int, productID int) (Cart, error)
	CreateCart(ctx context.Context, userID int, productID int, quantity int) (int, error)
	UpdateCart(ctx context.Context, userID int, productID int, quantity int) error
	DeleteCart(ctx context.Context, userID int, productID int) error
}

type CartRepositoryMySQL struct {
	DBConnection *sqlx.DB
}

func (repository CartRepositoryMySQL) GetCartDetail(ctx context.Context, userID int) ([]CartDetail, error) {
	var carts []CartDetail
	err := repository.DBConnection.SelectContext(ctx, &carts, `
		SELECT c.id, c.user_id, c.product_id, c.quantity, p.product_name, p.product_brand, p.stock, p.product_price, p.image_url
		FROM carts c
		LEFT JOIN products p ON c.product_id  = p.id
		WHERE  c.user_id = ?
	`, userID)
	return carts, err
}

func (repository CartRepositoryMySQL) GetCartByProductID(ctx context.Context, userID int, productID int) (Cart, error) {
	var cart Cart
	err := repository.DBConnection.GetContext(ctx, &cart, "SELECT id,user_id,product_id,quantity FROM carts WHERE user_id = ? AND product_id = ? LIMIT 1", userID, productID)
	return cart, err
}

func (repository CartRepositoryMySQL) CreateCart(ctx context.Context, userID int, productID int, quantity int) (int, error) {
	sqlResult, err := repository.DBConnection.ExecContext(ctx, "INSERT INTO carts (user_id, product_id, quantity) VALUE (?,?,?)", userID, productID, quantity)
	if err != nil {
		return 0, err
	}
	insertedId, err := sqlResult.LastInsertId()
	return int(insertedId), err
}

func (repository CartRepositoryMySQL) UpdateCart(ctx context.Context, userID int, productID int, quantity int) error {
	sqlResult, err := repository.DBConnection.ExecContext(ctx, "UPDATE carts SET quantity=? WHERE user_id = ? AND product_id = ?", quantity, userID, productID)
	if err != nil {
		return err
	}
	rowAffected, err := sqlResult.RowsAffected()
	if rowAffected == 0 {
		return fmt.Errorf("no any row affected , update not completed")
	}
	return err
}

func (repository CartRepositoryMySQL) DeleteCart(ctx context.Context, userID int, productID int) error {
	sqlResult, err := repository.DBConnection.ExecContext(ctx, "DELETE FROM carts WHERE user_id = ? AND product_id = ?", userID, productID)
	if err != nil {
		return err
	}
	rowAffected, err := sqlResult.RowsAffected()
	if rowAffected == 0 {
		return fmt.Errorf("no any row affected , delete not completed")
	}
	return err
}
