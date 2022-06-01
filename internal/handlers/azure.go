package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"rider-service/config"
	"rider-service/internal/core/domain"
	"rider-service/internal/core/interfaces"
	"rider-service/pkg/azure"
)

type azureHandler struct {
	serviceBus         *azure.ServiceBus
	service            interfaces.RiderService
	serviceAreaService interfaces.ServiceAreaService
	handlers           map[string]func(topic string, body []byte, handler *azureHandler) error
	config             *config.Config
	channel            chan bool
}

func NewAzure(serviceBus *azure.ServiceBus, service interfaces.RiderService, serviceAreaService interfaces.ServiceAreaService, config *config.Config) *azureHandler {
	return &azureHandler{
		serviceBus:         serviceBus,
		service:            service,
		serviceAreaService: serviceAreaService,
		handlers: map[string]func(topic string, body []byte, handler *azureHandler) error{
			"user.create":         userCreateOrUpdate,
			"user.update":         userCreateOrUpdate,
			"service_area.create": serviceAreaCreateOrUpdate,
			"service_area.update": serviceAreaCreateOrUpdate,
		},
		config: config,
	}
}

func serviceAreaCreateOrUpdate(topic string, body []byte, handler *azureHandler) error {
	var serviceArea domain.ServiceArea
	if err := json.Unmarshal(body, &serviceArea); err != nil {
		return err
	}

	if err := handler.serviceAreaService.SaveOrUpdateServiceArea(serviceArea); err != nil {
		return err
	}

	return nil
}

func userCreateOrUpdate(topic string, body []byte, handler *azureHandler) error {
	var user domain.User

	if err := json.Unmarshal(body, &user); err != nil {
		return err
	}

	if err := handler.service.SaveOrUpdateUser(context.Background(), user); err != nil {
		return err
	}

	return nil
}

func (handler *azureHandler) Listen() {

	receiver, err := handler.serviceBus.Client.NewReceiverForQueue(
		handler.config.AzureServiceBus.QueueName,
		nil,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	handler.channel = make(chan bool)

	go func() {
		for {
			select {
			case <-handler.channel:
				return
			default:
				msgs, err := receiver.ReceiveMessages(
					context.Background(),
					1,
					nil,
				)

				if err != nil {
					fmt.Println(err)
					return
				}

				for _, msg := range msgs {

					if msg.Subject != nil && *msg.Subject != "" {
						fun, exist := handler.handlers[*msg.Subject]

						if exist {
							err = fun(*msg.Subject, msg.Body, handler)
							if err == nil {
								_ = receiver.CompleteMessage(context.Background(), msg, nil)
								continue
							}
						}
					} else {
						fmt.Println("Message contains no subject: ", msg.MessageID)
						_ = receiver.CompleteMessage(context.Background(), msg, nil)
						continue
					}

					fmt.Println(err)
					_ = receiver.AbandonMessage(context.Background(), msg, nil)
				}
			}
		}
	}()
}

func (handler *azureHandler) Quit() {
	handler.channel <- true
}
