package broker

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp.Connection, error) {
	amqpURL := "amqp://guest:guest@localhost:5672" // TODO: Move to env variable

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CreateChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
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

	return &q, nil
}

func Publish(ch *amqp.Channel, q *amqp.Queue, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	return err
}

func Consume(ch *amqp.Channel, q *amqp.Queue) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		q.Name,
		"",
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
}
