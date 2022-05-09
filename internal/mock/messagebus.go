package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"rider-service/internal/core/domain"
)

type MessageBusPublisher struct {
	mock.Mock
}

func (m *MessageBusPublisher) CreateRider(ctx context.Context, rider domain.Rider) error {
	args := m.Called(rider)
	return args.Error(0)
}

func (m *MessageBusPublisher) UpdateRider(ctx context.Context, rider domain.Rider) error {
	args := m.Called(rider)
	return args.Error(0)
}

func (m *MessageBusPublisher) UpdateRiderLocation(ctx context.Context, serviceArea domain.ServiceArea, id string, newLocation domain.Location) error {
	args := m.Called(serviceArea, id, newLocation)
	return args.Error(0)
}
