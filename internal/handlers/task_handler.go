package handlers

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

type taskRepository struct {
}

func NewTaskRepository() TaskRepository {
	return &taskRepository{}
}

func (pr *taskRepository) GetTasksByProject(ctx *fiber.Ctx) error {
	return nil
}
func (pr *taskRepository) CreateTask(ctx *fiber.Ctx) error {
	return nil
}
func (pr *taskRepository) GetTaskByID(ctx *fiber.Ctx) error {
	return nil
}
func (pr *taskRepository) UpdateTask(ctx *fiber.Ctx) error {
	return nil
}
func (pr *taskRepository) DeleteTask(ctx *fiber.Ctx) error {
	return nil
}
