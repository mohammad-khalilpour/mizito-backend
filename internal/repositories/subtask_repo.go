package repositories

import "mizito/pkg/models"

type SubtaskRepository interface {
	GetSubtasksByTask(taskID uint) ([]models.Subtask, error)
	CreateSubtask(subtask *models.Subtask) (uint, error)
	GetSubtaskByID(subtaskID uint) (models.Subtask, error)
	UpdateSubtask(subtask *models.Subtask) (uint, error)
	DeleteSubtask(subtask *models.Subtask) (uint, error)
}

