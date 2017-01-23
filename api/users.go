package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mainflux/go-auth/domain"
	"github.com/mainflux/go-auth/services"
)

type userReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &userReq{}
	if err = json.Unmarshal(body, data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := services.RegisterUser(data.Username, data.Password)
	if err != nil {
		authErr := err.(*domain.AuthError)
		w.WriteHeader(authErr.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &userReq{}
	if err = json.Unmarshal(body, data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := services.Login(data.Username, data.Password)
	if err != nil {
		authErr := err.(*domain.AuthError)
		w.WriteHeader(authErr.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
