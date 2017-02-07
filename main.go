package conejo

import (
	"github.com/streadway/amqp"
)

/*
Connect to the RabbitMQ Sever using the server URI [string].
If a connections is established, return the connection [amqp.Connection].
*/
func Connect(amqpURI string) *amqp.Connection {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil
	} else {
		return conn
	}
}
