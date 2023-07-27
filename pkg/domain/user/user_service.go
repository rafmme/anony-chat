package domain

type UserService interface {
	CreateUser(user interface{}) (*User, error)
	FetchUser(...interface{}) (*User, error)
	FetchUserByID(id string) (*User, error)
	FetchUserByEmail(email string) (*User, error)
}
