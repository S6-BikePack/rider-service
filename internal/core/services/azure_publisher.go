package services

import (
	"context"
	"encoding/json"
	"fmt"
	"rider-service/config"
	"rider-service/internal/core/domain"
	"rider-service/pkg/azure"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type azurePublisher struct {
	serviceBus *azure.ServiceBus
	sender     *azservicebus.Sender
	config     *config.Config
}

func NewAzurePublisher(serviceBus *azure.ServiceBus, cfg *config.Config) *azurePublisher {
	return &azurePublisher{serviceBus: serviceBus, config: cfg}
}

func (az *azurePublisher) CreateRider(ctx context.Context, rider domain.Rider) error {
	return az.publishJson(ctx, "create", rider)
}

func (az *azurePublisher) UpdateRider(ctx context.Context, rider domain.Rider) error {
	return az.publishJson(ctx, "update", rider)
}

func (az *azurePublisher) UpdateRiderLocation(ctx context.Context, serviceArea domain.ServiceArea, id string, newLocation domain.Location) error {
	message := struct {
		Id       string
		Location domain.Location
	}{Id: id, Location: newLocation}

	return az.publishJson(ctx, serviceArea.Identifier+".update.location", message)
}

func (az *azurePublisher) publishJson(ctx context.Context, topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	topic = fmt.Sprintf("customer.%s", topic)

	sender, err := az.serviceBus.Client.NewSender(topic, nil)

	defer func(sender *azservicebus.Sender, ctx context.Context) {
		_ = sender.Close(ctx)
	}(sender, ctx)

	if err != nil {
		return err
	}

	err = sender.SendMessage(ctx, &azservicebus.Message{
		Body:    js,
		Subject: &topic,
	}, nil)

	if err != nil {
		return err
	}

	return err
}
