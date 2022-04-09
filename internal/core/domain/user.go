package domain

type User struct {
	ID       string
	Name     string
	LastName string
}

func NewUser(id, name, lastName string) User {
	return User{
		ID:       id,
		Name:     name,
		LastName: lastName,
	}
}
