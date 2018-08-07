package conejo

import (
	"github.com/streadway/amqp"
	"log"
)

/*
Create a channel using an established conn [amqp.Connection]. Once the channel
is established, return the channel [amqp.Channel].
*/
func createChannel(conn *amqp.Connection) (*amqp.Channel, error) {

	// Create the channel
	channel, err := conn.Channel()

	// Check for errors
	if err != nil {

		// Could not create channel
		log.Printf("CONEJO: Could not create channel - %q", err)
		return nil, err

	} else {

		return channel, nil
	}
}

func CloseChannel(channel *amqp.Channel) {

	channel.Close()

}
