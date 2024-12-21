package http_handlers

import "github.com/gofiber/fiber/v2"



type UserRepository interface {
	GetUsers(ctx *fiber.Ctx) error
	CreateUser(ctx *fiber.Ctx) error
	GetUserByID(ctx *fiber.Ctx) error
	UpdateUser(ctx *fiber.Ctx) error
	DeleteUser(ctx *fiber.Ctx) error
}


type userRepository struct {
	
}

func NewUserRepository() UserRepository{
	return &userRepository{}
}


func (pr *userRepository) GetUsers(ctx *fiber.Ctx) error {
	return nil
}
func (pr *userRepository) CreateUser(ctx *fiber.Ctx) error {
	return nil
}
func (pr *userRepository) GetUserByID(ctx *fiber.Ctx) error {
	return nil
}
func (pr *userRepository) UpdateUser(ctx *fiber.Ctx) error {
	return nil
}
func (pr *userRepository) DeleteUser(ctx *fiber.Ctx) error {
	return nil
}


