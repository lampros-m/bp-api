package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Product struct {
	Archetype
	CategoryId  sql.NullInt32  `db:"category_id" json:"category_id"`
	Price       float32        `db:"price" json:"price"`
	Description sql.NullString `db:"description" json:"description"`
}

type ProductExposed Product

func (o *ProductExposed) MarshalJSON() ([]byte, error) {
	type AliasProductExposed ProductExposed

	var updatedOnFormated string
	if o.UpdatedOn.Valid {
		updatedOnFormated = o.UpdatedOn.Time.Format("02-01-2006 15:04:05")
	}

	return json.Marshal(&struct {
		CreatedOn   string `json:"created_on"`
		UpdatedOn   string `json:"updated_on,omitempty"`
		CategoryId  int32  `json:"category_id"`
		Description string `json:"description"`
		*AliasProductExposed
	}{
		CreatedOn:           o.CreatedOn.Format("02-01-2006 15:04:05"),
		UpdatedOn:           updatedOnFormated,
		CategoryId:          o.CategoryId.Int32,
		Description:         o.Description.String,
		AliasProductExposed: (*AliasProductExposed)(o),
	})
}

func (o *Product) UnmarshalJSON(data []byte) error {
	type AliasProduct Product
	productOverwritten := &struct {
		CategoryId  int32  `json:"category_id"`
		Description string `json:"description"`
		*AliasProduct
	}{
		AliasProduct: (*AliasProduct)(o),
	}

	if err := json.Unmarshal(data, &productOverwritten); err != nil {
		return err
	}

	o.CategoryId = sql.NullInt32{
		Int32: productOverwritten.CategoryId,
		Valid: true,
	}
	o.Description = sql.NullString{
		String: productOverwritten.Description,
		Valid:  true,
	}

	return nil
}

func (o *ProductExposed) UnmarshalJSON(data []byte) error {
	type AliasProductExposed ProductExposed
	productExposedOverwritten := &struct {
		CategoryId  int32  `json:"category_id"`
		Description string `json:"description"`
		CreatedOn   string `json:"created_on"`
		UpdatedOn   string `json:"updated_on"`
		*AliasProductExposed
	}{
		AliasProductExposed: (*AliasProductExposed)(o),
	}

	if err := json.Unmarshal(data, &productExposedOverwritten); err != nil {
		return err
	}

	layoutTime := "02-01-2006 15:04:05"
	parsedCreatedOn, err := time.Parse(layoutTime, productExposedOverwritten.CreatedOn)
	if err != nil {
		return err
	}
	o.Archetype.CreatedOn = parsedCreatedOn

	parsedUpdatedOn, err := time.Parse(layoutTime, productExposedOverwritten.UpdatedOn)
	if err != nil {
		return err
	}

	o.Archetype.UpdatedOn = sql.NullTime{
		Time:  parsedUpdatedOn,
		Valid: true,
	}

	o.CategoryId = sql.NullInt32{
		Int32: productExposedOverwritten.CategoryId,
		Valid: true,
	}
	o.Description = sql.NullString{
		String: productExposedOverwritten.Description,
		Valid:  true,
	}

	return nil
}
