//go:build integration
// +build integration

package shipping_test

import (
	"store-service/internal/shipping"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_ShippingRepository(t *testing.T) {
	connection, err := sqlx.Connect("mysql", "user:password@(localhost:3306)/store")
	if err != nil {
		t.Fatalf("cannot tearup data err %s", err)
	}
	repository := shipping.ShippingRepositoryMySQL{
		DBConnection: connection,
	}

	t.Run("GetShippingMethodByID_Input_ID_1_Should_Be_ShippingMethod_Detail_No_Error", func(t *testing.T) {
		expected := shipping.ShippingMethodDetail{
			ID:          1,
			Name:        "Kerry",
			Description: "4-5 business days",
			Fee:         50,
		}
		ID := 1

		actualProduct, err := repository.GetShippingMethodByID(ID)
		assert.Equal(t, expected, actualProduct)
		assert.Equal(t, err, nil)
	})
}
