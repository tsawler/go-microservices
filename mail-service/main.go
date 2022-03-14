package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Config is the application Config, shared with functions by using it as a receiver
type Config struct {
	Mailer Mail
}

func main() {
	// create our configuration
	app := Config{
		Mailer: createMail(),
	}

	// define a server that listens on port 80 and uses our routes()
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

// createMail creates a variable of type Mail and populates its values.
// Typically, this kind of information comes from the environment, or from
// command line parameters.
func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	s := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return s
}
