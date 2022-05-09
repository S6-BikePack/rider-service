package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"rider-service/config"
	"rider-service/internal/core/domain"
	"rider-service/internal/mock"
	"rider-service/pkg/dto"
	"rider-service/pkg/logging"
	"strings"
	"testing"
)

type RestHandlerTestSuite struct {
	suite.Suite
	MockService *mock.RiderService
	TestHandler *HTTPHandler
	TestRouter  *gin.Engine
	Cfg         *config.Config
	TestData    struct {
		Rider    domain.Rider
		Location domain.Location
	}
}

func (suite *RestHandlerTestSuite) SetupSuite() {
	cfgPath := "../../test/rider.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	logger := logging.MockLogger{}

	mockService := new(mock.RiderService)

	router := gin.New()
	gin.SetMode(gin.TestMode)

	deliveryHandler := NewHTTPHandler(mockService, router, logger, cfg)
	deliveryHandler.SetupEndpoints()

	suite.Cfg = cfg
	suite.MockService = mockService
	suite.TestRouter = router
	suite.TestHandler = deliveryHandler
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

func (suite *RestHandlerTestSuite) SetupTest() {
	suite.MockService.ExpectedCalls = nil
}

func (suite *RestHandlerTestSuite) TestHandler_GetAll() {
	suite.MockService.On("GetAll").Return([]domain.Rider{suite.TestData.Rider}, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/riders", nil)
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.RiderListResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.Len(responseObject, 1)

	suite.EqualValues(suite.TestData.Rider.User.Name, responseObject[0].Name)
	suite.EqualValues(suite.TestData.Rider.UserID, responseObject[0].ID)
	suite.EqualValues(suite.TestData.Rider.ServiceAreaID, responseObject[0].ServiceAreaID)
	suite.EqualValues(suite.TestData.Rider.Status, responseObject[0].Status)
}

func (suite *RestHandlerTestSuite) TestHandler_GetAll_NoneFound() {
	suite.MockService.On("GetAll").Return([]domain.Rider{}, errors.New("Not found"))

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/riders", nil)
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusNotFound, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Get() {
	suite.MockService.On("Get", suite.TestData.Rider.UserID).Return(suite.TestData.Rider, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/riders/%s", suite.TestData.Rider.UserID), nil)
	request.Header.Set("X-User-Id", suite.TestData.Rider.UserID)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.RiderResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.Rider.UserID, responseObject.ID)
	suite.EqualValues(suite.TestData.Rider.ServiceArea, domain.ServiceArea(responseObject.ServiceArea))
	suite.EqualValues(suite.TestData.Rider.Capacity, domain.Dimensions(responseObject.Capacity))
	suite.EqualValues(suite.TestData.Rider.Location, domain.Location(responseObject.Location))
}

func (suite *RestHandlerTestSuite) TestHandler_Get_BadID() {
	suite.MockService.On("Get", "test").Return(domain.Rider{}, nil)

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/riders/%s", "test"), nil)
	request.Header.Set("X-User-Id", suite.TestData.Rider.UserID)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Get_NotFound() {
	suite.MockService.On("Get", suite.TestData.Rider.UserID).Return(domain.Rider{}, errors.New("Not found"))

	rr := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/riders/%s", suite.TestData.Rider.UserID), nil)
	request.Header.Set("X-User-Id", suite.TestData.Rider.UserID)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusNotFound, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Create() {
	suite.MockService.On("Create", suite.TestData.Rider.UserID, suite.TestData.Rider.ServiceAreaID, suite.TestData.Rider.Capacity).Return(suite.TestData.Rider, nil)

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateRider{
		ID:          suite.TestData.Rider.UserID,
		ServiceArea: suite.TestData.Rider.ServiceAreaID,
		Capacity:    dto.CreateDimensions(suite.TestData.Rider.Capacity),
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/riders", strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)
	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.RiderResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.Rider.UserID, responseObject.ID)
	suite.EqualValues(suite.TestData.Rider.ServiceAreaID, responseObject.ServiceArea.ID)
	suite.EqualValues(suite.TestData.Rider.Capacity, domain.Dimensions(responseObject.Capacity))
}

func (suite *RestHandlerTestSuite) TestHandler_Create_BadInput() {
	rr := httptest.NewRecorder()

	data, err := json.Marshal(struct {
		Test string
	}{
		Test: suite.TestData.Rider.UserID,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/riders", strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Create_CouldNotCreate() {
	suite.MockService.On("Create", suite.TestData.Rider.UserID, suite.TestData.Rider.ServiceAreaID, suite.TestData.Rider.Capacity).Return(domain.Rider{}, errors.New("could not create"))

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateRider{
		ID:          suite.TestData.Rider.UserID,
		ServiceArea: suite.TestData.Rider.ServiceAreaID,
		Capacity:    dto.CreateDimensions(suite.TestData.Rider.Capacity),
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPost, "/api/riders", strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_Update() {
	suite.MockService.On("Update", suite.TestData.Rider.UserID, 2, suite.TestData.Rider.ServiceAreaID, suite.TestData.Rider.Capacity).Return(suite.TestData.Rider, nil)

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateRider{
		ServiceArea: suite.TestData.Rider.ServiceAreaID,
		Capacity:    dto.CreateDimensions(suite.TestData.Rider.Capacity),
		Status:      2,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/riders/%s", suite.TestData.Rider.UserID), strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)
	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.RiderResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.Rider.UserID, responseObject.ID)
	suite.EqualValues(suite.TestData.Rider.ServiceAreaID, responseObject.ServiceArea.ID)
	suite.EqualValues(suite.TestData.Rider.Capacity, domain.Dimensions(responseObject.Capacity))
}

func (suite *RestHandlerTestSuite) TestHandler_Update_CouldNotCreate() {
	suite.MockService.On("Update", suite.TestData.Rider.UserID, 2, suite.TestData.Rider.ServiceAreaID, suite.TestData.Rider.Capacity).Return(domain.Rider{}, errors.New("could not update"))

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyCreateRider{
		ID:          suite.TestData.Rider.UserID,
		ServiceArea: suite.TestData.Rider.ServiceAreaID,
		Capacity:    dto.CreateDimensions(suite.TestData.Rider.Capacity),
		Status:      2,
	})

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/riders/%s", suite.TestData.Rider.UserID), strings.NewReader(string(data)))
	request.Header.Set("X-User-Claims", `{"admin": true}`)

	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func (suite *RestHandlerTestSuite) TestHandler_UpdateLocation() {
	suite.MockService.On("UpdateLocation", suite.TestData.Rider.UserID, suite.TestData.Location).Return(suite.TestData.Rider, nil)

	rr := httptest.NewRecorder()

	data, err := json.Marshal(dto.BodyLocation(suite.TestData.Location))

	suite.NoError(err)

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/riders/%s/location", suite.TestData.Rider.UserID), strings.NewReader(string(data)))
	request.Header.Set("X-User-Id", suite.TestData.Rider.UserID)
	suite.NoError(err)

	suite.TestRouter.ServeHTTP(rr, request)

	suite.Equal(http.StatusOK, rr.Code)

	var responseObject dto.RiderResponse
	err = json.NewDecoder(rr.Body).Decode(&responseObject)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.Rider.UserID, responseObject.ID)
	suite.EqualValues(suite.TestData.Rider.Location, responseObject.Location)
}

func TestIntegration_RestHandlerTestSuite(t *testing.T) {
	repoSuite := new(RestHandlerTestSuite)
	suite.Run(t, repoSuite)
}
