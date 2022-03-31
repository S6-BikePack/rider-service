package repositories

import (
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rider-service/internal/core/domain"
)

type cockroachdb struct {
	Connection *gorm.DB
}

func NewCockroachDB(connStr string) (*cockroachdb, error) {
	db, err := gorm.Open(postgres.Open(connStr))

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.Rider{})

	if err != nil {
		return nil, err
	}

	database := cockroachdb{
		Connection: db,
	}

	return &database, nil
}

func (repository *cockroachdb) Get(uid uuid.UUID) (domain.Rider, error) {
	var rider domain.Rider

	repository.Connection.Preload(clause.Associations).First(&rider, uid)

	return rider, nil
}

func (repository *cockroachdb) GetAll() ([]domain.Rider, error) {
	var riders []domain.Rider

	repository.Connection.Find(&riders)

	return riders, nil
}

func (repository *cockroachdb) Save(rider domain.Rider) (domain.Rider, error) {
	result := repository.Connection.Create(&rider)

	if result.Error != nil {
		return domain.Rider{}, result.Error
	}

	return rider, nil
}

func (repository *cockroachdb) Update(rider domain.Rider) (domain.Rider, error) {
	result := repository.Connection.Model(&rider).Updates(rider)

	if result.Error != nil {
		return domain.Rider{}, result.Error
	}

	return rider, nil
}
