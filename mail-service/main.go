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
		Host:        "mailhog",
		Port:        1025,
		Username:    "",
		Password:    "",
		Encryption:  "none",
		FromName:    "John Smith",
		FromAddress: "john.smith@example.com",
	}

	return s
}
