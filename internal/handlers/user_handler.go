package handlers

import (
	"mizito/internal/database"
	"strconv"

	"mizito/internal/repositories"
	"mizito/pkg/models"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	GetUsers(ctx *fiber.Ctx) error
	CreateUser(ctx *fiber.Ctx) error
	GetUserByID(ctx *fiber.Ctx) error
	UpdateUser(ctx *fiber.Ctx) error
	DeleteUser(ctx *fiber.Ctx) error
}

type userHandler struct {
	repo repositories.UserRepository
}

func NewUserHandler(postgreSql *database.DatabaseHandler) UserHandler {
	repo := repositories.NewUserRepository(postgreSql)
	return &userHandler{
		repo: repo,
	}
}

func (h *userHandler) GetUsers(ctx *fiber.Ctx) error {
	users, err := h.repo.GetAllUsers()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(users)
}

func (h *userHandler) CreateUser(ctx *fiber.Ctx) error {
	var user models.User

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, Email, and Password are required",
		})
	}

	userID, err := h.repo.CreateUser(&user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_id": userID,
	})
}

func (h *userHandler) GetUserByID(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.repo.GetUserByID(uint(userID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}

	if user == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	user.Password = ""

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (h *userHandler) UpdateUser(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var updateData struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
	}

	if err := ctx.BodyParser(&updateData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.repo.GetUserByID(uint(userID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}

	if user == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if updateData.Username != "" {
		user.Username = updateData.Username
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}

	updatedUserID, err := h.repo.UpdateUser(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id": updatedUserID,
	})
}

func (h *userHandler) DeleteUser(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	deletedUserID, err := h.repo.DeleteUser(uint(userID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	if deletedUserID == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
