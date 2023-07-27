package infrastructure

import (
	"time"

	domain "github.com/rafmme/anony-chat/pkg/domain/user"
	"github.com/rafmme/anony-chat/pkg/shared"
)

var UserRepo domain.UserRepository

type UserRepository struct {
}

func init() {
	UserRepo = new(UserRepository)
}

func (r *UserRepository) FindByID(id string) *domain.User {
	var fetchedUser domain.User
	shared.Database.First(&fetchedUser, "id = ?", id)
	return &fetchedUser
}

func (r *UserRepository) FindByEmail(email string) *domain.User {
	var fetchedUser domain.User
	shared.Database.First(&fetchedUser, "email = ?", email)
	return &fetchedUser
}

func (r *UserRepository) Find(where ...interface{}) *domain.User {
	var fetchedUser domain.User
	shared.Database.First(
		&fetchedUser, where...,
	)

	return &fetchedUser
}

func (r *UserRepository) Save(userData interface{}) *domain.User {
	time := time.Now().String()
	user := &domain.User{
		ID:        shared.CreateUUID(),
		Email:     userData.(*shared.UserSignupData).Email,
		Password:  userData.(*shared.UserSignupData).Password,
		CreatedAt: time,
		UpdatedAt: time,
	}

	shared.Database.Create(user)
	return user
}

func (r *UserRepository) DeleteByID(ID string) error {
	return nil
}
