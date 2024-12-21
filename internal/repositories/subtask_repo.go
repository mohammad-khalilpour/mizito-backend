package repositories

import "mizito/pkg/models"

type SubtaskRepository interface {
	GetSubtasksByTask(taskID uint) ([]models.Subtask, error)
	CreateSubtask(subtask *models.Subtask) (uint, error)
	GetSubtaskByID(subtaskID uint) (*models.Subtask, error)
	UpdateSubtask(subtask *models.Subtask) (uint, error)
	DeleteSubtask(subtask *models.Subtask) (uint, error)
}


type subtaskRepository struct {

}

func NewSubtaskRepository() SubtaskRepository{
	return &subtaskRepository{}
}


func (sr *subtaskRepository) GetSubtasksByTask(taskID uint) ([]models.Subtask, error) {
	return nil, nil
}
func (sr *subtaskRepository) CreateSubtask(subtask *models.Subtask) (uint, error) {
	return 0, nil
}
func (sr *subtaskRepository) GetSubtaskByID(subtaskID uint) (*models.Subtask, error) {
	return nil, nil
}
func (sr *subtaskRepository) UpdateSubtask(subtask *models.Subtask) (uint, error) {
	return 0, nil
}
func (sr *subtaskRepository) DeleteSubtask(subtask *models.Subtask) (uint, error) {
	return 0, nil
}

