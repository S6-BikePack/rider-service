package repositories

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rider-service/internal/core/domain"
)

type cockroachdb struct {
	Connection *gorm.DB
}

func NewCockroachDB(db *gorm.DB) (*cockroachdb, error) {
	err := db.AutoMigrate(&domain.Rider{})

	if err != nil {
		return nil, err
	}

	database := cockroachdb{
		Connection: db,
	}

	return &database, nil
}

func (repository *cockroachdb) Get(id string) (domain.Rider, error) {
	var rider domain.Rider

	repository.Connection.Preload(clause.Associations).First(&rider, "user_id = ?", id)

	return rider, nil
}

func (repository *cockroachdb) GetAll() ([]domain.Rider, error) {
	var riders []domain.Rider

	repository.Connection.Find(&riders)

	return riders, nil
}

func (repository *cockroachdb) Save(rider domain.Rider) (domain.Rider, error) {
	result := repository.Connection.Omit("User").Create(&rider)

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

func (repository *cockroachdb) SaveOrUpdateUser(user domain.User) error {
	updateResult := repository.Connection.Model(&user).Where("id = ?", user.ID).Updates(&user)

	if updateResult.RowsAffected == 0 {
		createResult := repository.Connection.Create(&user)

		if createResult.Error != nil {
			return errors.New("could not create user")
		}
	}

	if updateResult.Error != nil {
		return errors.New("could not update user")
	}

	return nil
}

func (repository *cockroachdb) GetUser(id string) (domain.User, error) {
	var user domain.User

	repository.Connection.Preload(clause.Associations).First(&user, "id = ?", id)

	if (user == domain.User{}) {
		return user, errors.New("user not found")
	}

	return user, nil
}
