package dto

import "rider-service/internal/core/domain"

type BodyUpdate struct {
	Status      int
	ServiceArea int
	Capacity    CreateDimensions
}

type ResponseUpdate domain.Rider

func BuildResponseUpdate(model domain.Rider) ResponseCreate {
	return ResponseCreate(model)
}
