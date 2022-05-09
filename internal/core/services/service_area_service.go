package services

import (
	"rider-service/internal/core/domain"
	"rider-service/internal/core/interfaces"
)

type serviceAreaService struct {
	serviceAreaRepository interfaces.ServiceAreaRepository
}

func NewServiceAreaService(serviceAreaRepository interfaces.ServiceAreaRepository) *serviceAreaService {
	return &serviceAreaService{
		serviceAreaRepository: serviceAreaRepository,
	}
}

func (s *serviceAreaService) SaveOrUpdateServiceArea(serviceArea domain.ServiceArea) error {
	return s.serviceAreaRepository.SaveOrUpdateServiceArea(serviceArea)
}
