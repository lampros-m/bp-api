package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Category struct {
	Archetype
	Position int `db:"place" json:"place"`
}

type CategoryExposed Category

func (o *CategoryExposed) MarshalJSON() ([]byte, error) {
	type AliasCategoryExposed CategoryExposed

	var updatedOnFormated string
	if o.UpdatedOn.Valid {
		updatedOnFormated = o.UpdatedOn.Time.Format("02-01-2006 15:04:05")
	}

	return json.Marshal(&struct {
		CreatedOn string `json:"created_on"`
		UpdatedOn string `json:"updated_on,omitempty"`
		*AliasCategoryExposed
	}{
		CreatedOn:            o.CreatedOn.Format("02-01-2006 15:04:05"),
		UpdatedOn:            updatedOnFormated,
		AliasCategoryExposed: (*AliasCategoryExposed)(o),
	})
}

func (o *CategoryExposed) UnmarshalJSON(data []byte) error {
	type AliasCategoryExposed CategoryExposed
	categoryExposedOverwritten := &struct {
		CreatedOn string `json:"created_on"`
		UpdatedOn string `json:"updated_on"`
		*AliasCategoryExposed
	}{
		AliasCategoryExposed: (*AliasCategoryExposed)(o),
	}

	if err := json.Unmarshal(data, &categoryExposedOverwritten); err != nil {
		return err
	}

	layoutTime := "02-01-2006 15:04:05"
	parsedCreatedOn, err := time.Parse(layoutTime, categoryExposedOverwritten.CreatedOn)
	if err != nil {
		return err
	}
	o.Archetype.CreatedOn = parsedCreatedOn

	parsedUpdatedOn, err := time.Parse(layoutTime, categoryExposedOverwritten.UpdatedOn)
	if err != nil {
		return err
	}

	o.Archetype.UpdatedOn = sql.NullTime{
		Time:  parsedUpdatedOn,
		Valid: true,
	}

	return nil
}
