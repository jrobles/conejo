package conejo

import (
	"github.com/streadway/amqp"
	"log"
)

var (
	ch *amqp.Channel
)

func Connect(amqpURI string) (conn *amqp.Connection) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Printf("ERROR: Could not connect to RabbitMQ server - %q", err)
		return nil
	} else {
		log.Printf("[CONEJO] Connected to RabbitMQ server - %s", amqpURI)
		ch, err = conn.Channel()
		if err != nil {
			log.Printf("ERROR: Could not declare channel %q", err)
			return nil
		} else {
			return conn
		}
	}
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

func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")
	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}

func Publish(conn *amqp.Connection, queue, exchange string) {

	var reliable = true
	var body = "DSDH HJD SHDJ DHJKLD HASJDKL HJKLDSA"

	err := declareExchange(exchange)
	if err != nil {
		log.Printf("ERROR: Could not declare exchange %q", err)
	} else {
		log.Printf("[CONEJO] Declared exchange")
		q, err := declareQueue(queue)
		if err != nil {
			log.Printf("ERROR: Could not declare queue %q", err)
		} else {
			log.Printf("[CONEJO] Declared queue")
			err = ch.QueueBind(
				q.Name, // queue name
				queue,
				exchange, // exchange
				false,
				nil,
			)
			if err != nil {
				log.Printf("ERROR: Could not bind queue '%s' to exchange '%s' using '%s'", queue, exchange, queue, err)
			} else {

				// Reliable publisher confirms require confirm.select support from to connection.
				if reliable {
					if err := ch.Confirm(false); err != nil {
						log.Printf("ERROR: Could not confirm - %s", err)
					}

					confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

					defer confirmOne(confirms)
				}

				if err = ch.Publish(
					exchange, // publish to an exchange
					q.Name,   // routing to 0 or more queues
					false,    // mandatory
					false,    // immediate
					amqp.Publishing{
						Headers:         amqp.Table{},
						ContentType:     "text/plain",
						ContentEncoding: "",
						Body:            []byte(body),
						DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
						Priority:        0,              // 0-9
						// a bunch of application/implementation-specific fields
					},
				); err != nil {
					log.Printf("ERROR: Could not publish message %s", err)
				} else {
					ch.Close()
				}

			}
		}
	}
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

		err := declareExchange(exchange)
		if err != nil {
			log.Printf("ERROR: Could not declare Exchange [%s] %q", exchange, err)
			return nil, err
		} else {

			q, err := declareQueue(queue)
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
