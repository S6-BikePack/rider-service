package repositories

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rider-service/config"
	"rider-service/internal/core/domain"
	"testing"
)

type RiderRepositoryTestSuite struct {
	suite.Suite
	TestDb   *gorm.DB
	TestRepo *riderRepository
	Cfg      *config.Config
	TestData struct {
		Rider    domain.Rider
		Location domain.Location
	}
}

func (suite *RiderRepositoryTestSuite) SetupSuite() {
	cfgPath := "../../test/rider.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database)
	db, err := gorm.Open(postgres.Open(dsn))
	db.Debug()

	if err != nil {
		panic(errors.WithStack(err))
	}

	repository, err := NewRiderRepository(db)

	if err != nil {
		panic(errors.WithStack(err))
	}

	db.Exec("DELETE FROM public.riders")
	db.Exec("DELETE FROM public.users")
	db.Exec("DELETE FROM public.service_areas")

	db.Exec("INSERT INTO public.users (id, name, last_name) VALUES ('test-id', 'test-name', 'test-lastname')")
	db.Exec("INSERT INTO public.service_areas (id, identifier) VALUES (1, 'test-area')")

	suite.Cfg = cfg
	suite.TestDb = db
	suite.TestRepo = repository
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

func (suite *RiderRepositoryTestSuite) TestRepository_Get() {
	suite.TestDb.Exec("INSERT INTO public.riders (user_id, status, service_area_id, width, height, depth, location) VALUES ('test-id', 1, 1, 100, 100, 100,'0101000020E61000000000000000000040000000000000F03F'::geometry(Point,4326))")

	result, err := suite.TestRepo.Get(context.Background(), suite.TestData.Rider.UserID)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.Rider, result)
}

func (suite *RiderRepositoryTestSuite) TestRepository_Get_NotFound() {
	_, err := suite.TestRepo.Get(context.Background(), "test")

	suite.Error(err)
}

func (suite *RiderRepositoryTestSuite) TestRepository_Save() {
	suite.TestDb.Exec("INSERT INTO public.users (id, name, last_name) VALUES ('test-id-2', 'test-name', 'test-lastname')")

	newRider := suite.TestData.Rider
	newRider.User = domain.User{}
	newRider.UserID = "test-id-2"

	_, err := suite.TestRepo.Save(context.Background(), newRider)

	suite.NoError(err)

	queryResult := domain.Rider{}
	suite.TestDb.Raw("SELECT * FROM public.riders WHERE user_id=?",
		newRider.UserID).Scan(&queryResult)

	suite.EqualValues(queryResult.Capacity, newRider.Capacity)
}

func (suite *RiderRepositoryTestSuite) TestRepository_Update() {
	updated := suite.TestData.Rider
	updated.UserID = "test-id-2"
	updated.Capacity = domain.Dimensions{
		Width:  1,
		Height: 1,
		Depth:  1,
	}
	updated.User = domain.User{}

	_, err := suite.TestRepo.Update(context.Background(), updated)

	suite.NoError(err)

	queryResult := domain.Rider{}
	suite.TestDb.Raw("SELECT * FROM public.riders WHERE user_id=?",
		updated.UserID).Scan(&queryResult)

	suite.EqualValues(queryResult.Capacity, updated.Capacity)
}

func (suite *RiderRepositoryTestSuite) TestRepository_SaveUser() {
	user := domain.User{
		ID:       "test-id-3",
		Name:     "test-name-3",
		LastName: "test-lastname-3",
	}

	err := suite.TestRepo.SaveOrUpdateUser(context.Background(), user)

	suite.NoError(err)

	queryResult := domain.User{}
	suite.TestDb.Raw("SELECT * FROM public.users WHERE id=?",
		user.ID).Scan(&queryResult)

	suite.EqualValues(queryResult, user)
}

func (suite *RiderRepositoryTestSuite) TestRepository_UpdateUser() {
	user := domain.User{
		ID:       "test-id-3",
		Name:     "new-name-3",
		LastName: "new-lastname-3",
	}

	err := suite.TestRepo.SaveOrUpdateUser(context.Background(), user)

	suite.NoError(err)

	queryResult := domain.User{}
	suite.TestDb.Raw("SELECT * FROM public.users WHERE id=?",
		user.ID).Scan(&queryResult)

	suite.EqualValues(queryResult, user)
}

func (suite *RiderRepositoryTestSuite) TestRepository_GetUser() {
	user := domain.User{
		ID:       "test-id",
		Name:     "test-name",
		LastName: "test-lastname",
	}

	result, err := suite.TestRepo.GetUser(context.Background(), user.ID)

	suite.NoError(err)

	suite.EqualValues(result, user)
}

func TestIntegration_RiderRepositoryTestSuite(t *testing.T) {
	repoSuite := new(RiderRepositoryTestSuite)
	suite.Run(t, repoSuite)
}
