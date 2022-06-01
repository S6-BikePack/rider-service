package repositories

import (
	"gorm.io/gorm"
	"rider-service/internal/core/domain"
)

type serviceAreaRepository struct {
	Connection *gorm.DB
}

func NewServiceAreaRepository(db *gorm.DB) (*serviceAreaRepository, error) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"postgis\";")
	
	err := db.AutoMigrate(&domain.ServiceArea{})

	if err != nil {
		return nil, err
	}

	database := serviceAreaRepository{
		Connection: db,
	}

	return &database, nil
}

func (repository *serviceAreaRepository) SaveOrUpdateServiceArea(serviceArea domain.ServiceArea) error {
	if repository.Connection.Model(&serviceArea).Where("id = ?", serviceArea.ID).Updates(&serviceArea).RowsAffected == 0 {
		create := repository.Connection.Create(&serviceArea)

		if create.Error != nil {
			return create.Error
		}
	}

	return nil
}
