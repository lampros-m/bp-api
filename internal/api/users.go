package api

import (
	"bestprice/bestprice-api/internal/config"
	"bestprice/bestprice-api/internal/data/model"
	"bestprice/bestprice-api/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (a *Api) Signup(w http.ResponseWriter, r *http.Request) {
	body, err := helper.ReadHttpRequestAndClose(r)
	if err != nil {
		apiError := BadRequestError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	var credentials model.Credentials

	err = json.Unmarshal(body, &credentials)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 8)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}
	credentials.Password = string(hashedPassword)

	err = a.Db.Users.AddUser(credentials)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	broadcast := Broadcast{
		Message: "User created successfully",
	}

	a.renderJson(w, r, http.StatusOK, broadcast)
}

func (a *Api) Login(w http.ResponseWriter, r *http.Request) {
	body, err := helper.ReadHttpRequestAndClose(r)
	if err != nil {
		apiError := BadRequestError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	var credentials model.Credentials

	err = json.Unmarshal(body, &credentials)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	user, err := a.Db.Users.IdentifyUser(credentials)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	config := config.NewConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtKey)
	if err != nil {
		apiError := InternalError(err)
		a.renderJson(w, r, apiError.Code, apiError)
		return
	}

	broadcast := Broadcast{
		Message:     fmt.Sprintf("Welcome %s", user.Username),
		Information: tokenString,
	}

	a.renderJson(w, r, http.StatusOK, broadcast)
}
