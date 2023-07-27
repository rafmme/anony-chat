package domain

type UserRepository interface {
	// FindByID retrieves a user by its ID
	FindByID(id string) (*User, error)

	// FindByEmail retrieves a user by its email
	FindByEmail(email string) (*User, error)

	//Find User
	Find(...interface{}) (*User, error)

	// Save saves a user to the repository
	Save(user interface{}) (*User, error)

	// DeleteByID deletes a user by its ID
	DeleteByID(ID string) error
}
