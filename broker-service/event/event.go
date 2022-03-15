package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func getExchangeName() string {
	return "logs_topic"
}

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		getExchangeName(), // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
}
