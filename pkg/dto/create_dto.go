package dto

import "rider-service/internal/core/domain"

type BodyCreate struct {
	Name   string
	Status int8
}

type ResponseCreate domain.Rider

func BuildResponseCreate(model domain.Rider) ResponseCreate {
	return ResponseCreate(model)
}
