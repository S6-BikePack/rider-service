package services

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/sdk/trace"
	"rider-service/config"
	"rider-service/internal/core/domain"
	"rider-service/internal/core/interfaces"
	"rider-service/internal/mock"
	"rider-service/pkg/rabbitmq"
	"testing"
)

type RabbitMQPublisherTestSuite struct {
	suite.Suite
	MockService   *mock.RiderService
	TestRabbitMQ  *rabbitmq.RabbitMQ
	TestPublisher interfaces.MessageBusPublisher
	Cfg           *config.Config
	TestData      struct {
		Rider    domain.Rider
		Location domain.Location
	}
}

func (suite *RabbitMQPublisherTestSuite) SetupSuite() {
	cfgPath := "../../../test/rider.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	mockService := new(mock.RiderService)

	rmqServer, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		panic(errors.WithStack(err))
	}

	tracer := trace.NewTracerProvider()

	rmqPublisher := NewRabbitMQPublisher(rmqServer, tracer, cfg)

	suite.Cfg = cfg
	suite.MockService = mockService
	suite.TestRabbitMQ = rmqServer
	suite.TestPublisher = rmqPublisher
	suite.TestData = struct {
		Rider    domain.Rider
		Location domain.Location
	}{
		Rider: domain.Rider{
			UserID: "test-id",
			User: domain.User{
				ID:       "test-id",
				Name:     "test-name",
				LastName: "test-lastname",
			},
			Status:        1,
			ServiceAreaID: 1,
			ServiceArea: domain.ServiceArea{
				ID:         1,
				Identifier: "test-area",
			},
			Capacity: domain.Dimensions{
				Width:  100,
				Height: 100,
				Depth:  100,
			},
			Location: domain.Location{
				Latitude:  1,
				Longitude: 2,
			},
		},
		Location: domain.Location{
			Latitude:  2,
			Longitude: 3,
		},
	}
}

func (suite *RabbitMQPublisherTestSuite) TestRabbitMQPublisher_CreateRider() {
	ch, err := suite.TestRabbitMQ.Connection.Channel()

	suite.NoError(err)

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	suite.NoError(err)

	err = ch.QueueBind(
		queue.Name,
		"rider.create",
		suite.Cfg.RabbitMQ.Exchange,
		false,
		nil)
	if err != nil {
		return
	}

	suite.NoError(err)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	suite.NoError(err)

	err = suite.TestPublisher.CreateRider(context.Background(), suite.TestData.Rider)

	suite.NoError(err)

	for msg := range msgs {
		suite.Equal("rider.create", msg.RoutingKey)

		var rider domain.Rider

		err = json.Unmarshal(msg.Body, &rider)
		suite.NoError(err)

		suite.Equal(suite.TestData.Rider, rider)

		err = msg.Ack(true)

		suite.NoError(err)

		err = ch.Close()

		suite.NoError(err)

		return
	}

}

func (suite *RabbitMQPublisherTestSuite) TestRabbitMQPublisher_UpdateRider() {
	ch, err := suite.TestRabbitMQ.Connection.Channel()

	suite.NoError(err)

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	suite.NoError(err)

	err = ch.QueueBind(
		queue.Name,
		"rider.update",
		suite.Cfg.RabbitMQ.Exchange,
		false,
		nil)
	if err != nil {
		return
	}

	suite.NoError(err)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	suite.NoError(err)

	err = suite.TestPublisher.UpdateRider(context.Background(), suite.TestData.Rider)

	suite.NoError(err)

	for msg := range msgs {
		suite.Equal("rider.update", msg.RoutingKey)

		var rider domain.Rider

		err = json.Unmarshal(msg.Body, &rider)
		suite.NoError(err)

		suite.Equal(suite.TestData.Rider, rider)

		err = msg.Ack(true)

		suite.NoError(err)

		err = ch.Close()

		suite.NoError(err)

		return
	}
}

func (suite *RabbitMQPublisherTestSuite) TestRabbitMQPublisher_UpdateRiderLocation() {
	ch, err := suite.TestRabbitMQ.Connection.Channel()

	suite.NoError(err)

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	suite.NoError(err)

	err = ch.QueueBind(
		queue.Name,
		"rider."+suite.TestData.Rider.ServiceArea.Identifier+".update.location",
		suite.Cfg.RabbitMQ.Exchange,
		false,
		nil)
	if err != nil {
		return
	}

	suite.NoError(err)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	suite.NoError(err)

	err = suite.TestPublisher.UpdateRiderLocation(context.Background(), suite.TestData.Rider.ServiceArea, suite.TestData.Rider.UserID, suite.TestData.Location)

	suite.NoError(err)

	for msg := range msgs {
		suite.Equal("rider."+suite.TestData.Rider.ServiceArea.Identifier+".update.location", msg.RoutingKey)

		var message struct {
			Id       string
			Location domain.Location
		}

		err = json.Unmarshal(msg.Body, &message)
		suite.NoError(err)

		suite.Equal(suite.TestData.Rider.UserID, message.Id)
		suite.Equal(suite.TestData.Location, message.Location)

		err = msg.Ack(true)

		suite.NoError(err)

		err = ch.Close()

		suite.NoError(err)

		return
	}
}

func TestIntegration_RabbitMQPublisherTestSuite(t *testing.T) {
	repoSuite := new(RabbitMQPublisherTestSuite)
	suite.Run(t, repoSuite)
}
