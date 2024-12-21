package repositories

import "mizito/pkg/models"


type TaskRepository interface {
	GetTasksByProject(projectID uint) ([]models.Task, error)
	CreateTask(task *models.Task) (uint, error)
	GetTaskByID(taskID uint) (*models.Task, error)
	UpdateTask(task *models.Task) (uint, error)
	DeleteTask(taskID uint) (uint, error)
}

type taskRepository struct {

}

func NewTaskRepository() TaskRepository{
	return &taskRepository{}
}


func (tr *taskRepository) GetTasksByProject(projectID uint) ([]models.Task, error) {
	return nil, nil
}
func (tr *taskRepository) CreateTask(task *models.Task) (uint, error) {
	return 0, nil
}
func (tr *taskRepository) GetTaskByID(taskID uint) (*models.Task, error) {
	return nil, nil
}
func (tr *taskRepository) UpdateTask(task *models.Task) (uint, error) {
	return 0, nil
}
func (tr *taskRepository) DeleteTask(taskID uint) (uint, error) {
	return 0, nil
}