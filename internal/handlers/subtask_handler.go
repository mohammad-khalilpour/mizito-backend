package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"mizito/internal/database"
	"mizito/internal/repositories"
	"mizito/pkg/models"
)

type SubtaskHandler interface {
	GetSubtasksByTask(ctx *fiber.Ctx) error
	CreateSubtask(ctx *fiber.Ctx) error
	GetSubtaskByID(ctx *fiber.Ctx) error
	UpdateSubtask(ctx *fiber.Ctx) error
	DeleteSubtask(ctx *fiber.Ctx) error
}

type subtaskHandler struct {
	repository repositories.SubtaskRepository
}

func NewSubtaskHandler(db *database.DatabaseHandler) SubtaskHandler {
	repo := repositories.NewSubtaskRepository(db)
	return &subtaskHandler{repository: repo}
}

func (sh *subtaskHandler) GetSubtasksByTask(ctx *fiber.Ctx) error {
	taskID, err := ctx.ParamsInt("task_id")
	if err != nil || taskID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	requestUserID := ctx.Locals("userID").(uint)

	subtasks, err := sh.repository.GetSubtasksByTask(uint(taskID), requestUserID)
	if err != nil {
		if err.Error() == "you don't have access to the project" {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch subtasks",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(subtasks)
}

func (sh *subtaskHandler) CreateSubtask(ctx *fiber.Ctx) error {
	var subtask models.Subtask
	if err := ctx.BodyParser(&subtask); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	requestUserID := ctx.Locals("userID").(uint)

	subtaskID, err := sh.repository.CreateSubtask(&subtask, requestUserID)
	if err != nil {
		if err.Error() == "you don't have access to the project" {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create subtask",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"subtask_id": subtaskID,
	})
}

func (sh *subtaskHandler) GetSubtaskByID(ctx *fiber.Ctx) error {
	subtaskID, err := ctx.ParamsInt("subtask_id")
	if err != nil || subtaskID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subtask ID",
		})
	}

	requestUserID := ctx.Locals("userID").(uint)

	subtask, err := sh.repository.GetSubtaskByID(uint(subtaskID), requestUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Subtask not found",
			})
		}
		if err.Error() == "you don't have access to the project" {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch subtask",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(subtask)
}

func (sh *subtaskHandler) UpdateSubtask(ctx *fiber.Ctx) error {
	var subtask models.Subtask
	if err := ctx.BodyParser(&subtask); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	requestUserID := ctx.Locals("userID").(uint)

	subtaskID, err := sh.repository.UpdateSubtask(&subtask, requestUserID)
	if err != nil {
		if err.Error() == "you don't have access to the project" {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update subtask",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"subtask_id": subtaskID,
	})
}

func (sh *subtaskHandler) DeleteSubtask(ctx *fiber.Ctx) error {
	subtaskID, err := ctx.ParamsInt("subtask_id")
	if err != nil || subtaskID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subtask ID",
		})
	}

	requestUserID := ctx.Locals("userID").(uint)

	deletedID, err := sh.repository.DeleteSubtask(uint(subtaskID), requestUserID)
	if err != nil {
		if err.Error() == "you don't have access to the project" {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete subtask",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"deleted_subtask_id": deletedID,
	})
}
