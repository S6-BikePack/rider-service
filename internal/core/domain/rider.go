package domain

type Rider struct {
	UserID      string `gorm:"primaryKey"`
	User        User
	Status      int
	ServiceArea int
	Capacity    Dimensions `gorm:"embedded"`
	Location    Location
}

func NewRider(user User, status int, serviceArea int, capacity Dimensions) Rider {
	return Rider{
		UserID:      user.ID,
		User:        user,
		Status:      status,
		ServiceArea: serviceArea,
		Capacity:    capacity,
	}
}
