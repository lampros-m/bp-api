package database

import (
	"bestprice/bestprice-api/internal/config"
	"bestprice/bestprice-api/internal/data/model"
	"bestprice/bestprice-api/internal/data/services"

	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DB struct {
	Connection *sqlx.DB
	Categories ICategoriesServicesDB
	Products   IProductsServicesDB
	Users      IUsersServicesDB
}

func NewDb(conf *config.Config) (*DB, error) {
	conn, err := sqlx.Connect("mysql", conf.DbPath)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot connect to db")
	}

	conn.SetMaxOpenConns(conf.DbMaxConnections)
	conn.SetMaxIdleConns(conf.DbMaxConnections)
	conn.SetConnMaxLifetime(300 * time.Second)

	db := DB{
		Connection: conn,
		Categories: services.NewCategoriesServices(conn),
		Products:   services.NewProductsServices(conn),
		Users:      services.NewUsersServices(conn),
	}

	return &db, nil
}

type ICategoriesServicesDB interface {
	GetCategories(string, string, string) ([]model.Category, error)
	GetCategory(string) (model.Category, error)
	AddCategory(model.Category) error
	UpdateCategory(string, model.Category) error
	DeleteCategory(string) error
}

type IProductsServicesDB interface {
	GetProducts(string, string, string) ([]model.Product, error)
	GetProduct(string) (model.Product, error)
	AddProduct(model.Product) error
	UpdateProduct(string, model.Product) error
	DeleteProduct(string) error
}

type IUsersServicesDB interface {
	AddUser(model.Credentials) error
	IdentifyUser(model.Credentials) (model.User, error)
}
