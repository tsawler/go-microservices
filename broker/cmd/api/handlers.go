package main

import (
	"broker/event"
	"encoding/json"
	"log"
	"net/http"
)

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
	Ok   bool   `json:"okay"`
}

func (app *Config) Home(w http.ResponseWriter, r *http.Request) {
	err := app.pushToQueue()
	if err != nil {
		log.Println(err)
	}
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Message = "Received request"

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}

func (app *Config) pushToQueue() error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		log.Println(err)
		return err
	}

	payload := Payload{
		Name: "my name",
		Data: "My data",
		Ok:   false,
	}

	j, _ := json.MarshalIndent(&payload, "", "    ")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
