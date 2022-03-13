package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tsawler/go-rabbit/lib/event"
)

func main() {
	var rabbitConn *amqp.Connection
	var counts int64
	var backOff = 1 * time.Second

	// don't continue until rabbitmq is ready
	for {
		connection, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not ready...")
			counts++
		} else {
			fmt.Println()
			rabbitConn = connection
			break
		}

		if counts > 5 {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Backing off for %d seconds...\n", int(math.Pow(float64(counts), 2)))
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)
		continue

	}

	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for RabbitMQ messages...")
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}
	err = consumer.Listen(os.Args[1:])
	if err != nil {
		log.Println(err)
	}
}
