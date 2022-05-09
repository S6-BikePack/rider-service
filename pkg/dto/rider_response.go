package dto

import "rider-service/internal/core/domain"

type riderResponseUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
}

type riderResponseArea struct {
	ID         int    `json:"id"`
	Identifier string `json:"identifier"`
}

type riderResponseCapacity struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Depth  int `json:"depth"`
}

type riderResponseLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RiderResponse struct {
	ID          string                `json:"id"`
	User        riderResponseUser     `json:"user"`
	Status      int                   `json:"status"`
	ServiceArea riderResponseArea     `json:"serviceArea"`
	Capacity    riderResponseCapacity `json:"capacity"`
	Location    riderResponseLocation `json:"location"`
}

func CreateRiderResponse(rider domain.Rider) RiderResponse {
	return RiderResponse{
		ID:          rider.UserID,
		User:        riderResponseUser(rider.User),
		Status:      rider.Status,
		ServiceArea: riderResponseArea(rider.ServiceArea),
		Capacity:    riderResponseCapacity(rider.Capacity),
		Location:    riderResponseLocation(rider.Location),
	}
}
