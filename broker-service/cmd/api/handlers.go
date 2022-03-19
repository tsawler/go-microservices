package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Payload is the type for data we push into RabbitMQ
type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

// Broker is a simple test handler for the broker
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	err := app.pushToQueue("broker_hit", r.RemoteAddr)
	if err != nil {
		log.Println(err)
	}

	var payload jsonResponse
	payload.Message = "Received request"

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write(out)
}

// BrokerAuth is the handler to authenticate using the authentication-service.
// We receive user credentials as JSON, and then post that JSON to the authentication-service
// to try to authenticate.
func (app *Config) BrokerAuth(w http.ResponseWriter, r *http.Request) {
	// create a variable matching the structure of the JSON we expect from the front end
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// read posted json into our variable
	_ = readJSON(w, r, &requestPayload)

	// create json we'll send to the authentication-service
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")

	// call the authentication-service; we need a request, so let's build one, and populate
	// its body with the jsonData we just created. First we get the correct url for our
	// auth service from our service map.
	authServiceURL := fmt.Sprintf("http://%s/authenticate", app.GetServiceURL("auth"))

	// now build the request and set header
	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	// call the service
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode == http.StatusUnauthorized {
		_ = errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusAccepted {
		_ = errorJSON(w, errors.New("error calling auth service"), http.StatusBadRequest)
		return
	}

	// create variable we'll read the response.Body from the authentication-service into
	var jsonFromService jsonResponse

	// decode the json we get from the authentication-service into our variable
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if jsonFromService.Error {
		// invalid login
		_ = errorJSON(w, err, http.StatusUnauthorized)
		_ = app.pushToQueue("authentication", fmt.Sprintf("invalid login for %s", requestPayload.Email))
		return
	}

	// log action
	_ = app.pushToQueue("authentication", fmt.Sprintf("valid login for %s", requestPayload.Email))

	// send json back to our end user
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write(out)
}

// MailMessagePayload is the type for JSON describing a message to be sent
type MailMessagePayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// SendMailMessage sends a mail message which is received as JSON
func (app *Config) SendMailMessage(w http.ResponseWriter, r *http.Request) {
	var msg MailMessagePayload
	_ = readJSON(w, r, &msg)

	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail-service; we need a request, so let's build one, and populate
	// its body with the jsonData we just created. First we get the correct server
	// to call from our service map.
	mailServiceURL := fmt.Sprintf("http://%s/send", app.GetServiceURL("mail"))

	// now post to the mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		_ = errorJSON(w, errors.New("error calling mail service"), http.StatusBadRequest)
		return
	}

	// send json back to our end user
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write(out)

}

// HandleSubmission handles a JSON payload that describes an action to take,
// processes it, and sends it where it needs to go
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	// TODO - handle log, mail, auth,
}

// pushToQueue pushes a message into RabbitMQ
func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		log.Println(err)
		return err
	}

	payload := Payload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "    ")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
