package api

import "net/http"

func (a *Api) ping(w http.ResponseWriter, r *http.Request) {
	message := Broadcast{
		Message: "I'm alive",
	}
	a.renderJson(w, r, http.StatusOK, message)
}
