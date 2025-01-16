package repositories

import (
	"errors"
	"gorm.io/gorm"
	"mizito/internal/database"
	"mizito/internal/repositories/utils"
	"mizito/pkg/models"
)

type SubtaskRepository interface {
	GetSubtasksByTask(taskID uint, requestUserID uint) ([]models.Subtask, error)
	CreateSubtask(subtask *models.Subtask, requestUserID uint) (uint, error)
	GetSubtaskByID(subtaskID uint, requestUserID uint) (*models.Subtask, error)
	UpdateSubtask(subtask *models.Subtask, requestUserID uint) (uint, error)
	DeleteSubtask(subtask *models.Subtask, requestUserID uint) (uint, error)
}

type subtaskRepository struct {
	permissionRepo utils.PermissionRepository
	DB             *gorm.DB
}

func NewSubtaskRepository(postgreSql *database.DatabaseHandler) SubtaskRepository {
	permissionRepo := utils.NewPermissionRepository(postgreSql)
	return &subtaskRepository{DB: postgreSql.DB, permissionRepo: permissionRepo}
}

// GetSubtasksByTask fetches all subtasks for a given task if the user is an admin of the task.
func (sr *subtaskRepository) GetSubtasksByTask(taskID uint, requestUserID uint) ([]models.Subtask, error) {
	// Check if the user has admin permission for the task
	if !sr.permissionRepo.CheckUserHasAccessToTask(taskID, requestUserID) {
		return nil, errors.New("you don't have access to the project")
	}

	var subtasks []models.Subtask
	if err := sr.DB.Where("task_id = ?", taskID).Find(&subtasks).Error; err != nil {
		return nil, err
	}
	return subtasks, nil
}

// CreateSubtask creates a new subtask under a task if the user is an admin of the task.
func (sr *subtaskRepository) CreateSubtask(subtask *models.Subtask, requestUserID uint) (uint, error) {
	// Check if the user has admin permission for the task
	if !sr.permissionRepo.CheckUserIsAdminOfTask(subtask.TaskID, requestUserID) {
		return 0, errors.New("you don't have access to the project")
	}

	if err := sr.DB.Create(subtask).Error; err != nil {
		return 0, err
	}
	return subtask.ID, nil
}

// GetSubtaskByID fetches a subtask by its ID if the user is an admin of the associated task.
func (sr *subtaskRepository) GetSubtaskByID(subtaskID uint, requestUserID uint) (*models.Subtask, error) {

	var subtask models.Subtask
	if err := sr.DB.First(&subtask, subtaskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if !sr.permissionRepo.CheckUserHasAccessToTask(subtask.TaskID, requestUserID) {
		return nil, errors.New("you don't have access to the project")
	}

	return &subtask, nil
}

// UpdateSubtask updates an existing subtask if the user is an admin of the associated task.
func (sr *subtaskRepository) UpdateSubtask(subtask *models.Subtask, requestUserID uint) (uint, error) {
	// Check if the user has admin permission for the task
	if !sr.permissionRepo.CheckUserHasAccessToTask(subtask.TaskID, requestUserID) {
		return 0, errors.New("you don't have access to the project")
	}

	if err := sr.DB.Save(subtask).Error; err != nil {
		return 0, err
	}
	return subtask.ID, nil
}

// DeleteSubtask deletes a subtask if the user is an admin of the associated task.
func (sr *subtaskRepository) DeleteSubtask(subtask *models.Subtask, requestUserID uint) (uint, error) {
	// Check if the user has admin permission for the task
	if !sr.permissionRepo.CheckUserIsAdminOfProject(subtask.TaskID, requestUserID) {
		return 0, errors.New("you don't have access to the project")
	}

	if err := sr.DB.Delete(subtask).Error; err != nil {
		return 0, err
	}
	return subtask.ID, nil
}
