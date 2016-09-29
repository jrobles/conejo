package conejo

import (
	"github.com/streadway/amqp"
)

func declareQueue(queue string) (q amqp.Queue, err error) {
	q, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return q, err
}
