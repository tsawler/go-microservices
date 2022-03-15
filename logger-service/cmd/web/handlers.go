package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const authServiceURL = "http://authentication-service/authenticate"

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	log.Println("Received request")
	var requestPayload JSONPayload
	_ = readJSON(w, r, &requestPayload)

	collection := app.Mongo.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      requestPayload.Name,
		Data:      requestPayload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println(err)
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Message = "logged"
	_ = writeJSON(w, http.StatusAccepted, resp)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	render(w, "login.page.gohtml")
}

func (app *Config) LoginPagePost(w http.ResponseWriter, r *http.Request) {
	// authentication-service
	_ = app.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	requestPayload.Email = email
	requestPayload.Password = password

	// create json we'll send to the authentication-service
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")

	// call the authentication-service
	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode == http.StatusUnauthorized {
		log.Println("wrong status code")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if response.StatusCode != http.StatusAccepted {
		log.Println("did not get status accepted")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user struct {
		ID        int       `json:"id"`
		Email     string    `json:"email"`
		FirstName string    `json:"first_name,omitempty"`
		LastName  string    `json:"last_name,omitempty"`
		Password  string    `json:"-"`
		Active    int       `json:"active"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	err = readJSON(w, r, &user)
	if err != nil {
		log.Println("error reading json from auth:", err)
	}

	// set up session & log user in
	app.Session.Put(r.Context(), "userID", user.ID)

	w.Write([]byte("Auth worked!"))
}
