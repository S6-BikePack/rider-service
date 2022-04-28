package ports

import (
	"context"
	"rider-service/internal/core/domain"
)

type RiderService interface {
	GetAll(ctx context.Context) ([]domain.Rider, error)
	Get(ctx context.Context, id string) (domain.Rider, error)
	Create(ctx context.Context, userId string, serviceArea int, capacity domain.Dimensions) (domain.Rider, error)
	Update(ctx context.Context, id string, status int, serviceArea int, capacity domain.Dimensions) (domain.Rider, error)
	UpdateLocation(ctx context.Context, id string, location domain.Location) (domain.Rider, error)
	SaveOrUpdateUser(ctx context.Context, user domain.User) error
}
