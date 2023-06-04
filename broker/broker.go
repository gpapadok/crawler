package broker

import (
	"context"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp.Connection, error) {
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s",
		os.Getenv("BROKER_USER"),
		os.Getenv("BROKER_PASSWORD"),
		os.Getenv("BROKER_HOST"),
		os.Getenv("BROKER_PORT"),
	)

	return amqp.Dial(amqpURI)
}

func CreateChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	return conn.Channel()
}

func CreateQueue(ch *amqp.Channel, name string) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		name,
		false, // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &q, err
}

func Publish(ch *amqp.Channel, q *amqp.Queue, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
}

func Consume(ch *amqp.Channel, q *amqp.Queue) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		q.Name,
		"",
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
}
