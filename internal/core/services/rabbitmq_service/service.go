package rabbitmq_service

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"rider-service/internal/core/domain"
	"rider-service/pkg/rabbitmq"
)

type rabbitmqPublisher struct {
	rabbitmq *rabbitmq.RabbitMQ
	tracer   trace.Tracer
}

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ, tracerProvider trace.TracerProvider) *rabbitmqPublisher {
	return &rabbitmqPublisher{rabbitmq: rabbitmq, tracer: tracerProvider.Tracer("RabbitMQ.Publisher")}
}

func (rmq *rabbitmqPublisher) CreateRider(ctx context.Context, rider domain.Rider) error {
	return rmq.publishJson(ctx, "rider.create", rider)
}

func (rmq *rabbitmqPublisher) UpdateRider(ctx context.Context, rider domain.Rider) error {
	return rmq.publishJson(ctx, "rider.update", rider)
}

func (rmq *rabbitmqPublisher) UpdateRiderLocation(ctx context.Context, serviceArea domain.ServiceArea, id string, newLocation domain.Location) error {
	message := struct {
		Id       string
		Location domain.Location
	}{Id: id, Location: newLocation}

	return rmq.publishJson(ctx, "rider."+serviceArea.Identifier+".update.location", message)
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
