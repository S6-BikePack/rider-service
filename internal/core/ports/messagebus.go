package ports

import (
	"context"
	"rider-service/internal/core/domain"
)

type MessageBusPublisher interface {
	CreateRider(ctx context.Context, rider domain.Rider) error
	UpdateRider(ctx context.Context, rider domain.Rider) error
	UpdateRiderLocation(ctx context.Context, serviceArea domain.ServiceArea, id string, newLocation domain.Location) error
}
