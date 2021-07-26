package api

import (
	"bestprice/bestprice-api/internal/data/model"
	"bestprice/bestprice-api/internal/helper"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	REDISTAG_CATEGORY        = "category.id:"
	REDISTAG_CATEGORIES_LIST = "category.list:"
)

func (a *Api) listCategories(w http.ResponseWriter, r *http.Request) {
	var categories []model.Category
	var err error

	params := r.URL.Query()

	sortParameter := strings.Join(params["sort"], ",")
	sortingSql, err := helper.ExportSqlSorting(sortParameter)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	limitParameter := strings.Join(params["limit"], ",")
	offsetParameter := strings.Join(params["offset"], ",")
	limitSql, offsetSql, err := helper.ExportSqlPagination(limitParameter, offsetParameter)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	redisTagKey := fmt.Sprintf("%ssort[%s].limit[%s].offset[%s]", REDISTAG_CATEGORIES_LIST, sortParameter, limitParameter, offsetParameter)
	flagRedisKeyFound := false

	value, err := a.Redis.RedisServices.Cache.Get(redisTagKey)
	if err == nil {
		flagRedisKeyFound = true
	}

	if flagRedisKeyFound {
		a.renderJson(w, r, http.StatusOK, value)
		return
	}

	categories, err = a.Db.Categories.GetCategories(sortingSql, limitSql, offsetSql)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	categoriesRedisValue, err := json.Marshal(helper.ExposeCategories(categories))
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	err = a.Redis.RedisServices.Cache.Set(redisTagKey, string(categoriesRedisValue))
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	a.renderJson(w, r, http.StatusOK, helper.ExposeCategories(categories))
}

func (a *Api) readCategory(w http.ResponseWriter, r *http.Request) {
	categoryId := mux.Vars(r)["id"]

	if !helper.SupportedSqlString(categoryId) {
		apiError := BadRequestError(errors.New("Not supported characters in id provided"))
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	redisTagKey := REDISTAG_CATEGORY + categoryId
	flagRedisKeyFound := false

	value, err := a.Redis.RedisServices.Cache.Get(redisTagKey)
	if err == nil {
		flagRedisKeyFound = true
	}

	if flagRedisKeyFound {
		a.renderJson(w, r, http.StatusOK, value)
		return
	}

	category, err := a.Db.Categories.GetCategory(categoryId)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	categoryRedisValue, err := json.Marshal(func(c model.CategoryExposed) *model.CategoryExposed { return &c }(model.CategoryExposed(category)))
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	err = a.Redis.RedisServices.Cache.Set(redisTagKey, string(categoryRedisValue))
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	a.renderJson(w, r, http.StatusOK, func(c model.CategoryExposed) *model.CategoryExposed { return &c }(model.CategoryExposed(category)))
}

func (a *Api) createCategory(w http.ResponseWriter, r *http.Request) {
	body, err := helper.ReadHttpRequestAndClose(r)
	if err != nil {
		apiError := BadRequestError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	var category model.Category
	err = json.Unmarshal(body, &category)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	err = a.Db.Categories.AddCategory(category)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	a.renderJson(w, r, http.StatusOK, Broadcast{Message: "Category added"})
}

func (a *Api) updateCategory(w http.ResponseWriter, r *http.Request) {
	categoryId := mux.Vars(r)["id"]

	if !helper.SupportedSqlString(categoryId) {
		apiError := BadRequestError(errors.New("Not supported characters in id provided"))
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	redisTagKey := REDISTAG_CATEGORY + categoryId
	err := a.Redis.RedisServices.Cache.Delete(redisTagKey)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	body, err := helper.ReadHttpRequestAndClose(r)
	if err != nil {
		apiError := BadRequestError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	var category model.Category
	err = json.Unmarshal(body, &category)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	err = a.Db.Categories.UpdateCategory(categoryId, category)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	keysMustDeleted, err := a.Redis.Client.Keys(REDISTAG_CATEGORIES_LIST + "*").Result()
	for i := range keysMustDeleted {
		err = a.Redis.RedisServices.Cache.Delete(keysMustDeleted[i])
		if err != nil {
			apiError := InternalError(err)
			a.renderJson(w, r, apiError.Code, apiError)
			return
		}
	}

	a.renderJson(w, r, http.StatusOK, Broadcast{Message: "Category updated successfully"})
}

func (a *Api) deleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryId := mux.Vars(r)["id"]

	if !helper.SupportedSqlString(categoryId) {
		apiError := BadRequestError(errors.New("Not supported characters in id provided"))
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	redisTagKey := REDISTAG_CATEGORY + categoryId
	err := a.Redis.RedisServices.Cache.Delete(redisTagKey)
	if err != nil {
		apiError := BadRequestError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	err = a.Db.Categories.DeleteCategory(categoryId)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	keysMustDeleted, err := a.Redis.Client.Keys(REDISTAG_CATEGORIES_LIST + "*").Result()
	for i := range keysMustDeleted {
		err = a.Redis.RedisServices.Cache.Delete(keysMustDeleted[i])
		if err != nil {
			apiError := InternalError(err)
			a.renderJson(w, r, apiError.Code, apiError)
			return
		}
	}

	a.renderJson(w, r, http.StatusOK, Broadcast{Message: "Category deleted successfully"})
}
