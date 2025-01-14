package basic

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mizito/internal/database"
	"mizito/pkg/models"
)

type BasicRepository interface {
	AuthenticateUser(username string, password string) (bool, uint, error)
}

func NewBasicHandler(dbHandler *database.DatabaseHandler) BasicRepository {
	return &basicRepository{
		db: dbHandler.DB,
	}
}

type basicRepository struct {
	db *gorm.DB
}

func (jr *basicRepository) AuthenticateUser(username string, password string) (bool, uint, error) {
	var user models.User
	err := jr.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, 0, nil
	} else if err != nil {
		return false, 0, fmt.Errorf("failed to retrieve user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, 0, nil
	}
	return true, user.ID, nil
}
