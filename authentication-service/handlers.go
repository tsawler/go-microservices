package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := readJSON(w, r, &requestPayload)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate against database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		_ = errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		_ = errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", requestPayload.Email),
		//Data: User{
		//	ID:        1,
		//	FirstName: "Jack",
		//	LastName:  "Smith",
		//	Email:     "jack@smith.com",
		//	Active:    1,
		//},
		Data: user,
	}

	_ = writeJSON(w, http.StatusAccepted, payload)
}
