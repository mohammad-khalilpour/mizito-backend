package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"mizito/internal/database"
	"mizito/internal/repositories"
	"mizito/pkg/models"
	"strconv"
)

type TaskHandler interface {
	GetTasksByProject(ctx *fiber.Ctx) error
	CreateTask(ctx *fiber.Ctx) error
	GetTaskByID(ctx *fiber.Ctx) error
	UpdateTask(ctx *fiber.Ctx) error
	DeleteTask(ctx *fiber.Ctx) error
}

type taskHandler struct {
	repository repositories.TaskRepository
}

func NewTaskHandler(db *database.DatabaseHandler) TaskHandler {
	repo := repositories.NewTaskRepository(db)
	return &taskHandler{repository: repo}
}

// GetTasksByProject fetches all tasks for a given project ID
func (th *taskHandler) GetTasksByProject(ctx *fiber.Ctx) error {
	projectID, err := strconv.ParseUint(ctx.Params("project_id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	requestUserID := ctx.Locals("userID").(uint)

	tasks, err := th.repository.GetTasksByProject(uint(projectID), requestUserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(tasks)
}

// CreateTask creates a new task
func (th *taskHandler) CreateTask(ctx *fiber.Ctx) error {
	var task models.Task
	if err := ctx.BodyParser(&task); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	requestUserID := ctx.Locals("userID").(uint)

	taskID, err := th.repository.CreateTask(&task, requestUserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"task_id": taskID})
}

// GetTaskByID fetches a task by its ID
func (th *taskHandler) GetTaskByID(ctx *fiber.Ctx) error {
	taskID, err := strconv.ParseUint(ctx.Params("task_id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	requestUserID := ctx.Locals("userID").(uint)

	task, err := th.repository.GetTaskByID(uint(taskID), requestUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(task)
}

// UpdateTask updates an existing task
func (th *taskHandler) UpdateTask(ctx *fiber.Ctx) error {
	var task models.Task
	if err := ctx.BodyParser(&task); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	requestUserID := ctx.Locals("userID").(uint)

	taskID, err := th.repository.UpdateTask(&task, requestUserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"updated_task_id": taskID})
}

// DeleteTask deletes a task by its ID
func (th *taskHandler) DeleteTask(ctx *fiber.Ctx) error {
	taskID, err := strconv.ParseUint(ctx.Params("task_id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	requestUserID := ctx.Locals("userID").(uint)

	deletedTaskID, err := th.repository.DeleteTask(uint(taskID), requestUserID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"deleted_task_id": deletedTaskID})
}
