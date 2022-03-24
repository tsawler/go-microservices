package main

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"os"
	"strconv"
)

const webPort = "80"

// Config is the application Config, shared with functions by using it as a receiver
type Config struct {
	Mailer Mail
	Etcd   *clientv3.Client
}

func main() {
	// create our configuration
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail-service on port", webPort)

	// define a server that listens on port 80 and uses our routes()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// connect to etcd and register service
	app.registerService()
	defer app.Etcd.Close()

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
