package dto

type CreateDimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Depth  int `json:"depth"`
}

type BodyCreateRider struct {
	ID          string           `json:"id"`
	ServiceArea int              `json:"serviceArea"`
	Capacity    CreateDimensions `json:"capacity"`
	Status      int              `json:"status"`
}