package services

import (
	"bestprice/bestprice-api/internal/data/model"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UsersServices struct {
	DB         *sqlx.DB
	SqlQueries map[string]string
}

func NewUsersServices(db *sqlx.DB) *UsersServices {
	return &UsersServices{
		DB: db,
		SqlQueries: map[string]string{
			"insert_user": `INSERT INTO users (id, username, pass) 
							VALUES (UUID(), ?, ?);`,

			"select_user_by_username": `SELECT *
										FROM users
										WHERE 1=1
										AND username=?
										LIMIT 1;`,
		},
	}
}

func (o *UsersServices) AddUser(credentials model.Credentials) error {
	_, err := o.DB.Exec(o.SqlQueries["insert_user"], credentials.Username, credentials.Password)
	if err != nil {
		return err
	}

	return nil
}

func (o *UsersServices) IdentifyUser(credentials model.Credentials) (model.User, error) {
	var user, outUser model.User
	var err error

	rows, err := o.DB.Queryx(o.SqlQueries["select_user_by_username"], credentials.Username)
	if err != nil {
		return outUser, err
	}

	if !rows.Next() {
		return outUser, errors.New("Could not find the provided username")
	}

	err = rows.StructScan(&user)
	if err != nil {
		return outUser, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		return outUser, errors.New("User could not be indentified")
	}

	outUser = user
	return user, nil
}
