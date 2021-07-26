package api

import (
	"bestprice/bestprice-api/internal/config"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func ReadJWTToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Invalid authorization request")
	}
	return parts[1], nil
}

func Authenticator(next http.HandlerFunc, a *Api) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := ReadJWTToken(r)
		if err != nil {
			apiError := UnauthorizedError(err)
			a.renderJson(w, r, apiError.Code, apiError)
			return
		}

		if tokenStr == "" {
			err = errors.New("Jwt token is empty")
			apiError := UnauthorizedError(err)
			a.renderJson(w, r, apiError.Code, apiError)
			return
		}

		config := config.NewConfig()
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})
		if err != nil {
			apiError := InternalError(err)
			a.renderJson(w, r, apiError.Code, apiError)
			return
		}

		if !tkn.Valid {
			err = errors.New("The token provided is not valid or expired")
			apiError := UnauthorizedError(err)
			a.renderJson(w, r, apiError.Code, apiError)
			return
		}

		next(w, r)
	}
}
