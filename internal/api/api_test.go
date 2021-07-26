package api

import (
	"bestprice/bestprice-api/internal/data/cache"
	"bestprice/bestprice-api/internal/data/database"
	"bestprice/bestprice-api/internal/data/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	Api Api
}

func (s *Suite) SetupTest() {
	s.Api = Api{
		Db:    database.NewDbMock(),
		Redis: cache.NewRedisHandlerMock(),
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

// ping
func (s *Suite) TestPing() {
	req, err := http.NewRequest("GET", "", nil)
	assert.Nil(s.T(), err)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(s.Api.ping)
	handler.ServeHTTP(recorder, req)
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
	expected := `{"Message":"I'm alive"}`
	assert.Equal(s.T(), expected, recorder.Body.String())
}

// categories
func (s *Suite) TestListCategories() {
	req, err := http.NewRequest("GET", "", nil)
	assert.Nil(s.T(), err)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(s.Api.listCategories)
	handler.ServeHTTP(recorder, req)
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
	bodyBytes := []byte(recorder.Body.String())
	var exposedCategories []model.CategoryExposed
	err = json.Unmarshal(bodyBytes, &exposedCategories)
	assert.Nil(s.T(), err)
	assert.Len(s.T(), exposedCategories, 2)
}

// products
func (s *Suite) TestListProducts() {
	req, err := http.NewRequest("GET", "", nil)
	assert.Nil(s.T(), err)
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(s.Api.listProducts)
	handler.ServeHTTP(recorder, req)
	assert.Equal(s.T(), http.StatusOK, recorder.Code)
	bodyBytes := []byte(recorder.Body.String())
	var exposedProducts []model.ProductExposed
	err = json.Unmarshal(bodyBytes, &exposedProducts)
	assert.Nil(s.T(), err)
	assert.Len(s.T(), exposedProducts, 4)
}
