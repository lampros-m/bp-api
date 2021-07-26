package database

import (
	"bestprice/bestprice-api/internal/data/services"
)

func NewDbMock() *DB {
	db := DB{
		Categories: services.NewCategoriesServicesMock(),
		Products:   services.NewProductsServicesMock(),
		//Users:      services.NewUsersServices(),
	}

	return &db
}
