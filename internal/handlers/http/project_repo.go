package http_handlers

import (
	"github.com/gofiber/fiber/v2"
)


type ProjectRepository interface {
	GetProjectsByUser(ctx *fiber.Ctx) error
	CreateProject(ctx *fiber.Ctx) error
	GetProjectByID(ctx *fiber.Ctx)	error
	UpdateProject(ctx *fiber.Ctx) error
	DeleteProject(ctx *fiber.Ctx) error
}

func NewProjectRepository() ProjectRepository{
	return nil
}