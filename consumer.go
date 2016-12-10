package conejo

import (
	"github.com/streadway/amqp"
	"log"
)

func Consume(conn *amqp.Connection, exchange Exchange, queue Queue, consumerTag string) (chan string, error) {

	data := make(chan string)

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("ERROR: Could not declare channel %q", err)
		return nil, err
	} else {

		log.Printf("Got Channel")
		defer ch.Close()

		err := declareExchange(exchange, ch)
		if err != nil {
			log.Printf("ERROR: Could not declare Exchange [%s] %q", exchange.Name, err)
			return nil, err
		} else {

			q, err := declareQueue(queue, ch)
			if err != nil {
				log.Printf("ERROR: Could not declare queue [%s] %q", queue.Name, err)
				return nil, err
			} else {

				err = ch.QueueBind(
					q.Name,        // queue name
					"image",       // routing key @TODO
					exchange.Name, // exchange
					false,
					nil,
				)
				if err != nil {
					log.Printf("ERROR: Could not bind [%s] queue to [%s] exhange %q", queue, exchange.Name, err)
					return nil, err
				} else {

					log.Printf("Queue %s declared", queue)
					err = ch.Qos(
						1,     // prefetch count
						0,     // prefetch size
						false, // global
					)
					if err != nil {
						log.Printf("ERROR: %q", err)
						return nil, err
					}
					msgs, err := ch.Consume(
						q.Name,      // queue
						consumerTag, // consumer
						false,       // auto-ack
						false,       // exclusive
						false,       // no-local
						false,       // no-wait
						nil,         // args
					)
					if err != nil {
						log.Printf("ERROR: Could not consume messages on [%s] queue %q", queue, err)
						return nil, err
					}

					forever := make(chan bool)

					go func() {
						for d := range msgs {
							d.Ack(false)
							//downloadFile(string(d.Body))
							data <- string(d.Body)
						}
					}()
					log.Printf("Consumer tag %s", consumerTag)
					<-forever
					return data, nil
				}

			}

		}

	}
	return data, nil
}
