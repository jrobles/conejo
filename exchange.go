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

func declareExchange(exchange Exchange, channel *amqp.Channel) error {
	err := channel.ExchangeDeclare(
		exchange.Name,        // name
		exchange.Type,        // type
		exchange.Durable,     // durable
		exchange.AutoDeleted, // auto-deleted
		exchange.Internal,    // internal
		exchange.NoWait,      // noWait
		exchange.Arguments,   // arguments
	)
	return err
}
