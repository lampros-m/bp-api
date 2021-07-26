package services

import (
	"bestprice/bestprice-api/internal/data/model"
	"database/sql"
	"time"
)

type CategoriesServicesMock struct {
}

func NewCategoriesServicesMock() *CategoriesServicesMock {
	return &CategoriesServicesMock{}
}

func (o *CategoriesServicesMock) GetCategories(orderInput string, limitInput string, offsetInput string) ([]model.Category, error) {
	categories := []model.Category{
		model.Category{
			Archetype: model.Archetype{
				Id:        1,
				Title:     "category1",
				UrlImage:  "https://www.category1.com",
				CreatedOn: time.Now(),
				UpdatedOn: sql.NullTime{Time: time.Now(), Valid: true},
			},
			Position: 1,
		},
		model.Category{
			Archetype: model.Archetype{
				Id:        2,
				Title:     "category2",
				UrlImage:  "https://www.category2.com",
				CreatedOn: time.Now(),
				UpdatedOn: sql.NullTime{Time: time.Now(), Valid: true},
			},
			Position: 2,
		},
	}

	return categories, nil
}

func (o *CategoriesServicesMock) GetCategory(idInput string) (model.Category, error) {
	var category model.Category
	return category, nil
}

func (o *CategoriesServicesMock) AddCategory(category model.Category) error {
	return nil
}

func (o *CategoriesServicesMock) UpdateCategory(categoryId string, category model.Category) error {
	return nil
}

func (o *CategoriesServicesMock) DeleteCategory(categoryId string) error {
	return nil
}
