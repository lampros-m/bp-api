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
	REDISTAG_PRODUCT       = "product.id:"
	REDISTAG_PRODUCTS_LIST = "product.list:"
)

func (a *Api) listProducts(w http.ResponseWriter, r *http.Request) {
	var products []model.Product
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

	redisTagKey := fmt.Sprintf("%ssort[%s].limit[%s].offset[%s]", REDISTAG_PRODUCTS_LIST, sortParameter, limitParameter, offsetParameter)
	flagRedisKeyFound := false

	value, err := a.Redis.RedisServices.Cache.Get(redisTagKey)
	if err == nil {
		flagRedisKeyFound = true
	}

	if flagRedisKeyFound {
		a.renderJson(w, r, http.StatusOK, value)
		return
	}

	products, err = a.Db.Products.GetProducts(sortingSql, limitSql, offsetSql)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	productsRedisValue, err := json.Marshal(helper.ExposeProducts(products))
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	err = a.Redis.RedisServices.Cache.Set(redisTagKey, string(productsRedisValue))
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	a.renderJson(w, r, http.StatusOK, helper.ExposeProducts(products))
}

func (a *Api) readProduct(w http.ResponseWriter, r *http.Request) {
	productId := mux.Vars(r)["id"]

	if !helper.SupportedSqlString(productId) {
		apiError := BadRequestError(errors.New("Not supported characters in id provided"))
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	redisTagKey := REDISTAG_PRODUCT + productId
	flagRedisKeyFound := false

	value, err := a.Redis.RedisServices.Cache.Get(redisTagKey)
	if err == nil {
		flagRedisKeyFound = true
	}

	if flagRedisKeyFound {
		a.renderJson(w, r, http.StatusOK, value)
		return
	}

	product, err := a.Db.Products.GetProduct(productId)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	productRedisValue, err := json.Marshal(func(p model.ProductExposed) *model.ProductExposed { return &p }(model.ProductExposed(product)))
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	err = a.Redis.RedisServices.Cache.Set(redisTagKey, string(productRedisValue))
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	a.renderJson(w, r, http.StatusOK, func(p model.ProductExposed) *model.ProductExposed { return &p }(model.ProductExposed(product)))
}

func (a *Api) createProduct(w http.ResponseWriter, r *http.Request) {
	body, err := helper.ReadHttpRequestAndClose(r)
	if err != nil {
		apiError := BadRequestError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	var product model.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	err = a.Db.Products.AddProduct(product)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	a.renderJson(w, r, http.StatusOK, Broadcast{Message: "Product added"})
}

func (a *Api) updateProduct(w http.ResponseWriter, r *http.Request) {
	productId := mux.Vars(r)["id"]

	if !helper.SupportedSqlString(productId) {
		apiError := BadRequestError(errors.New("Not supported characters in id provided"))
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	redisTagKey := REDISTAG_PRODUCT + productId
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

	var product model.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	err = a.Db.Products.UpdateProduct(productId, product)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	keysMustDeleted, err := a.Redis.Client.Keys(REDISTAG_PRODUCTS_LIST + "*").Result()
	for i := range keysMustDeleted {
		err = a.Redis.RedisServices.Cache.Delete(keysMustDeleted[i])
		if err != nil {
			apiError := InternalError(err)
			a.renderJson(w, r, apiError.Code, apiError)
			return
		}
	}

	a.renderJson(w, r, http.StatusOK, Broadcast{Message: "Product updated successfully"})
}

func (a *Api) deleteProduct(w http.ResponseWriter, r *http.Request) {
	productId := mux.Vars(r)["id"]

	if !helper.SupportedSqlString(productId) {
		apiError := BadRequestError(errors.New("Not supported characters in id provided"))
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	redisTagKey := REDISTAG_PRODUCT + productId
	err := a.Redis.RedisServices.Cache.Delete(redisTagKey)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	err = a.Db.Products.DeleteProduct(productId)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	keysMustDeleted, err := a.Redis.Client.Keys(REDISTAG_PRODUCTS_LIST + "*").Result()
	for i := range keysMustDeleted {
		err = a.Redis.RedisServices.Cache.Delete(keysMustDeleted[i])
		if err != nil {
			apiError := InternalError(err)
			a.renderJson(w, r, apiError.Code, apiError)
			return
		}
	}

	a.renderJson(w, r, http.StatusOK, Broadcast{Message: "Product deleted successfully"})
}
