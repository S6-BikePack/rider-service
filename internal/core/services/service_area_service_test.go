package services

import (
	"github.com/stretchr/testify/suite"
	"rider-service/internal/core/domain"
	"rider-service/internal/core/interfaces"
	"rider-service/internal/mock"
	"testing"
)

type ServiceAreaServiceTestSuite struct {
	suite.Suite
	MockRepository *mock.ServiceAreaRepository
	TestService    interfaces.ServiceAreaService
	TestData       struct {
		ServiceArea domain.ServiceArea
	}
}

func (suite *ServiceAreaServiceTestSuite) SetupSuite() {
	repository := new(mock.ServiceAreaRepository)

	srv := NewServiceAreaService(repository)

	suite.MockRepository = repository
	suite.TestService = srv
	suite.TestData = struct {
		ServiceArea domain.ServiceArea
	}{
		ServiceArea: domain.ServiceArea{
			ID:         1,
			Identifier: "test-area",
		},
	}
}

func (suite *ServiceAreaServiceTestSuite) TestServiceAreaService_SaveOrUpdateServiceArea() {
	suite.MockRepository.On("SaveOrUpdateServiceArea", suite.TestData.ServiceArea).Return(nil)

	err := suite.TestService.SaveOrUpdateServiceArea(suite.TestData.ServiceArea)

	if err != nil {
		return
	}

	suite.MockRepository.AssertCalled(suite.T(), "SaveOrUpdateServiceArea", suite.TestData.ServiceArea)
}

func TestUnit_ServiceAreaServiceTestSuite(t *testing.T) {
	repoSuite := new(ServiceAreaServiceTestSuite)
	suite.Run(t, repoSuite)
}
