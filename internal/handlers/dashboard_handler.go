package handlers

import (
	"github.com/gofiber/fiber/v2"
	"mizito/internal/database"
	"mizito/internal/repositories"
)

type DashboardHandler struct {
	Repo repositories.DashboardRepository
}

func NewDashboardHandler(postgreSql *database.DatabaseHandler) *DashboardHandler {
	repo := repositories.NewDashboardRepository(postgreSql)
	return &DashboardHandler{Repo: repo}
}

func (h *DashboardHandler) GetDashboardDetails(ctx *fiber.Ctx) error {
	// Extract userID from Locals
	requestUserID, ok := ctx.Locals("userID").(uint)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Call the repository method
	username, profilePicture, todoCount, coworkers, todoList, projectList, err := h.Repo.GetDashboardDetails(requestUserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Construct the response
	response := fiber.Map{
		"username":        username,
		"profile_picture": profilePicture,
		"todo_count":      todoCount,
		"coworkers":       coworkers,
		"todo_list":       todoList,
		"project_list":    projectList,
	}

	// Return the response
	return ctx.Status(fiber.StatusOK).JSON(response)
}
