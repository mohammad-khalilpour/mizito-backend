package repositories

import "mizito/pkg/models"

type UserRepository interface {
	GetUsers() ([]models.User, error)
	CreateUser(user *models.User) (uint, error)
	GetUserByID(userID uint) (*models.User, error)
	UpdateUser(user *models.User) (uint, error)
	DeleteUser(userID uint) (uint, error)
	GetUserMessages(userID uint) ([]models.Message, error)
}

type userRepository struct {
}

func NewUserRepository() TaskRepository {
	return &taskRepository{}
}

func (tr *taskRepository) GetUsers() ([]models.User, error) {
	return nil, nil
}
func (tr *taskRepository) CreateUser(user *models.User) (uint, error) {
	return 0, nil
}
func (tr *taskRepository) GetUserByID(userID uint) (*models.User, error) {
	return nil, nil
}
func (tr *taskRepository) UpdateUser(user *models.User) (uint, error) {
	return 0, nil
}
func (tr *taskRepository) DeleteUser(userID uint) (uint, error) {
	return 0, nil
}
func (tr *taskRepository) GetUserMessages(userID uint) ([]models.Message, error) {
	return nil, nil
}
