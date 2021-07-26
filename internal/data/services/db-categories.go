package services

import (
	"bestprice/bestprice-api/internal/data/model"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type CategoriesServices struct {
	DB         *sqlx.DB
	SqlQueries map[string]string
}

func NewCategoriesServices(db *sqlx.DB) *CategoriesServices {
	return &CategoriesServices{
		DB: db,
		SqlQueries: map[string]string{
			"select_all_categories_ordered_limited": `SELECT * 
													FROM categories 
													WHERE 1=1
													%s
													LIMIT ?,?;`,

			"select_category_by_id": `SELECT *
									FROM categories
									WHERE 1=1
									AND id=?
									LIMIT 1;`,

			"insert_category": `INSERT INTO categories (title, place, url_image) 
								VALUES (?,?,?);`,

			"update_category": `UPDATE categories
								SET title=?, place=?, url_image=?, updated_on=?
								WHERE id=?;`,

			"delete_category": `DELETE 
								FROM categories
								WHERE id=?;`,
		},
	}
}

func (o *CategoriesServices) GetCategories(orderInput string, limitInput string, offsetInput string) ([]model.Category, error) {
	var categories []model.Category
	var err error
	const (
		ORDER_BY = "ORDER BY place ASC"
		OFFSET   = "0"
		LIMIT    = "18446744073709551615" // bigint largest value (2^64-1) - ugly but it's proposed in MySQL guideline https://dev.mysql.com/doc/refman/8.0/en/select.html
	)

	order := ORDER_BY
	if orderInput != "" {
		order = orderInput
	}
	q := fmt.Sprintf(o.SqlQueries["select_all_categories_ordered_limited"], order)

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
		return categories, err
	}

	var category model.Category
	for rows.Next() {
		err = rows.StructScan(&category)
		if err != nil {
			return nil, err
			//TODO : Handle errors with channel
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (o *CategoriesServices) GetCategory(idInput string) (model.Category, error) {
	var category model.Category
	var err error

	rows, err := o.DB.Queryx(o.SqlQueries["select_category_by_id"], idInput)
	if err != nil {
		return category, err
	}

	if !rows.Next() {
		return category, errors.New("Could not find the provided id")
	}

	err = rows.StructScan(&category)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (o *CategoriesServices) AddCategory(category model.Category) error {
	_, err := o.DB.Exec(o.SqlQueries["insert_category"], category.Title, category.Position, category.UrlImage)
	if err != nil {
		return err
	}

	return nil
}

func (o *CategoriesServices) UpdateCategory(categoryId string, category model.Category) error {
	rows, err := o.DB.Exec(o.SqlQueries["update_category"], category.Title, category.Position, category.UrlImage, time.Now(), categoryId)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if rowsAffected == 0 && err == nil {
		return errors.New("Could not update the record with this specific id")
	}

	return nil
}

func (o *CategoriesServices) DeleteCategory(categoryId string) error {
	rows, err := o.DB.Exec(o.SqlQueries["delete_category"], categoryId)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if rowsAffected == 0 && err == nil {
		return errors.New("Could not delete the record with this specific id")
	}

	return nil
}
