package ports

import (
	"rider-service/internal/core/domain"
)

type MessageBusPublisher interface {
	CreateRider(rider domain.Rider) error
	UpdateRider(rider domain.Rider) error
	UpdateRiderLocation(id string, newLocation domain.Location) error
}
