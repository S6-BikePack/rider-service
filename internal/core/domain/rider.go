package domain

type Rider struct {
	UserID        string `gorm:"primaryKey"`
	User          User
	Status        int
	ServiceAreaID int `json:"-"`
	ServiceArea   ServiceArea
	Capacity      Dimensions `gorm:"embedded"`
	Location      Location
}

func NewRider(user User, status int, serviceArea int, capacity Dimensions) Rider {
	return Rider{
		UserID:        user.ID,
		User:          user,
		Status:        status,
		ServiceAreaID: serviceArea,
		Capacity:      capacity,
	}
}
