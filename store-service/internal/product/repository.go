package product

import (
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	GetProducts(keyword string, limit string, offset string) (ProductResult, error)
	GetProductByID(ID int) (ProductDetail, error)
	UpdateStock(productID int, quantity int) error
}

type ProductRepositoryMySQL struct {
	DBConnection *sqlx.DB
}

func (repository ProductRepositoryMySQL) GetProducts(keyword string, limit string, offset string) (ProductResult, error) {
	var products []Product
	var total int
	if keyword == "" {
		err := repository.DBConnection.Select(&products, "SELECT id,product_name,product_price,image_url FROM products LIMIT ? OFFSET ?", limit, offset)
		if err == nil {
			err = repository.DBConnection.Get(&total, "SELECT count(*) FROM products")
		}

		return ProductResult{
			Total:    total,
			Products: products,
		}, err
	}
	err := repository.DBConnection.Select(&products, "SELECT id,product_name,product_price,image_url FROM products WHERE product_name LIKE ? LIMIT ? OFFSET ?", "%"+keyword+"%", limit, offset)
	if err == nil {
		err = repository.DBConnection.Get(&total, "SELECT count(*) FROM products WHERE product_name LIKE ?", "%"+keyword+"%")
	}

	return ProductResult{
		Total:    total,
		Products: products,
	}, err
}

func (productRepository ProductRepositoryMySQL) GetProductByID(ID int) (ProductDetail, error) {
	result := ProductDetail{}
	err := productRepository.DBConnection.Get(&result, "SELECT id,product_name,product_price,stock,image_url,product_brand FROM products WHERE id=?", ID)
	return result, err
}

func (productRepository ProductRepositoryMySQL) UpdateStock(productID int, stock int) error {
	_, err := productRepository.DBConnection.Exec(`UPDATE products SET stock = stock-? WHERE id=?`, stock, productID)
	return err
}
