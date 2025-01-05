package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type ProjectCrudRepository interface {
	GetProjectByID(ctx *fiber.Ctx) error
	UpdateProject(ctx *fiber.Ctx) error
	DeleteProject(ctx *fiber.Ctx) error
	CreateProject(ctx *fiber.Ctx) error
}

type ProjectRepository interface {
	ProjectCrudRepository
	GetProjectsByUser(ctx *fiber.Ctx) error
	AddUserToProject(ctx *fiber.Ctx) error
	AssignTask(ctx *fiber.Ctx) error
}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

type projectRepository struct {
}

func (pr *projectRepository) GetProjectsByUser(ctx *fiber.Ctx) error {
	return nil
}
func (pr *projectRepository) CreateProject(ctx *fiber.Ctx) error {
	return nil
}
func (pr *projectRepository) GetProjectByID(ctx *fiber.Ctx) error {
	return nil
}
func (pr *projectRepository) UpdateProject(ctx *fiber.Ctx) error {
	return nil
}
func (pr *projectRepository) DeleteProject(ctx *fiber.Ctx) error {
	return nil
}
func (pr *projectRepository) AddUserToProject(ctx *fiber.Ctx) error {
	return nil
}
func (pr *projectRepository) AssignTask(ctx *fiber.Ctx) error {
	return nil
}
