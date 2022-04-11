package handlers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
	"rider-service/internal/core/domain"
	"rider-service/internal/core/ports"
	"rider-service/pkg/rabbitmq"
)

type rabbitmqHandler struct {
	rabbitmq *rabbitmq.RabbitMQ
	service  ports.RiderService
	handlers map[string]func(topic string, body []byte, handler *rabbitmqHandler) error
}

func NewRabbitMQ(rabbitmq *rabbitmq.RabbitMQ, service ports.RiderService) *rabbitmqHandler {
	return &rabbitmqHandler{
		rabbitmq: rabbitmq,
		service:  service,
		handlers: map[string]func(topic string, body []byte, handler *rabbitmqHandler) error{
			"user.create": UserCreateOrUpdate,
			"user.update": UserCreateOrUpdate,
		},
	}
}

func UserCreateOrUpdate(topic string, body []byte, handler *rabbitmqHandler) error {
	var user domain.User

	if err := json.Unmarshal(body, &user); err != nil {
		return err
	}

	if err := handler.service.SaveOrUpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (handler *rabbitmqHandler) Listen(queue string) {

	q, err := handler.rabbitmq.Channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	for _, s := range maps.Keys(handler.handlers) {
		err = handler.rabbitmq.Channel.QueueBind(
			q.Name,
			s,
			"topics",
			false,
			nil)
		if err != nil {
			return
		}
	}

	msgs, err := handler.rabbitmq.Channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			fun, exist := handler.handlers[msg.RoutingKey]

			if exist {
				err = fun(msg.RoutingKey, msg.Body, handler)
				if err == nil {
					msg.Ack(false)
					continue
				}
			}

			fmt.Println(err)
			msg.Nack(false, true)
		}
	}()

	<-forever
}

type MessageHandler struct {
	topic   string
	handler func(topic string, body []byte, handler *rabbitmqHandler) error
}
