package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"rider-service/config"
	"rider-service/internal/core/domain"
	"rider-service/internal/mock"
	"rider-service/pkg/rabbitmq"
	"testing"
)

type RabbitMQHandlerTestSuite struct {
	suite.Suite
	MockRiderService       *mock.RiderService
	MockServiceAreaService *mock.ServiceAreaService
	TestRabbitMQ           *rabbitmq.RabbitMQ
	TestHandler            *rabbitmqHandler
	Cfg                    *config.Config
	TestData               struct {
		User        domain.User
		ServiceArea domain.ServiceArea
	}
}

func (suite *RabbitMQHandlerTestSuite) SetupSuite() {
	cfgPath := "../../test/rider.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	mockRiderService := new(mock.RiderService)
	mockServiceAreaService := new(mock.ServiceAreaService)

	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		panic(errors.WithStack(err))
	}

	handler := NewRabbitMQ(rabbitMQ, mockRiderService, mockServiceAreaService, cfg)

	suite.Cfg = cfg
	suite.MockRiderService = mockRiderService
	suite.MockServiceAreaService = mockServiceAreaService
	suite.TestHandler = handler
	suite.TestRabbitMQ = rabbitMQ
	suite.TestData = struct {
		User        domain.User
		ServiceArea domain.ServiceArea
	}{
		User: domain.User{
			ID:       "test-id",
			Name:     "test-name",
			LastName: "test-lastname",
		},
		ServiceArea: domain.ServiceArea{
			ID:         1,
			Identifier: "test-area",
		},
	}

	go suite.TestHandler.Listen()
}

func (suite *RabbitMQHandlerTestSuite) SetupTest() {
	suite.MockRiderService.ExpectedCalls = nil
}

func (suite *RabbitMQHandlerTestSuite) TearDownSuite() {
	suite.TestRabbitMQ.Close()
}

func (suite *RabbitMQHandlerTestSuite) TestHandler_ServiceAreaCreateOrUpdate() {
	suite.MockServiceAreaService.On("SaveOrUpdateServiceArea", mock2.Anything).Return(nil)
	suite.MockRiderService.On("SaveOrUpdateUser", mock2.Anything).Return(nil)

	err := publishJson(suite.TestRabbitMQ, suite.Cfg.RabbitMQ.Exchange, "service_area.create", suite.TestData.ServiceArea)

	suite.NoError(err)

	for len(suite.MockServiceAreaService.Calls) < 1 {
	}

	suite.MockServiceAreaService.AssertCalled(suite.T(), "SaveOrUpdateServiceArea", suite.TestData.ServiceArea)
}

func (suite *RabbitMQHandlerTestSuite) TestHandler_UserCreateOrUpdate_Create() {
	suite.MockServiceAreaService.On("SaveOrUpdateServiceArea", mock2.Anything).Return(nil)
	suite.MockRiderService.On("SaveOrUpdateUser", mock2.Anything).Return(nil)

	err := publishJson(suite.TestRabbitMQ, suite.Cfg.RabbitMQ.Exchange, "user.create", suite.TestData.User)

	suite.NoError(err)

	for len(suite.MockRiderService.Calls) < 1 {
	}

	suite.MockRiderService.AssertCalled(suite.T(), "SaveOrUpdateUser", suite.TestData.User)
}

func (suite *RabbitMQHandlerTestSuite) TestHandler_UserCreateOrUpdate_Update() {
	suite.MockServiceAreaService.On("SaveOrUpdateServiceArea", mock2.Anything).Return(nil)
	suite.MockRiderService.On("SaveOrUpdateUser", mock2.Anything).Return(nil)

	err := publishJson(suite.TestRabbitMQ, suite.Cfg.RabbitMQ.Exchange, "user.update", suite.TestData.User)

	suite.NoError(err)

	for len(suite.MockRiderService.Calls) < 1 {
	}

	suite.MockRiderService.AssertCalled(suite.T(), "SaveOrUpdateUser", suite.TestData.User)
}

func publishJson(rabbitmq *rabbitmq.RabbitMQ, exchange, topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	err = rabbitmq.Channel.Publish(
		exchange,
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

func TestIntegration_RabbitMQHandlerTestSuite(t *testing.T) {
	repoSuite := new(RabbitMQHandlerTestSuite)
	suite.Run(t, repoSuite)
}
