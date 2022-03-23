package riderrepo

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rider-service/internal/core/domain"
	"strconv"
)

type cockroachdb struct {
	Connection *gorm.DB
}

func NewCockroachDB(connStr string) (*cockroachdb, error) {
	db, err := gorm.Open(postgres.Open(connStr))

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&RiderDao{})

	if err != nil {
		return nil, err
	}

	database := cockroachdb{
		Connection: db,
	}

	return &database, nil
}

func (repository *cockroachdb) Get(id string) (domain.Rider, error) {
	var riderDao RiderDao
	uid, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return domain.Rider{}, nil
	}

	repository.Connection.Preload(clause.Associations).First(&riderDao, uid)

	return riderDao.ToDomain()
}

func (repository *cockroachdb) GetAll() ([]domain.Rider, error) {
	var ridersDao []RiderDao
	repository.Connection.Find(&ridersDao)

	var riders []domain.Rider

	for _, v := range ridersDao {
		domainModel, err := v.ToDomain()

		if err != nil {
			return nil, err
		}

		riders = append(riders, domainModel)
	}

	return riders, nil
}

func (repository *cockroachdb) Save(rider domain.Rider) (domain.Rider, error) {
	var riderDao RiderDao
	riderDao.FromDomain(rider)

	result := repository.Connection.Create(&riderDao)

	if result.Error != nil {
		return domain.Rider{}, result.Error
	}

	return riderDao.ToDomain()
}

func (repository *cockroachdb) Update(rider domain.Rider) (domain.Rider, error) {
	var riderDao RiderDao
	riderDao.FromDomain(rider)

	result := repository.Connection.Model(&riderDao).Updates(riderDao)

	if result.Error != nil {
		return domain.Rider{}, result.Error
	}

	return riderDao.ToDomain()
}
