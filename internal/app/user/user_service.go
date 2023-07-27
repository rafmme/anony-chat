package application

import (
	domain "github.com/rafmme/anony-chat/pkg/domain/user"
)

type UserService struct {
	UserRepository domain.UserRepository
}

func (userService *UserService) CreateUser(user interface{}) (*domain.User, error) {
	newUserData, err := userService.UserRepository.Save(user)

	if newUserData != nil {
		return newUserData, nil
	}

	return nil, err
}

func (userService *UserService) FetchUser(where ...interface{}) (*domain.User, error) {
	existingUserData, err := userService.UserRepository.Find(where)

	if existingUserData != nil {
		return existingUserData, nil
	}

	return nil, err
}

func (userService *UserService) FetchUserByID(id string) (*domain.User, error) {
	existingUserData, err := userService.UserRepository.FindByID(id)

	if existingUserData != nil {
		return existingUserData, nil
	}

	return nil, err
}

func (userService *UserService) FetchUserByEmail(email string) (*domain.User, error) {
	existingUserData, err := userService.UserRepository.FindByEmail(email)

	if existingUserData != nil {
		return existingUserData, nil
	}

	return nil, err
}
