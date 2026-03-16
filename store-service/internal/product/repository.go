package product

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	GetProducts(ctx context.Context, keyword string, limit string, offset string) (ProductResult, error)
	GetProductByID(ctx context.Context, ID int) (ProductDetail, error)
	UpdateStock(ctx context.Context, productID int, quantity int) error
}

type ProductRepositoryMySQL struct {
	DBConnection *sqlx.DB
}

func (repository ProductRepositoryMySQL) GetProducts(ctx context.Context, keyword string, limit string, offset string) (ProductResult, error) {
	var products []Product
	var total int
	if keyword == "" {
		err := repository.DBConnection.SelectContext(ctx, &products, "SELECT id,product_name,product_price,image_url FROM products LIMIT ? OFFSET ?", limit, offset)
		if err == nil {
			err = repository.DBConnection.GetContext(ctx, &total, "SELECT count(*) FROM products")
		}

		return ProductResult{
			Total:    total,
			Products: products,
		}, err
	}
	err := repository.DBConnection.SelectContext(ctx, &products, "SELECT id,product_name,product_price,image_url FROM products WHERE product_name LIKE ? LIMIT ? OFFSET ?", "%"+keyword+"%", limit, offset)
	if err == nil {
		err = repository.DBConnection.GetContext(ctx, &total, "SELECT count(*) FROM products WHERE product_name LIKE ?", "%"+keyword+"%")
	}

	return ProductResult{
		Total:    total,
		Products: products,
	}, err
}

func (productRepository ProductRepositoryMySQL) GetProductByID(ctx context.Context, ID int) (ProductDetail, error) {
	result := ProductDetail{}
	err := productRepository.DBConnection.GetContext(ctx, &result, "SELECT id,product_name,product_price,stock,image_url,product_brand FROM products WHERE id=?", ID)
	return result, err
}

func (productRepository ProductRepositoryMySQL) UpdateStock(ctx context.Context, productID int, stock int) error {
	_, err := productRepository.DBConnection.ExecContext(ctx, `UPDATE products SET stock = stock-? WHERE id=?`, stock, productID)
	return err
}
