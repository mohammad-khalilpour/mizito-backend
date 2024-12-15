package handlers

import "github.com/gofiber/fiber/v2"



type UserRepository interface {
	GetUsers(ctx *fiber.Ctx) error
	CreateUser(ctx *fiber.Ctx) error
	GetUserByID(ctx *fiber.Ctx) error
	UpdateUser(ctx *fiber.Ctx) error
	DeleteUser(ctx *fiber.Ctx) error
}


func NewUserRepository() UserRepository{
	return nil
}