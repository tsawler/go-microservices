package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"time"
)

// authServiceURL is the url to the authentication service. Since we're using
// Docker, we specify the appropriate entry from docker-compose.yml
const authServiceURL = "http://authentication-service/authenticate"

// JSONPayload is the type for JSON posted to this API
type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// WriteLog is the handler to accept a post request consisting of json payload,
// and then write it to Mongo
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// insert the data
	err := app.logEvent(requestPayload.Name, requestPayload.Data)
	if err != nil {
		log.Println(err)
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// create the response we'll send back as JSON
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	// write the response back as JSON
	_ = app.writeJSON(w, http.StatusAccepted, resp)
}

// Logout logs the user out and redirects them to the login page
func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// log the event
	_ = app.logEvent("authentication", fmt.Sprintf("%s logged out of the logger service", app.Session.GetString(r.Context(), "email")))

	// clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	// redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// LoginPage displays the login page
func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

// LoginPagePost handles user login. Note that it calls the authentication microservice
func (app *Config) LoginPagePost(w http.ResponseWriter, r *http.Request) {
	// it's always good to regenerate the session on login/logout
	_ = app.Session.RenewToken(r.Context())

	// parse the posted form data
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// get email and password from form post
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// create a variable we'll post to the auth service as JSON
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	requestPayload.Email = email
	requestPayload.Password = password

	// create json we'll send to the authentication-service
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")

	// call the authentication-service; we need a request, so let's build one, and populate
	// its body with the jsonData we just created
	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	c := &http.Client{}
	response, err := c.Do(request)
	if err != nil {
		log.Println(err)
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode == http.StatusUnauthorized {
		log.Println("wrong status code", response.StatusCode)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else if response.StatusCode != http.StatusAccepted {
		log.Println("did not get status accepted")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Read the body of the response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// define a type that matches the JSON we're getting from the response body
	type userPayload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Data    struct {
			ID        int       `json:"id"`
			Email     string    `json:"email"`
			FirstName string    `json:"first_name"`
			LastName  string    `json:"last_name"`
			Active    int       `json:"active"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"data"`
	}

	// declare a variable we can unmarshal the JSON into
	var user userPayload

	// since we received the request from a remote host with http.Client,
	// we need to build a new request with the body we received and pass it to
	// app.readJSON
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))

	// read the JSON
	err = app.readJSON(w, req, &user)
	if err != nil {
		log.Println(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// log the event
	_ = app.logEvent("authentication", fmt.Sprintf("%s logged into the logger service", user.Data.Email))

	// set up session & log user in
	app.Session.Put(r.Context(), "userID", user.Data.ID)
	app.Session.Put(r.Context(), "email", user.Data.Email)

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

// Dashboard displays the dashboard page
func (app *Config) Dashboard(w http.ResponseWriter, r *http.Request) {
	// get the list of all log entries from mongo
	logs, err := app.Models.LogEntry.All()
	if err != nil {
		log.Println("Error getting all log entries")
		app.clientError(w, http.StatusBadRequest)
	}

	templateData := make(map[string]any)
	templateData["logs"] = logs

	app.render(w, r, "dashboard.page.gohtml", &TemplateData{
		Data: templateData,
	})
}

// DisplayOne is the handler to display a single log entry
func (app *Config) DisplayOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	entry, err := app.Models.LogEntry.GetOne(id)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	templateData := make(map[string]any)
	templateData["entry"] = entry

	app.render(w, r, "entry.page.gohtml", &TemplateData{
		Data: templateData,
	})
}

// DeleteAll drops everything in the collection and redirects to the same page
func (app *Config) DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := app.Models.LogEntry.DropCollection()
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

// UpdateTimeStamp just demos how to update a document
func (app *Config) UpdateTimeStamp(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	entry, err := app.Models.LogEntry.GetOne(id)
	if err != nil {
		log.Println("Error getting record:", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	res, err := entry.Update()
	if err != nil {
		log.Println("Error updating", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	log.Println("Result in handler:", res.ModifiedCount)

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
