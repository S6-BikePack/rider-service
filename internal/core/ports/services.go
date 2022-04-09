package ports

import (
	"rider-service/internal/core/domain"
)

type RiderService interface {
	GetAll() ([]domain.Rider, error)
	Get(id string) (domain.Rider, error)
	Create(userId string, status int8) (domain.Rider, error)
	Update(id string, status int8) (domain.Rider, error)
	UpdateLocation(id string, location domain.Location) (domain.Rider, error)
	SaveOrUpdateUser(user domain.User) error
}
