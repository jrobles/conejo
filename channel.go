package conejo

import (
	"github.com/streadway/amqp"
)

func createChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	} else {
		return channel, nil
	}
}
