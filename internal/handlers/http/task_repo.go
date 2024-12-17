package http_handlers

import (
	"github.com/gofiber/fiber/v2"
)


type TaskRepository interface {
	GetTasksByProject(ctx *fiber.Ctx) error
	CreateTask(ctx *fiber.Ctx) error
	GetTaskByID(ctx *fiber.Ctx) error
	UpdateTask(ctx *fiber.Ctx) error
	DeleteTask(ctx *fiber.Ctx) error
}

func NewTaskRepository() TaskRepository{
	return nil
}