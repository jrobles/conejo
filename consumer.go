package conejo

import (
	"github.com/streadway/amqp"
	"log"
)

func Consume(conn *amqp.Connection, queue Queue, exchange Exchange, consumerTag string, cb chan string) error {

	channel, err := conn.Channel()
	if err != nil {
		log.Printf("[CONEJO] Could not declare channel %q", err)
	}
	defer channel.Close()

	err = declareExchange(exchange, channel)
	if err != nil {
		log.Printf("ERROR: Could not declare Exchange [%s] %q", exchange.Name, err)
		return err
	} else {

		err = declareQueue(queue, channel)
		if err != nil {
			log.Printf("ERROR: Could not declare queue [%s] %q", queue.Name, err)
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
				log.Printf("ERROR: Could not bind [%s] queue to [%s] exhange %q", queue.Name, exchange.Name, err)
				return err
			} else {

				log.Printf("Queue %s declared", queue.Name)
				err = channel.Qos(
					1,     // prefetch count
					0,     // prefetch size
					false, // global
				)
				if err != nil {
					log.Printf("ERROR: %q", err)
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
					log.Printf("ERROR: Could not consume messages on [%s] queue %q", queue.Name, err)
					return err
				}

				forever := make(chan bool)
				go func() {
					for d := range msgs {
						d.Ack(false)
						cb <- string(d.Body)
					}
				}()
				log.Printf("Consumer tag %s", consumerTag)
				<-forever
				return nil

			} // Bound to queue

		} // Queue declared

	}
	return nil
}
