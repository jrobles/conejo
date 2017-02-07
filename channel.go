package conejo

import (
	"github.com/streadway/amqp"
)

/*
Create a channel using an established conn [amqp.Connection]. Once the channel
is established, return the channel [amqp.Channel].
*/
func createChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	} else {
		return channel, nil
	}
}
