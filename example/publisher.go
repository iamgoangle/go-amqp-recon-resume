package main

import (
	"log"

	"github.com/iamgoangle/go-amqp-recon-resume/rabbitmq"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := rabbitmq.NewDial("amqp://admin:1234@localhost:5672/")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
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

	body := "Hello, World"
	err = ch.Publish(
		"direct.exchange", // exchange
		"direct.routing",  // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Panic(err)
	}

	log.Printf(" [x] Sent %s", body)
}
