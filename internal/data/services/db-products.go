package services

import (
	"bestprice/bestprice-api/internal/data/model"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type ProductsServices struct {
	DB         *sqlx.DB
	SqlQueries map[string]string
}

func NewProductsServices(db *sqlx.DB) *ProductsServices {
	return &ProductsServices{
		DB: db,
		SqlQueries: map[string]string{
			"select_all_products_ordered_limited": `SELECT * 
													FROM products 
													WHERE 1=1
													%s
													LIMIT ?,?;`,
			"select_product_by_id": `SELECT *
									FROM products
									WHERE 1=1
									AND id=?
									LIMIT 1;`,

			"insert_product": `INSERT INTO products (category_id, title, url_image, price, description) 
								VALUES (?,?,?,?,?);`,

			"update_product": `UPDATE products
								SET category_id=?, title=?, url_image=?, price=?, description=?, updated_on=?
								WHERE id=?;`,

			"delete_product": `DELETE 
								FROM products
								WHERE id=?;`,
		},
	}
}

func (o *ProductsServices) GetProducts(orderInput string, limitInput string, offsetInput string) ([]model.Product, error) {
	var products []model.Product
	var err error
	const (
		ORDER_BY = "ORDER BY id ASC"
		OFFSET   = "0"
		LIMIT    = "18446744073709551615" // bigint largest value (2^64-1) - ugly but it's proposed in MySQL guideline https://dev.mysql.com/doc/refman/8.0/en/select.html
	)

	order := ORDER_BY
	if orderInput != "" {
		order = orderInput
	}
	q := fmt.Sprintf(o.SqlQueries["select_all_products_ordered_limited"], order)

	limit := LIMIT
	if limitInput != "" {
		limit = limitInput
	}
	offset := OFFSET
	if offsetInput != "" {
		offset = offsetInput
	}

	rows, err := o.DB.Queryx(q, offset, limit)
	if err != nil {
		return products, err
	}

	var product model.Product
	for rows.Next() {
		err = rows.StructScan(&product)
		if err != nil {
			return nil, err
			//TODO : Handle errors with channel
		}
		products = append(products, product)
	}

	return products, nil
}

func (o *ProductsServices) GetProduct(idInput string) (model.Product, error) {
	var product model.Product
	var err error

	rows, err := o.DB.Queryx(o.SqlQueries["select_product_by_id"], idInput)
	if err != nil {
		return product, err
	}

	if !rows.Next() {
		return product, errors.New("Could not find the provided id")
	}

	err = rows.StructScan(&product)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (o *ProductsServices) AddProduct(product model.Product) error {
	_, err := o.DB.Exec(o.SqlQueries["insert_product"], product.CategoryId, product.Title, product.UrlImage, product.Price, product.Description)
	if err != nil {
		return err
	}

	return nil
}

func (o *ProductsServices) UpdateProduct(productId string, product model.Product) error {
	rows, err := o.DB.Exec(o.SqlQueries["update_product"], product.CategoryId, product.Title, product.UrlImage, product.Price, product.Description, time.Now(), productId)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if rowsAffected == 0 && err == nil {
		return errors.New("Could not update the record with this specific id")
	}

	return nil
}

func (o *ProductsServices) DeleteProduct(productId string) error {
	rows, err := o.DB.Exec(o.SqlQueries["delete_product"], productId)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if rowsAffected == 0 && err == nil {
		return errors.New("Could not delete the record with this specific id")
	}

	return nil

}
