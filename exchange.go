package conejo

import (
	"github.com/streadway/amqp"
)

type Exchange struct {
	Name        string
	Type        string
	Durable     bool
	AutoDeleted bool
	Internal    bool
	NoWait      bool
	Arguments   amqp.Table
}

func declareExchange(e Exchange) error {
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
