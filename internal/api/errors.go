package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorJsonMarshable struct {
	error
}

type ApiError struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Err     ErrorJsonMarshable `json:"error"`
}

func (o ErrorJsonMarshable) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Error())
}

func InternalError(err error) ApiError {
	log.Println(err)
	code := http.StatusInternalServerError

	return ApiError{
		code,
		http.StatusText(code),
		ErrorJsonMarshable{err},
	}
}

func UnauthorizedError(err error) ApiError {
	log.Println(err)
	code := http.StatusUnauthorized

	return ApiError{
		code,
		http.StatusText(code),
		ErrorJsonMarshable{err},
	}
}

func BadRequestError(err error) ApiError {
	log.Println(err)
	code := http.StatusBadRequest

	return ApiError{
		code,
		http.StatusText(code),
		ErrorJsonMarshable{err},
	}
}

func NotFoundError(err error) ApiError {
	log.Println(err)
	code := http.StatusNotFound

	return ApiError{
		code,
		http.StatusText(code),
		ErrorJsonMarshable{err},
	}
}
