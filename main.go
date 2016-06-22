package robimq

import (
	"github.com/streadway/amqp"
	"log"
)

func Connect(amqpURI string) (ch *amqp.Connection) {
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Printf("ERROR: Could not connect to RabbitMQ server - %q", err)
		return nil
	} else {
		log.Printf("Connected to RabbitMQ server - %s", amqpURI)
		return connection
	}
}

func declareExchange(ch *amqp.Channel, exchange string) error {
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

func declareQueue(ch *amqp.Channel, queue string) (q amqp.Queue, err error) {
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return q, err
}

func Consume(conn *amqp.Connection, exchange, queue, consumerTag string) (chan string, error) {

	data := make(chan string)

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("ERROR: Could not declare channel %q", err)
		return nil, err
	} else {

		log.Printf("Got Channel")
		defer ch.Close()

		err := declareExchange(ch, exchange)
		if err != nil {
			log.Printf("ERROR: Could not declare Exchange [%s] %q", exchange, err)
			return nil, err
		} else {

			q, err := declareQueue(ch, queue)
			if err != nil {
				log.Printf("ERROR: Could not declare queue [%s] %q", queue, err)
				return nil, err
			} else {

				err = ch.QueueBind(
					q.Name,   // queue name
					"image",  // routing key @TODO
					exchange, // exchange
					false,
					nil,
				)
				if err != nil {
					log.Printf("ERROR: Could not bind [%s] queue to [%s] exhange %q", queue, exchange, err)
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
