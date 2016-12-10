package conejo

import (
	"github.com/streadway/amqp"
	"log"
)

func Connect(amqpURI string) (conn *amqp.Connection) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Printf("[CONEJO] Could not connect to RabbitMQ server - %q", err)
		return nil
	} else {
		log.Printf("[CONEJO] Connected to RabbitMQ server - %s", amqpURI)
		return conn
	}
}

func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")
	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}

func Publish(conn *amqp.Connection, queue Queue, exchange Exchange, body string) {

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[CONEJO] Could not declare channel %q", err)
	}

	var reliable = true

	err = declareExchange(exchange, ch)
	if err != nil {
		log.Printf("ERROR: Could not declare exchange %q", err)
	} else {
		log.Printf("[CONEJO] Declared exchange")
		q, err := declareQueue(queue, ch)
		if err != nil {
			log.Printf("ERROR: Could not declare queue %q", err)
		} else {
			log.Printf("[CONEJO] Declared queue")
			err = ch.QueueBind(
				queue.Name,    // queue name
				queue.Name,    // @TODO - FIX ME!!!
				exchange.Name, // exchange
				false,
				nil,
			)
			if err != nil {
				log.Printf("ERROR: Could not bind queue '%s' to exchange '%s' using '%s' - ", queue.Name, exchange.Name, queue.Name, err)
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
					exchange.Name, // publish to an exchange
					q.Name,        // routing to 0 or more queues
					false,         // mandatory
					false,         // immediate
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
					//ch.Close()
					defer ch.Close()
				}

			}
		}
	}
}

/*
func Consume(conn *amqp.Connection, exchangeName string, queue Queue, consumerTag string) (chan string, error) {

	data := make(chan string)
	c := Exchange{
		Name:        exchangeName,
		Type:        "topic",
		Durable:     true,
		AutoDeleted: false,
		Internal:    false,
		NoWait:      false,
		Arguments:   nil,
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("ERROR: Could not declare channel %q", err)
		return nil, err
	} else {

		log.Printf("Got Channel")
		defer ch.Close()

		err := declareExchange(c)
		if err != nil {
			log.Printf("ERROR: Could not declare Exchange [%s] %q", exchangeName, err)
			return nil, err
		} else {

			q, err := declareQueue(queue)
			if err != nil {
				log.Printf("ERROR: Could not declare queue [%s] %q", queue.Name, err)
				return nil, err
			} else {

				err = ch.QueueBind(
					q.Name,       // queue name
					"image",      // routing key @TODO
					exchangeName, // exchange
					false,
					nil,
				)
				if err != nil {
					log.Printf("ERROR: Could not bind [%s] queue to [%s] exhange %q", queue, exchangeName, err)
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
*/
