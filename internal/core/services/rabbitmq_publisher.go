package services

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"rider-service/config"
	"rider-service/internal/core/domain"
	"rider-service/pkg/rabbitmq"
)

type rabbitmqPublisher struct {
	rabbitmq *rabbitmq.RabbitMQ
	tracer   trace.Tracer
	config   *config.Config
}

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ, tracerProvider trace.TracerProvider, cfg *config.Config) *rabbitmqPublisher {
	return &rabbitmqPublisher{rabbitmq: rabbitmq, tracer: tracerProvider.Tracer("RabbitMQ.Publisher"), config: cfg}
}

func (rmq *rabbitmqPublisher) CreateRider(ctx context.Context, rider domain.Rider) error {
	return rmq.publishJson(ctx, "create", rider)
}

func (rmq *rabbitmqPublisher) UpdateRider(ctx context.Context, rider domain.Rider) error {
	return rmq.publishJson(ctx, "update", rider)
}

func (rmq *rabbitmqPublisher) UpdateRiderLocation(ctx context.Context, serviceArea domain.ServiceArea, id string, newLocation domain.Location) error {
	message := struct {
		Id       string
		Location domain.Location
	}{Id: id, Location: newLocation}

	return rmq.publishJson(ctx, serviceArea.Identifier+".update.location", message)
}

func (rmq *rabbitmqPublisher) publishJson(ctx context.Context, topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	_, span := rmq.tracer.Start(ctx, "publish")

	span.AddEvent(
		"Published message to rabbitmq",
		trace.WithAttributes(
			attribute.String("topic", topic),
			attribute.String("body", string(js))))
	span.End()

	err = rmq.rabbitmq.Channel.Publish(
		rmq.config.RabbitMQ.Exchange,
		fmt.Sprintf("rider.%s", topic),
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
