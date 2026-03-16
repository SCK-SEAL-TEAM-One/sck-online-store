package shipping

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ShippingRepository interface {
	GetShippingMethodByID(ctx context.Context, ID int) (ShippingMethodDetail, error)
}

type ShippingRepositoryMySQL struct {
	DBConnection *sqlx.DB
}

func (shippingRepository ShippingRepositoryMySQL) GetShippingMethodByID(ctx context.Context, ID int) (ShippingMethodDetail, error) {
	result := ShippingMethodDetail{}
	err := shippingRepository.DBConnection.GetContext(ctx, &result, "SELECT id,name,description,fee FROM shipping_methods WHERE id=?", ID)
	return result, err
}
