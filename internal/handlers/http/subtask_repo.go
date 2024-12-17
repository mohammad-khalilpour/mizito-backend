package http_handlers

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


func NewSubtaskRepository() SubtaskRepository{
	return nil
}