package repositories

import "mizito/pkg/models"

type UserRepository interface {
	GetUsers() ([]models.User, error)
	CreateUser(user *models.User) (uint, error)
	GetUserByID(userID uint) (models.User, error)
	UpdateUser(user *models.User) (uint, error)
	DeleteUser(userID uint) (uint, error)
}