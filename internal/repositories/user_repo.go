package repositories

import (
	"errors"

	"mizito/internal/database"
	"mizito/pkg/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) (uint, error)
	GetUserByID(userID uint) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUser(user *models.User) (uint, error)
	DeleteUser(userID uint) (uint, error)
}

type userRepository struct {
	db *database.DatabaseHandler
}

func NewUserRepository(db *database.DatabaseHandler) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) CreateUser(user *models.User) (uint, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword

	if err := ur.db.DB.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (ur *userRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := ur.db.DB.First(&user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := ur.db.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *userRepository) UpdateUser(user *models.User) (uint, error) {
	allowedUpdates := map[string]interface{}{
		"Username": user.Username,
		"Email":    user.Email,
	}

	if err := ur.db.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(allowedUpdates).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (ur *userRepository) DeleteUser(userID uint) (uint, error) {
	result := ur.db.DB.Delete(&models.User{}, userID)
	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, nil
	}

	return userID, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
