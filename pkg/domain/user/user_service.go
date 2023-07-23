package domain

type UserService interface {
	CreateUser(id, email, password string) (*User, error)
}
