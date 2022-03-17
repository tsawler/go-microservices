package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	var rabbitConn *amqp.Connection
	var counts int64
	log.Println("--------------------------")
	log.Println("Starting broker-service...")

	// don't continue until rabbitmq is ready
	for {
		connection, err := amqp.Dial("amqp://guest:guest@rabbitmq")
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
			os.Exit(1)
		}
		fmt.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}

	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Println("Starting broker service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	srv.ListenAndServe()
}
