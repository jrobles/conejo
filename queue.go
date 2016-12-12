package conejo

import (
	"github.com/streadway/amqp"
)

type Queue struct {
	Name      string
	Durable   bool
	Delete    bool
	Exclusive bool
	NoWait    bool
	Arguments amqp.Table
}

func declareQueue(queue Queue, channel *amqp.Channel) error {
	_, err := channel.QueueDeclare(
		queue.Name,      // name
		queue.Durable,   // durable
		queue.Delete,    // delete when unused
		queue.Exclusive, // exclusive
		queue.NoWait,    // no-wait
		queue.Arguments, // arguments
	)
	return err
}
