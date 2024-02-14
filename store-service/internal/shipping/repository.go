package shipping

import (
	"github.com/jmoiron/sqlx"
)

type ShippingRepository interface {
	GetShippingMethodByID(ID int) (ShippingMethodDetail, error)
}

type ShippingRepositoryMySQL struct {
	DBConnection *sqlx.DB
}

func (shippingRepository ShippingRepositoryMySQL) GetShippingMethodByID(ID int) (ShippingMethodDetail, error) {
	result := ShippingMethodDetail{}
	err := shippingRepository.DBConnection.Get(&result, "SELECT id,name,description,fee FROM shipping_methods WHERE id=?", ID)
	return result, err
}
