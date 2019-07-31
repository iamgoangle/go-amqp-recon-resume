package main

import (
	"log"

	"github.com/iamgoangle/go-amqp-recon-resume/rabbitmq"
)

func main() {
	conn, err := rabbitmq.NewDial("amqp://admin:1234@localhost:5672/")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	ch, err := conn.NewChannel()
	if err != nil {
		log.Panic(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"direct.exchange", // name
		"direct",          // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Panic(err)
	}

	q, err := ch.QueueDeclare(
		"test.golf", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Panic(err)
	}

	err = ch.QueueBind(
		q.Name,            // queue name
		"direct.routing",  // routing key
		"direct.exchange", // exchange
		false,
		nil)
	if err != nil {
		log.Panic(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		log.Panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
