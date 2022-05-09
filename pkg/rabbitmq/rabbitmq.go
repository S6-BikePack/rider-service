package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"rider-service/config"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(cfg *config.Config) (*RabbitMQ, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)

	conn, err := amqp.Dial(connStr)

	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(
		cfg.RabbitMQ.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	err = channel.Qos(
		1,
		0,
		false,
	)

	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}, nil
}

func (r *RabbitMQ) Close() {
	_ = r.Connection.Close()
}
