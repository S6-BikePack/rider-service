package interfaces

import (
	"context"
	"rider-service/internal/core/domain"
)

type RiderRepository interface {
	GetAll(ctx context.Context) ([]domain.Rider, error)
	Get(ctx context.Context, id string) (domain.Rider, error)
	Save(ctx context.Context, rider domain.Rider) (domain.Rider, error)
	Update(ctx context.Context, rider domain.Rider) (domain.Rider, error)
	SaveOrUpdateUser(ctx context.Context, user domain.User) error
	GetUser(ctx context.Context, id string) (domain.User, error)
}

type ServiceAreaRepository interface {
	SaveOrUpdateServiceArea(serviceArea domain.ServiceArea) error
}
