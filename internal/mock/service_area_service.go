package mock

import (
	"github.com/stretchr/testify/mock"
	"rider-service/internal/core/domain"
)

type ServiceAreaService struct {
	mock.Mock
}

func (m *ServiceAreaService) SaveOrUpdateServiceArea(serviceArea domain.ServiceArea) error {
	args := m.Called(serviceArea)
	return args.Error(0)
}
