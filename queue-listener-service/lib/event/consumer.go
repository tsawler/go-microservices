package event

import (
	"encoding/json"
	"fmt"
	"log"
	"net/rpc"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer for receiving AMPQ events
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

// NewConsumer returns a new Consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Listen will listen for all new Queue publications
// and print them to the console.
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err = ch.QueueBind(
			q.Name,
			s,
			getExchangeName(),
			false,
			nil,
		)
		fmt.Println("Exchange name:", getExchangeName())

		if err != nil {
			log.Println(err)
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			log.Println("Got value from json:", payload.Name)
			// Do something with the payload
			go handlePayload(payload)
		}
	}()

	log.Printf("[*] Waiting for message [Exchange, Queue][%s, %s]. To exit press CTRL+C", getExchangeName(), q.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload) {
	// logic to process payload goes in here
	switch payload.Name {
	case "broker_hit":
		res, err := rpcPush("LogInfo", payload.Data)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Response from RPC:", res)

	default:

	}
}

func rpcPush(function string, data string) (string, error) {
	c, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		log.Println(err)
		return "", err
	}

	fmt.Println("Connected via rpc...")
	var result string
	err = c.Call("RPCServer."+function, data, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}
