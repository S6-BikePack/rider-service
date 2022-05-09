package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"rider-service/internal/core/domain"
)

type RiderRepository struct {
	mock.Mock
}

func (m *RiderRepository) GetAll(ctx context.Context) ([]domain.Rider, error) {
	args := m.Called()
	return args.Get(0).([]domain.Rider), args.Error(1)
}

func (m *RiderRepository) Get(ctx context.Context, id string) (domain.Rider, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Rider), args.Error(1)
}

func (m *RiderRepository) Save(ctx context.Context, rider domain.Rider) (domain.Rider, error) {
	args := m.Called(rider)
	return args.Get(0).(domain.Rider), args.Error(1)
}

func (m *RiderRepository) Update(ctx context.Context, rider domain.Rider) (domain.Rider, error) {
	args := m.Called(rider)
	return args.Get(0).(domain.Rider), args.Error(1)
}

func (m *RiderRepository) SaveOrUpdateUser(ctx context.Context, user domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *RiderRepository) GetUser(ctx context.Context, id string) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}
