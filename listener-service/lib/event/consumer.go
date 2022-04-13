package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer is the type used for receiving AMPQ events
type Consumer struct {
	conn      *amqp.Connection
	queueName string
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

// setup opens a channel and declares the exchange
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

// Payload is the type used for pushing events to RabbitMQ
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Listen will listen for all new queue publications
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
			// get the JSON payload and unmarshal it into a variable
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			// do something with the payload
			go handlePayload(payload)
		}
	}()

	log.Printf("[*] Waiting for message [Exchange, Queue][%s, %s].", getExchangeName(), q.Name)
	<-forever
	return nil
}

// handlePayload takes an action based on the name of an event in the queue
func handlePayload(payload Payload) {
	// logic to process payload goes in here
	switch payload.Name {
	case "broker_hit":
		// just a test to make sure everything works
		res, err := rpcPushToLogger("LogInfo", payload)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Response from RPC:", res)

	case "auth", "authentication":
		// we are trying to authenticate someone
		err := authenticate(payload)
		if err != nil {
			log.Println(err)
		}

	// you can have as many cases as you want here, but naturally you'll have to write the logic
	// to connect to a given microservice

	default:
		// log whatever we get
		res, err := rpcPushToLogger("LogInfo", payload)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Response from RPC:", res)
	}
}

// rpcPushToLogger pushes data to the logger-service via RPC, where
// it gets stored into a mongo database
func rpcPushToLogger(function string, data any) (string, error) {
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

// authenticate is a stub that we'll never actually use, but it is here
// as we get used to how to interact with services
func authenticate(payload Payload) error {
	// TODO actually authenticate via JSON
	log.Printf("Got payload of %v", payload)
	return nil
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
