package services

import (
	"bestprice/bestprice-api/internal/data/model"
	"database/sql"
	"time"
)

type ProductsServicesMock struct {
}

func NewProductsServicesMock() *ProductsServicesMock {
	return &ProductsServicesMock{}
}

func (o *ProductsServicesMock) GetProducts(orderInput string, limitInput string, offsetInput string) ([]model.Product, error) {
	products := []model.Product{
		model.Product{
			Archetype: model.Archetype{
				Id:        1,
				Title:     "product1",
				UrlImage:  "https://www.product1.com",
				CreatedOn: time.Now(),
				UpdatedOn: sql.NullTime{Time: time.Now(), Valid: true},
			},
			CategoryId:  sql.NullInt32{Int32: 1, Valid: true},
			Price:       1.1,
			Description: sql.NullString{String: "theproduct1", Valid: true},
		},
		model.Product{
			Archetype: model.Archetype{
				Id:        2,
				Title:     "product2",
				UrlImage:  "https://www.product2.com",
				CreatedOn: time.Now(),
				UpdatedOn: sql.NullTime{Time: time.Now(), Valid: true},
			},
			CategoryId:  sql.NullInt32{Int32: 1, Valid: true},
			Price:       2.2,
			Description: sql.NullString{String: "theproduct2", Valid: true},
		},
		model.Product{
			Archetype: model.Archetype{
				Id:        3,
				Title:     "product3",
				UrlImage:  "https://www.product3.com",
				CreatedOn: time.Now(),
				UpdatedOn: sql.NullTime{Time: time.Now(), Valid: true},
			},
			CategoryId:  sql.NullInt32{Int32: 2, Valid: true},
			Price:       3.3,
			Description: sql.NullString{String: "theproduct3", Valid: true},
		},
		model.Product{
			Archetype: model.Archetype{
				Id:        4,
				Title:     "product4",
				UrlImage:  "https://www.product4.com",
				CreatedOn: time.Now(),
				UpdatedOn: sql.NullTime{Time: time.Now(), Valid: true},
			},
			CategoryId:  sql.NullInt32{Int32: 2, Valid: true},
			Price:       4.4,
			Description: sql.NullString{String: "theproduct4", Valid: true},
		},
	}
	return products, nil
}

func (o *ProductsServicesMock) GetProduct(idInput string) (model.Product, error) {
	var product model.Product
	return product, nil
}

func (o *ProductsServicesMock) AddProduct(product model.Product) error {
	return nil
}

func (o *ProductsServicesMock) UpdateProduct(productId string, product model.Product) error {
	return nil
}

func (o *ProductsServicesMock) DeleteProduct(productId string) error {
	return nil
}
