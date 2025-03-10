package handlers

import (
	"github.com/gofiber/fiber/v2"
	"mizito/internal/database"
	"mizito/internal/repositories"
	"mizito/pkg/models"
	"strconv"
)

type ProjectCrudHandler interface {
	GetProjectByID(ctx *fiber.Ctx) error
	UpdateProject(ctx *fiber.Ctx) error
	DeleteProject(ctx *fiber.Ctx) error
	CreateProject(ctx *fiber.Ctx) error
}

type ProjectHandler interface {
	ProjectCrudHandler
	GetProjectsByUser(ctx *fiber.Ctx) error
	GetUsersByProjectID(ctx *fiber.Ctx) error
	AddUserToProject(ctx *fiber.Ctx) error
}

func NewProjectHandler(db *database.DatabaseHandler) ProjectHandler {
	repo := repositories.NewProjectRepository(db)
	return &projectHandler{repository: repo}
}

type projectHandler struct {
	repository repositories.ProjectRepository
}

func (pr *projectHandler) GetProjectsByUser(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	projects, repoErr := pr.repository.GetProjectsByUser(userID)
	if repoErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": repoErr.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(projects)
}

func (pr *projectHandler) GetUsersByProjectID(ctx *fiber.Ctx) error {
	projectID := ctx.Params("project_id")
	requestUserID := ctx.Locals("userID").(uint)
	if projectID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "projectID is required",
		})
	}

	parsedProjectID, err := strconv.ParseUint(projectID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid projectID format",
		})
	}

	users, repoErr := pr.repository.GetUsersByProjectID(uint(parsedProjectID), requestUserID)
	if repoErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": repoErr.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(users)
}

func (pr *projectHandler) CreateProject(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)
	project := new(models.Project)
	if err := ctx.BodyParser(project); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	projectID, repoErr := pr.repository.CreateProject(project, userID)
	if repoErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": repoErr.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"projectID": projectID,
	})
}

func (pr *projectHandler) GetProjectByID(ctx *fiber.Ctx) error {
	projectID := ctx.Params("project_id")
	requestUserID := ctx.Locals("userID").(uint)
	if projectID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "projectID is required",
		})
	}

	parsedProjectID, err := strconv.ParseUint(projectID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid projectID format",
		})
	}

	project, repoErr := pr.repository.GetProjectByID(uint(parsedProjectID), requestUserID)
	if repoErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": repoErr.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(project)
}

func (pr *projectHandler) UpdateProject(ctx *fiber.Ctx) error {
	projectID := ctx.Params("projectID")
	requestUserID := ctx.Locals("userID").(uint)
	if projectID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "projectID is required",
		})
	}

	parsedProjectID, err := strconv.ParseUint(projectID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid projectID format",
		})
	}

	updatedProject := new(models.Project)
	if err := ctx.BodyParser(updatedProject); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	_, repoErr := pr.repository.UpdateProject(uint(parsedProjectID), updatedProject, requestUserID)
	if repoErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": repoErr.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Project updated successfully",
	})
}

func (pr *projectHandler) DeleteProject(ctx *fiber.Ctx) error {
	projectID := ctx.Params("projectID")
	requestUserID := ctx.Locals("userID").(uint)

	if projectID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "projectID is required",
		})
	}

	parsedProjectID, err := strconv.ParseUint(projectID, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid projectID format",
		})
	}

	_, repoErr := pr.repository.DeleteProject(uint(parsedProjectID), requestUserID)
	if repoErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": repoErr.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Project deleted successfully",
	})
}

func (pr *projectHandler) AddUserToProject(ctx *fiber.Ctx) error {
	type RequestBody struct {
		UserID uint `json:"userID"`
	}
	requestUserID := ctx.Locals("userID").(uint)

	projectID, err := ctx.ParamsInt("project_id")
	if err != nil || projectID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	var requestBody RequestBody
	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	userID := requestBody.UserID
	if userID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID cannot be zero",
		})
	}

	repoErr := pr.repository.AddUserToProject(uint(projectID), userID, requestUserID)
	if repoErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": repoErr.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User added to the project successfully",
	})
}
