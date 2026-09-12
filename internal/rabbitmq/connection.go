package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RabbitMQConnect(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	go func() {
		closeErr := <-conn.NotifyClose(make(chan *amqp.Error))
		if closeErr != nil {
			log.Printf("RabbitMQ connection closed: %v", closeErr)
		}
	}()

	return conn, nil
}
