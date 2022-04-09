package domain

type Rider struct {
	UserID   string `gorm:"primaryKey"`
	User     User
	Status   int8
	Location Location `gorm:"embedded"`
}

func NewRider(user User, status int8, location Location) Rider {
	return Rider{
		UserID:   user.ID,
		User:     user,
		Status:   status,
		Location: location,
	}
}
