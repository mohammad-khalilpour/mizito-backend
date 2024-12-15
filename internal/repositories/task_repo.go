package repositories

import "mizito/pkg/models"


type TaskRepository interface {
	GetTasksByProject(projectID uint) ([]models.Task, error)
	CreateTask(task *models.Task) (uint, error)
	GetTaskByID(taskID uint) (models.Task, error)
	UpdateTask(task *models.Task) (uint, error)
	DeleteTask(taskID uint) (uint, error)
}