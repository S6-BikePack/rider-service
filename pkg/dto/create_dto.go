package dto

import "rider-service/internal/core/domain"

type CreateDimensions struct {
	Width  int
	Height int
	Depth  int
}

type BodyCreate struct {
	ID          string
	ServiceArea int
	Capacity    CreateDimensions
}

type ResponseCreate domain.Rider

func BuildResponseCreate(model domain.Rider) ResponseCreate {
	return ResponseCreate(model)
}
