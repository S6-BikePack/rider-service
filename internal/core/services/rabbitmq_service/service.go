package rabbitmq_service

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"rider-service/internal/core/domain"
	"rider-service/pkg/rabbitmq"
)

type rabbitmqPublisher rabbitmq.RabbitMQ

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ) *rabbitmqPublisher {
	return &rabbitmqPublisher{Connection: rabbitmq.Connection, Channel: rabbitmq.Channel}
}

func (rmq *rabbitmqPublisher) CreateRider(rider domain.Rider) error {
	return rmq.publishJson("rider.create", rider)
}

func (rmq *rabbitmqPublisher) UpdateRider(rider domain.Rider) error {
	return rmq.publishJson("rider.update", rider)
}

func (rmq *rabbitmqPublisher) UpdateRiderLocation(id string, newLocation domain.Location) error {
	message := struct {
		id       string
		location domain.Location
	}{id: id, location: newLocation}

	return rmq.publishJson("rider.update.location", message)
}

func (rmq *rabbitmqPublisher) publishJson(topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	err = rmq.Channel.Publish(
		"topics",
		topic,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         js,
		},
	)

	return err
}
