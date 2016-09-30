package conejo

import (
	"github.com/streadway/amqp"
)

type exchange struct {
	Name        string
	Type        string
	Durable     bool
	AutoDeleted bool
	Internal    bool
	NoWait      bool
	Arguments   amqp.Table
}

func declareExchange(exchange string) error {
	err := ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // noWait
		nil,      // arguments
	)
	return err
}

func _declareExchange(e exchange) error {
	err := ch.ExchangeDeclare(
		e.Name,        // name
		e.Type,        // type
		e.Durable,     // durable
		e.AutoDeleted, // auto-deleted
		e.Internal,    // internal
		e.NoWait,      // noWait
		e.Arguments,   // arguments
	)
	return err
}
