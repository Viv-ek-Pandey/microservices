package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestBody)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	//validate user
	user, err := app.Models.User.GetByEmail(requestBody.Email)
	if err != nil {
		app.errorJson(w, errors.New("invalid credential"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestBody.Password)
	if err != nil || !valid {
		app.errorJson(w, errors.New("invalid credential"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}
	app.writeJson(w, http.StatusAccepted, payload)
}
