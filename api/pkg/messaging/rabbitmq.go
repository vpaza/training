package messaging

import (
	"context"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection
var channel *amqp091.Channel
var queue map[string]*amqp091.Queue

func init() {
	queue = make(map[string]*amqp091.Queue)
}

func Connect(dsn string) error {
	c, err := amqp091.Dial(dsn)
	if err != nil {
		return err
	}

	conn = c
	channel, err = conn.Channel()
	if err != nil {
		return err
	}

	return nil
}

func DeclareQueue(name string) error {
	q, err := channel.QueueDeclare(name, false, false, false, false, nil)
	if err != nil {
		return err
	}

	queue[name] = &q
	return nil
}

func Publish(queue string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := channel.PublishWithContext(ctx,
		"",
		queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)

	return err
}
