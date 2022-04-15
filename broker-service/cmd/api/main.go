package main

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"os"
	"time"
)

// webPort the port that we listen on for api calls
const webPort = "80"

// Config is the type we'll use as a receiver to share application
// configuration around our app.
type Config struct {
	Rabbit         *amqp.Connection
	Etcd           *clientv3.Client
	LogServiceURLs map[string]string
	//MailServiceURLs map[string]string
	//AuthServiceURLs map[string]string
}

func main() {
	// don't continue until rabbitmq is ready
	rabbitConn, err := connectToRabbit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// don't continue until etcd is ready
	//etcConn, err := connectToEtcd()
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//defer etcConn.Close()

	app := Config{
		Rabbit: rabbitConn,
		//Etcd:   etcConn,
	}

	// get service urls
	//app.getServiceURLs()

	// watch service urls
	//go app.watchEtcd()

	log.Println("Starting broker service on port", webPort)

	// define the http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// connectToRabbit tries to connect to RabbitMQ, for up to 30 seconds
func connectToRabbit() (*amqp.Connection, error) {
	var rabbitConn *amqp.Connection
	var counts int64
	var rabbitURL = os.Getenv("RABBIT_URL")

	for {
		connection, err := amqp.Dial(rabbitURL)
		if err != nil {
			fmt.Println("rabbitmq not ready...")
			counts++
		} else {
			fmt.Println()
			rabbitConn = connection
			break
		}

		if counts > 15 {
			fmt.Println(err)
			return nil, errors.New("cannot connect to rabbit")
		}
		fmt.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
	fmt.Println("Connected to RabbitMQ!")
	return rabbitConn, nil
}
