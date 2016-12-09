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

func declareQueue(q Queue) (aq amqp.Queue, err error) {
	aq, err = ch.QueueDeclare(
		q.Name,      // name
		q.Durable,   // durable
		q.Delete,    // delete when unused
		q.Exclusive, // exclusive
		q.NoWait,    // no-wait
		q.Arguments, // arguments
	)
	return aq, err
}
