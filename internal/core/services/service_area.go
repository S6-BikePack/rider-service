package services

import (
	"rider-service/internal/core/domain"
	"rider-service/internal/core/ports"
)

type serviceAreaService struct {
	serviceAreaRepository ports.ServiceAreaRepository
}

func NewServiceAreaService(serviceAreaRepository ports.ServiceAreaRepository) *serviceAreaService {
	return &serviceAreaService{
		serviceAreaRepository: serviceAreaRepository,
	}
}

func (s *serviceAreaService) SaveOrUpdateServiceArea(serviceArea domain.ServiceArea) error {
	return s.serviceAreaRepository.SaveOrUpdateServiceArea(serviceArea)
}
