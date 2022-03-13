package main

import (
	"fmt"
	"log"
	"net/http"
)

type config struct {
	Mailer Mail
}

func main() {
	app := config{
		Mailer: createMail(),
	}

	go app.Mailer.ListenForMail()

	srv := &http.Server{
		Addr:    ":80",
		Handler: app.routes(),
	}
	fmt.Println("Starting mail service on port 80")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func createMail() Mail {
	s := Mail{
		Domain:      "localhost",
		Templates:   "./templates",
		Host:        "mailhog",
		Port:        1025,
		Username:    "",
		Password:    "",
		Encryption:  "none",
		FromName:    "John Smith",
		FromAddress: "john.smith@example.com",
		Jobs:        make(chan Message, 5),
		Results:     make(chan Result, 5),
		API:         "",
		APIKey:      "",
		APIUrl:      "",
	}

	return s
}
