package conejo

import (
	"github.com/streadway/amqp"
	"log"
)

func Connect(amqpURI string) (conn *amqp.Connection) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Printf("[CONEJO] Could not connect to RabbitMQ server - %q", err)
		return nil
	} else {
		log.Printf("[CONEJO] Connected to RabbitMQ server - %s", amqpURI)
		return conn
	}
}
