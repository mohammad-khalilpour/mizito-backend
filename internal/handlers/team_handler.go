package handlers

import (
	"fmt"
	"mizito/internal/database"
	"mizito/internal/repositories"
	"mizito/pkg/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TeamHandler interface {
	GetTeams(ctx *fiber.Ctx) error
	GetTeamByID(ctx *fiber.Ctx) error
	GetProjectsByTeam(ctx *fiber.Ctx) error
	AddUsersToTeam(ctx *fiber.Ctx) error
	DeleteUsersFromTeam(ctx *fiber.Ctx) error
	CreateTeam(ctx *fiber.Ctx) error
	UpdateTeam(ctx *fiber.Ctx) error
	DeleteTeam(ctx *fiber.Ctx) error
}

type teamHandler struct {
	repo repositories.TeamRepository
}

func NewTeamHandler(postgreSql *database.DatabaseHandler) TeamHandler {
	repo := repositories.NewTeamRepository(postgreSql)
	return &teamHandler{
		repo: repo,
	}
}

// GetTeams retrieves all teams for the authenticated user
func (h *teamHandler) GetTeams(ctx *fiber.Ctx) error {
	// Assume userID is obtained from middleware (e.g., JWT)
	userID := ctx.Locals("userID").(uint)

	teams, err := h.repo.GetTeams(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve teams",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(teams)
}

// CreateTeam creates a new team
func (h *teamHandler) CreateTeam(ctx *fiber.Ctx) error {
	var team models.Team

	// Parse request body
	if err := ctx.BodyParser(&team); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if team.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Team name is required",
		})
	}

	// Optionally, handle initial members
	// Example: team.Members = [...] (already handled in repository)

	teamID, err := h.repo.CreateTeam(&team)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create team",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"team_id": teamID,
	})
}

// GetTeamByID retrieves a team by its ID
func (h *teamHandler) GetTeamByID(ctx *fiber.Ctx) error {
	teamIDParam := ctx.Params("id")
	teamID, err := strconv.ParseUint(teamIDParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	team, err := h.repo.GetTeamByID(uint(teamID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve team",
		})
	}

	if team == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Team not found",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(team)
}

// GetProjectsByTeam retrieves all projects for a specific team
func (h *teamHandler) GetProjectsByTeam(ctx *fiber.Ctx) error {
	teamIDParam := ctx.Params("id")
	teamID, err := strconv.ParseUint(teamIDParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	projects, err := h.repo.GetProjectsByTeam(uint(teamID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve projects",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(projects)
}

// AddUsersToTeam adds users to a team with a specified role
func (h *teamHandler) AddUsersToTeam(ctx *fiber.Ctx) error {
	teamIDParam := ctx.Params("id")
	teamID, err := strconv.ParseUint(teamIDParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	var payload struct {
		UserIDs []uint      `json:"user_ids" validate:"required,min=1"`
		Role    models.Role `json:"role" validate:"required,oneof=admin member"`
	}

	// Parse request body
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate payload
	if len(payload.UserIDs) == 0 || payload.Role == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User IDs and role are required",
		})
	}

	addedCount, err := h.repo.AddUsersToTeam(payload.UserIDs, uint(teamID), payload.Role)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to add users to team: %v", err),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"added_users": addedCount,
	})
}

// DeleteUsersFromTeam removes users from a team
func (h *teamHandler) DeleteUsersFromTeam(ctx *fiber.Ctx) error {
	teamIDParam := ctx.Params("id")
	teamID, err := strconv.ParseUint(teamIDParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	var payload struct {
		UserIDs []uint `json:"user_ids" validate:"required,min=1"`
	}

	// Parse request body
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate payload
	if len(payload.UserIDs) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User IDs are required",
		})
	}

	deletedCount, err := h.repo.DeleteUsersFromTeam(payload.UserIDs, uint(teamID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to delete users from team: %v", err),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"deleted_users": deletedCount,
	})
}

// UpdateTeam updates a team's details
func (h *teamHandler) UpdateTeam(ctx *fiber.Ctx) error {
	teamIDParam := ctx.Params("id")
	teamID, err := strconv.ParseUint(teamIDParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	var updateData struct {
		Name string `json:"name,omitempty"`
		// Add other fields as necessary
	}

	// Parse request body
	if err := ctx.BodyParser(&updateData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Fetch the existing team
	team, err := h.repo.GetTeamByID(uint(teamID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve team",
		})
	}

	if team == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Team not found",
		})
	}

	// Update fields if provided
	if updateData.Name != "" {
		team.Name = updateData.Name
	}
	// Add other fields as necessary

	updatedTeamID, err := h.repo.UpdateTeam(team)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update team",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"team_id": updatedTeamID,
	})
}

// DeleteTeam deletes a team by its ID
func (h *teamHandler) DeleteTeam(ctx *fiber.Ctx) error {
	teamIDParam := ctx.Params("id")
	teamID, err := strconv.ParseUint(teamIDParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	deletedTeamID, err := h.repo.DeleteTeam(uint(teamID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to delete team: %v", err),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Team deleted successfully",
		"team_id": deletedTeamID,
	})
}
