package dto

import "rider-service/internal/core/domain"

type ridersResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Status        int    `json:"status"`
	ServiceAreaID int    `json:"serviceArea"`
}

func createRidersResponse(rider domain.Rider) ridersResponse {
	return ridersResponse{
		ID:            rider.UserID,
		Name:          rider.User.Name,
		Status:        rider.Status,
		ServiceAreaID: rider.ServiceArea.ID,
	}
}

type RiderListResponse []*ridersResponse

func CreateServiceAreaListResponse(riders []domain.Rider) RiderListResponse {
	response := RiderListResponse{}
	for _, r := range riders {
		rider := createRidersResponse(r)
		response = append(response, &rider)
	}
	return response
}
