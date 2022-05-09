package services

import (
	"context"
	"github.com/pkg/errors"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"rider-service/internal/core/domain"
	"rider-service/internal/core/interfaces"
	"rider-service/internal/mock"
	"testing"
)

type RiderServiceTestSuite struct {
	suite.Suite
	MockRepository *mock.RiderRepository
	MockPublisher  *mock.MessageBusPublisher
	TestService    interfaces.RiderService
	TestData       struct {
		Rider    domain.Rider
		Location domain.Location
	}
}

func (suite *RiderServiceTestSuite) SetupSuite() {
	repository := new(mock.RiderRepository)
	publisher := new(mock.MessageBusPublisher)

	srv := NewRiderService(repository, publisher)

	suite.MockRepository = repository
	suite.MockPublisher = publisher
	suite.TestService = srv
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

func (suite *RiderServiceTestSuite) SetupTest() {
	suite.MockPublisher.ExpectedCalls = nil
	suite.MockRepository.ExpectedCalls = nil
}

func (suite *RiderServiceTestSuite) TestRiderService_GetAll() {
	suite.MockRepository.On("GetAll").Return([]domain.Rider{suite.TestData.Rider}, nil)

	result, err := suite.TestService.GetAll(context.Background())

	suite.NoError(err)

	suite.MockRepository.AssertCalled(suite.T(), "GetAll")
	suite.Equal(1, len(result))
	suite.EqualValues(suite.TestData.Rider, result[0])
}

func (suite *RiderServiceTestSuite) TestRiderService_Get() {
	suite.MockRepository.On("Get", suite.TestData.Rider.UserID).Return(suite.TestData.Rider, nil)

	result, err := suite.TestService.Get(context.Background(), suite.TestData.Rider.UserID)

	suite.NoError(err)

	suite.MockRepository.AssertCalled(suite.T(), "Get", suite.TestData.Rider.UserID)
	suite.EqualValues(suite.TestData.Rider, result)
}

func (suite *RiderServiceTestSuite) TestRiderService_Get_NotFound() {
	suite.MockRepository.On("Get", suite.TestData.Rider.UserID).Return(domain.Rider{}, errors.New("could not find rider"))

	result, err := suite.TestService.Get(context.Background(), suite.TestData.Rider.UserID)

	suite.Error(err)
	suite.EqualValues(domain.Rider{}, result)

	suite.MockRepository.AssertCalled(suite.T(), "Get", suite.TestData.Rider.UserID)
}

func (suite *RiderServiceTestSuite) TestRiderService_Create() {
	suite.MockRepository.On("GetUser", suite.TestData.Rider.UserID).Return(suite.TestData.Rider.User, nil)
	suite.MockRepository.On("Save", mock2.Anything).Return(suite.TestData.Rider, nil)
	suite.MockPublisher.On("CreateRider", suite.TestData.Rider).Return(nil)

	result, err := suite.TestService.Create(context.Background(), suite.TestData.Rider.UserID, suite.TestData.Rider.ServiceAreaID, suite.TestData.Rider.Capacity)

	suite.NoError(err)

	suite.MockPublisher.AssertCalled(suite.T(), "CreateRider", suite.TestData.Rider)
	suite.EqualValues(suite.TestData.Rider, result)
}

func (suite *RiderServiceTestSuite) TestRiderService_Create_UserNotFound() {
	suite.MockRepository.On("GetUser", suite.TestData.Rider.UserID).Return(domain.User{}, errors.New("user not found"))

	_, err := suite.TestService.Create(context.Background(), suite.TestData.Rider.UserID, suite.TestData.Rider.ServiceAreaID, suite.TestData.Rider.Capacity)

	suite.MockRepository.AssertNotCalled(suite.T(), "Save")
	suite.Error(err)
}

func (suite *RiderServiceTestSuite) TestRiderService_Create_CouldNotSave() {
	suite.MockRepository.On("GetUser", suite.TestData.Rider.UserID).Return(suite.TestData.Rider.User, nil)
	suite.MockRepository.On("Save", mock2.Anything).Return(domain.Rider{}, errors.New("could not save rider"))
	suite.MockPublisher.On("CreateRider", suite.TestData.Rider).Return(nil)

	_, err := suite.TestService.Create(context.Background(), suite.TestData.Rider.UserID, suite.TestData.Rider.ServiceAreaID, suite.TestData.Rider.Capacity)

	suite.Error(err)

	suite.MockPublisher.AssertNotCalled(suite.T(), "CreateRider")
}

func (suite *RiderServiceTestSuite) TestRiderService_Update() {
	updated := suite.TestData.Rider
	updated.Capacity = domain.Dimensions{Width: 5, Height: 5, Depth: 5}

	suite.MockRepository.On("Get", suite.TestData.Rider.UserID).Return(suite.TestData.Rider, nil)
	suite.MockRepository.On("Update", updated).Return(updated, nil)
	suite.MockPublisher.On("UpdateRider", updated).Return(nil)

	result, err := suite.TestService.Update(context.Background(), suite.TestData.Rider.UserID, suite.TestData.Rider.Status, suite.TestData.Rider.ServiceAreaID, updated.Capacity)

	suite.NoError(err)

	suite.MockPublisher.AssertCalled(suite.T(), "UpdateRider", updated)
	suite.EqualValues(updated, result)
}

func (suite *RiderServiceTestSuite) TestRiderService_UpdateLocation() {
	updated := suite.TestData.Rider
	updated.Location = suite.TestData.Location

	suite.MockRepository.On("Get", suite.TestData.Rider.UserID).Return(suite.TestData.Rider, nil)
	suite.MockRepository.On("Update", updated).Return(updated, nil)
	suite.MockPublisher.On("UpdateRiderLocation", updated.ServiceArea, updated.UserID, updated.Location).Return(nil)

	result, err := suite.TestService.UpdateLocation(context.Background(), suite.TestData.Rider.UserID, updated.Location)

	suite.NoError(err)

	suite.MockPublisher.AssertCalled(suite.T(), "UpdateRiderLocation", updated.ServiceArea, updated.UserID, updated.Location)
	suite.EqualValues(updated, result)
}

func TestUnit_RiderServiceTestSuite(t *testing.T) {
	repoSuite := new(RiderServiceTestSuite)
	suite.Run(t, repoSuite)
}
