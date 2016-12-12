package conejo

import (
	"github.com/streadway/amqp"
)

func Consume(conn *amqp.Connection, queue Queue, exchange Exchange, consumerTag string, cb chan string) error {

	channel, err := createChannel(conn)
	if err != nil {
		return err
	}
	defer channel.Close()

	err = declareExchange(exchange, channel)
	if err != nil {
		return err
	} else {

		err = declareQueue(queue, channel)
		if err != nil {
			return err
		} else {

			err = channel.QueueBind(
				queue.Name,    // queue name
				queue.Name,    // routing key @TODO
				exchange.Name, // exchange
				false,
				nil,
			)
			if err != nil {
				return err
			} else {

				err = channel.Qos(
					1,     // prefetch count
					0,     // prefetch size
					false, // global
				)
				if err != nil {
					return err
				}

				msgs, err := channel.Consume(
					queue.Name,  // queue
					consumerTag, // consumer
					false,       // auto-ack
					false,       // exclusive
					false,       // no-local
					false,       // no-wait
					nil,         // args
				)
				if err != nil {
					return err
				}

				forever := make(chan bool)
				go func() {
					for d := range msgs {
						d.Ack(false)
						cb <- string(d.Body)
					}
				}()
				<-forever
				return nil

			} // Bound to queue

		} // Queue declared

	}
	return nil
}
