package repositories

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rider-service/config"
	"rider-service/internal/core/domain"
	"rider-service/pkg/logging"
	"testing"
)

type ServiceAreaRepositoryTestSuite struct {
	suite.Suite
	TestDb   *gorm.DB
	TestRepo *serviceAreaRepository
	Cfg      *config.Config
	TestData struct {
		ServiceArea domain.ServiceArea
	}
}

func (suite *ServiceAreaRepositoryTestSuite) SetupSuite() {
	cfgPath := "../../test/rider.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	logger, err := logging.NewSugaredOtelZap(cfg)
	defer func(logger *logging.OtelzapSugaredLogger) {
		_ = logger.Close()
	}(logger)

	if err != nil {
		panic(errors.WithStack(err))
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic(errors.WithStack(err))
	}

	repository, err := NewServiceAreaRepository(db)

	if err != nil {
		panic(errors.WithStack(err))
	}

	db.Exec("DELETE FROM public.service_areas")

	suite.Cfg = cfg
	suite.TestDb = db
	suite.TestRepo = repository
	suite.TestData = struct {
		ServiceArea domain.ServiceArea
	}{
		ServiceArea: domain.ServiceArea{
			ID:         1,
			Identifier: "test-area",
		},
	}
}

func (suite *ServiceAreaRepositoryTestSuite) TestRepository_SaveServiceArea() {

	err := suite.TestRepo.SaveOrUpdateServiceArea(suite.TestData.ServiceArea)

	suite.NoError(err)

	queryResult := domain.ServiceArea{}
	suite.TestDb.Raw("SELECT * FROM public.service_areas WHERE id=?",
		suite.TestData.ServiceArea.ID).Scan(&queryResult)

	suite.EqualValues(queryResult.Identifier, suite.TestData.ServiceArea.Identifier)
}

func (suite *ServiceAreaRepositoryTestSuite) TestRepository_UpdateServiceArea() {
	updated := suite.TestData.ServiceArea
	updated.Identifier = "new-area"

	err := suite.TestRepo.SaveOrUpdateServiceArea(updated)

	suite.NoError(err)

	queryResult := domain.ServiceArea{}
	suite.TestDb.Raw("SELECT * FROM public.service_areas WHERE id=?",
		suite.TestData.ServiceArea.ID).Scan(&queryResult)

	suite.EqualValues(queryResult.Identifier, updated.Identifier)
}

func TestIntegration_ServiceAreaRepositoryTestSuite(t *testing.T) {
	repoSuite := new(ServiceAreaRepositoryTestSuite)
	suite.Run(t, repoSuite)
}
