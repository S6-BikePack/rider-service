package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rider-service/internal/core/domain"
)

type riderRepository struct {
	Connection *gorm.DB
}

func NewRiderRepository(db *gorm.DB) (*riderRepository, error) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"postgis\";")

	err := db.AutoMigrate(&domain.ServiceArea{}, &domain.Rider{})

	if err != nil {
		return nil, err
	}

	database := riderRepository{
		Connection: db,
	}

	return &database, nil
}

func (repository *riderRepository) Get(ctx context.Context, id string) (domain.Rider, error) {
	var rider domain.Rider

	repository.Connection.WithContext(ctx).Preload(clause.Associations).First(&rider, "user_id = ?", id)

	if (rider == domain.Rider{}) {
		return domain.Rider{}, errors.New("could not find rider")
	}

	return rider, nil
}

func (repository *riderRepository) GetAll(ctx context.Context) ([]domain.Rider, error) {
	var riders []domain.Rider

	repository.Connection.WithContext(ctx).Find(&riders)

	return riders, nil
}

func (repository *riderRepository) Save(ctx context.Context, rider domain.Rider) (domain.Rider, error) {
	result := repository.Connection.WithContext(ctx).Omit("User").Create(&rider)

	if result.Error != nil {
		return domain.Rider{}, result.Error
	}

	return rider, nil
}

func (repository *riderRepository) Update(ctx context.Context, rider domain.Rider) (domain.Rider, error) {
	result := repository.Connection.WithContext(ctx).Model(&rider).Updates(rider)

	if result.Error != nil {
		return domain.Rider{}, result.Error
	}

	return rider, nil
}

func (repository *riderRepository) SaveOrUpdateUser(ctx context.Context, user domain.User) error {
	updateResult := repository.Connection.WithContext(ctx).Model(&user).Where("id = ?", user.ID).Updates(&user)

	if updateResult.RowsAffected == 0 {
		createResult := repository.Connection.WithContext(ctx).Create(&user)

		if createResult.Error != nil {
			return errors.New("could not create user")
		}
	}

	if updateResult.Error != nil {
		return errors.New("could not update user")
	}

	return nil
}

func (repository *riderRepository) GetUser(ctx context.Context, id string) (domain.User, error) {
	var user domain.User

	repository.Connection.WithContext(ctx).Preload(clause.Associations).First(&user, "id = ?", id)

	if (user == domain.User{}) {
		return user, errors.New("user not found")
	}

	return user, nil
}
