package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"rider-service/internal/core/domain"
)

type RiderService struct {
	mock.Mock
}

func (m *RiderService) GetAll(ctx context.Context) ([]domain.Rider, error) {
	args := m.Called()
	return args.Get(0).([]domain.Rider), args.Error(1)
}

func (m *RiderService) Get(ctx context.Context, id string) (domain.Rider, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Rider), args.Error(1)
}

func (m *RiderService) Create(ctx context.Context, userId string, serviceArea int, capacity domain.Dimensions) (domain.Rider, error) {
	args := m.Called(userId, serviceArea, capacity)
	return args.Get(0).(domain.Rider), args.Error(1)
}

func (m *RiderService) Update(ctx context.Context, id string, status int, serviceArea int, capacity domain.Dimensions) (domain.Rider, error) {
	args := m.Called(id, status, serviceArea, capacity)
	return args.Get(0).(domain.Rider), args.Error(1)
}

func (m *RiderService) UpdateLocation(ctx context.Context, id string, location domain.Location) (domain.Rider, error) {
	args := m.Called(id, location)
	return args.Get(0).(domain.Rider), args.Error(1)
}

func (m *RiderService) SaveOrUpdateUser(ctx context.Context, user domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}
