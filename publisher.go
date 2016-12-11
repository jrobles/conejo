package conejo

import (
	"github.com/streadway/amqp"
	"log"
)

func Publish(conn *amqp.Connection, queue Queue, exchange Exchange, body string) {

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("[CONEJO] Could not declare channel %q", err)
	}
	ch.Confirm(false)
	ack, nack := ch.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	err = declareExchange(exchange, ch)
	if err != nil {
		log.Printf("ERROR: Could not declare exchange %q", err)
	} else {
		log.Printf("[CONEJO] Declared exchange")
		_, err := declareQueue(queue, ch)
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

				if err = ch.Publish(
					exchange.Name, // publish to an exchange
					queue.Name,    // routing to 0 or more queues
					false,         // mandatory
					false,         // immediate
					amqp.Publishing{
						Headers:         amqp.Table{},
						ContentType:     "text/plain",
						ContentEncoding: "",
						Body:            []byte(body),
						DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
						Priority:        0,              // 0-9
					},
				); err != nil {
					log.Printf("ERROR: Could not publish message %s", err)
				} else {
					select {
					case tag := <-ack:
						log.Println("Acked ", tag)
					case tag := <-nack:
						log.Println("Nack alert! ", tag)
					}
					defer ch.Close()
				}

			}
		}
	}
}
