package conejo

import (
	"github.com/streadway/amqp"
)

func Connect(amqpURI string) *amqp.Connection {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil
	} else {
		return conn
	}
}
