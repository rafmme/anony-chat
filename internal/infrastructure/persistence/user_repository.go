package infrastructure

import (
	domain "github.com/rafmme/anony-chat/pkg/domain/user"
)

type UserRepository struct {
	users map[string]*domain.User
}

func NewUserRepository() *UserRepository {
	return nil /* &UserRepository{
		users: make(map[string]*domain.User),
	} */
}

func (r *UserRepository) FindByID(ID string) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepository) Save(user *domain.User) error {
	return nil
}

func (r *UserRepository) DeleteByID(ID string) error {
	return nil
}
