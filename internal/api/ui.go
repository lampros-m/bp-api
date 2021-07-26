package api

import (
	"bestprice/bestprice-api/internal/config"
	"bestprice/bestprice-api/internal/data/cache"
	"bestprice/bestprice-api/internal/data/database"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	Router  *mux.Router
	Cfg     *config.Config
	Db      *database.DB
	Redis   *cache.RedisHanlder
	Address string
	ApiError
}

type Broadcast struct {
	Message     string
	Information string `json:",omitempty"`
}

func NewApi(db *database.DB, redisHandler *cache.RedisHanlder, conf *config.Config) *Api {
	api := Api{
		Db:      db,
		Redis:   redisHandler,
		Address: conf.ApiAddress,
		Router:  mux.NewRouter(),
	}
	api.setRoutesV1()
	return &api
}

func (a *Api) setRoutesV1() {
	v1 := a.Router.PathPrefix("/v1").Subrouter()
	v1.NotFoundHandler = http.HandlerFunc(a.notFound)

	v1.HandleFunc("/", a.ping).Methods("GET")

	v1.HandleFunc("/signup", a.Signup).Methods("POST")
	v1.HandleFunc("/login", a.Login).Methods("POST")

	v1.HandleFunc("/categories", a.listCategories).Methods("GET")
	v1.HandleFunc("/categories/{id}", a.readCategory).Methods("GET")
	v1.HandleFunc("/categories", Authenticator(a.createCategory, a)).Methods("POST")
	v1.HandleFunc("/categories/{id}", Authenticator(a.updateCategory, a)).Methods("PUT")
	v1.HandleFunc("/categories/{id}", Authenticator(a.deleteCategory, a)).Methods("DELETE")

	v1.HandleFunc("/products", a.listProducts).Methods("GET")
	v1.HandleFunc("/products/{id}", a.readProduct).Methods("GET")
	v1.HandleFunc("/products", Authenticator(a.createProduct, a)).Methods("POST")
	v1.HandleFunc("/products/{id}", Authenticator(a.updateProduct, a)).Methods("PUT")
	v1.HandleFunc("/products/{id}", Authenticator(a.deleteProduct, a)).Methods("DELETE")
}

func (a *Api) Run() {
	log.Println("Listening to", a.Address)
	log.Fatal(http.ListenAndServe(a.Address, a.Router))
}

func (a *Api) renderJson(w http.ResponseWriter, r *http.Request, statusCode int, input interface{}) {
	var response []byte
	var err error

	if inputStr, ok := input.(string); ok {
		response = []byte(inputStr)
	} else {
		response, err = json.Marshal(input)
		if err != nil {
			apiError := InternalError(err)
			statusCode = apiError.Code
			response, _ = json.Marshal(apiError)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func (a *Api) notFound(w http.ResponseWriter, r *http.Request) {
	err := errors.New("This page could not be found")
	apiError := NotFoundError(err)
	a.renderJson(w, r, apiError.Code, apiError)
	return
}
