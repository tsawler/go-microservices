package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	var rabbitConn *amqp.Connection
	var counts int64

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

	srv := &http.Server{
		Addr:    ":80",
		Handler: app.routes(),
	}

	srv.ListenAndServe()
}
