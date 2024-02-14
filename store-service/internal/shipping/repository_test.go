//go:build integration
// +build integration

package product_test

import (
	"store-service/internal/product"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_ShippingRepository(t *testing.T) {
	connection, err := sqlx.Connect("mysql", "user:password@(localhost:3305)/store")
	if err != nil {
		t.Fatalf("cannot tearup data err %s", err)
	}
	repository := product.ProductRepositoryMySQL{
		DBConnection: connection,
	}

	t.Run("GetShippingMethodByID_Input_ID_1_Should_Be_ShippingMethod_Detail_No_Error", func(t *testing.T) {
		expected := product.ProductDetail{
			ID:          1,
			Name:        "Kerry",
			Description: "4-5 business days",
			Fee:         50,
		}
		ID := 2

		actualProduct, err := repository.GetProductByID(ID)
		assert.Equal(t, expected, actualProduct)
		assert.Equal(t, err, nil)
	})
}
