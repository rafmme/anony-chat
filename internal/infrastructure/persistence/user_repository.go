package infrastructure

import (
	domain "github.com/rafmme/anony-chat/pkg/domain/user"
	"github.com/rafmme/anony-chat/pkg/shared"
)

type UserRepository struct {
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	var fetchedUser domain.User
	result := shared.Database.First(&fetchedUser, "id = ?", id)

	if err := result.Error; err != nil {
		return nil, err
	}

	return &fetchedUser, nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var fetchedUser domain.User
	result := shared.Database.First(&fetchedUser, "email = ?", email)

	if err := result.Error; err != nil {
		return nil, err
	}

	return &fetchedUser, nil
}

func (r *UserRepository) Find(where ...interface{}) (*domain.User, error) {
	var fetchedUser domain.User
	result := shared.Database.First(
		&fetchedUser, where...,
	)

	if err := result.Error; err != nil {
		return nil, err
	}

	return &fetchedUser, nil
}

func (r *UserRepository) Save(userData interface{}) (*domain.User, error) {
	user := &domain.User{
		ID:       shared.CreateUUID(),
		Email:    userData.(*shared.UserSignupData).Email,
		Password: userData.(*shared.UserSignupData).Password,
	}

	result := shared.Database.Create(user)

	if err := result.Error; err != nil {
		return nil, err
	}

	return result.Value.(*domain.User), nil
}

func (r *UserRepository) DeleteByID(ID string) error {
	return nil
}
