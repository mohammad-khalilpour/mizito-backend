package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type SubtaskRepository interface {
	GetSubtasksByTask(ctx *fiber.Ctx) error
	CreateSubtask(ctx *fiber.Ctx) error
	GetSubtaskByID(ctx *fiber.Ctx) error
	UpdateSubtask(ctx *fiber.Ctx) error
	DeleteSubtask(ctx *fiber.Ctx) error
}

type subtaskRepository struct {
}

func NewSubtaskRepository() SubtaskRepository {
	return &subtaskRepository{}
}

func (pr *subtaskRepository) GetSubtasksByTask(ctx *fiber.Ctx) error {
	return nil
}
func (pr *subtaskRepository) CreateSubtask(ctx *fiber.Ctx) error {
	return nil
}
func (pr *subtaskRepository) GetSubtaskByID(ctx *fiber.Ctx) error {
	return nil
}
func (pr *subtaskRepository) UpdateSubtask(ctx *fiber.Ctx) error {
	return nil
}
func (pr *subtaskRepository) DeleteSubtask(ctx *fiber.Ctx) error {
	return nil
}
