package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

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
