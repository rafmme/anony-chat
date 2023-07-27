package domain

type UserRepository interface {
	// FindByID retrieves a user by its ID
	FindByID(id string) *User

	// FindByEmail retrieves a user by its email
	FindByEmail(email string) *User

	//Find User
	Find(...interface{}) *User

	// Save saves a user to the repository
	Save(user interface{}) *User

	// DeleteByID deletes a user by its ID
	DeleteByID(ID string) error
}
