package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"mizito/internal/database"
	"mizito/internal/repositories/utils"
	"mizito/pkg/models"
)

type TaskRepository interface {
	GetTasksByProject(projectID uint, requestUserID uint) ([]models.Task, error)
	CreateTask(task *models.Task, requestUserID uint) (uint, error)
	GetTaskByID(taskID uint, requestUserID uint) (*models.Task, error)
	UpdateTask(task *models.Task, requestUserID uint) (uint, error)
	DeleteTask(taskID uint, requestUserID uint) (uint, error)
	AssignTask(UserID uint, TaskTitle uint, requestUserID uint) error
}

type taskRepository struct {
	permissionRepo utils.PermissionRepository
	DB             *gorm.DB
}

func NewTaskRepository(postgreSql *database.DatabaseHandler) TaskRepository {
	permissionRepo := utils.NewPermissionRepository(postgreSql)
	return &taskRepository{DB: postgreSql.DB, permissionRepo: permissionRepo}
}

func (tr *taskRepository) GetTasksByProject(projectID uint, requestUserID uint) ([]models.Task, error) {
	// Check if the user has permission to view tasks for the project
	if !tr.permissionRepo.CheckUserHasAccessToTask(projectID, requestUserID) {
		return nil, errors.New("you don't have access to the project")
	}

	// Fetch tasks from the database
	var tasks []models.Task
	err := tr.DB.Preload("Subtasks").
		Preload("Members").
		Preload("Reports").
		Where("project_id = ?", projectID).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (tr *taskRepository) CreateTask(task *models.Task, requestUserID uint) (uint, error) {
	// Check if the user is an admin of the project associated with the task
	isAdmin := tr.permissionRepo.CheckUserIsAdminOfProject(task.ProjectID, requestUserID)
	if !isAdmin {
		return 0, errors.New("user does not have permission to create tasks for this project")
	}

	// Start a transaction to ensure atomicity
	tx := tr.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Add the task to the database
	if err := tx.Create(task).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Ensure the task is added to the project's tasks
	project := &models.Project{}
	if err := tx.Where("id = ?", task.ProjectID).First(project).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	project.ProjectTasks = append(project.ProjectTasks, *task)

	// Save the updated project
	if err := tx.Save(project).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return task.ID, nil
}

func (tr *taskRepository) GetTaskByID(taskID uint, requestUserID uint) (*models.Task, error) {
	// Check if the user has permission to view tasks for the project
	if !tr.permissionRepo.CheckUserHasAccessToTask(taskID, requestUserID) {
		return nil, errors.New("you don't have access to the project")
	}

	var task models.Task
	err := tr.DB.First(&task, taskID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return &task, nil
}

func (tr *taskRepository) UpdateTask(task *models.Task, requestUserID uint) (uint, error) {
	// Check if the user has permission to update the task
	isAdmin := tr.permissionRepo.CheckUserIsAdminOfTask(task.ID, requestUserID)
	if !isAdmin {
		return 0, errors.New("user does not have permission to change this task")
	}

	// Verify that the task exists in the database
	var existingTask models.Task
	err := tr.DB.First(&existingTask, task.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("task not found")
		}
		return 0, err
	}

	// Update the task in the database
	err = tr.DB.Model(&existingTask).Updates(task).Error
	if err != nil {
		return 0, err
	}

	return task.ID, nil
}

func (tr *taskRepository) DeleteTask(taskID uint, requestUserID uint) (uint, error) {
	// Check if the user has permission to delete the task
	isAdmin := tr.permissionRepo.CheckUserIsAdminOfTask(taskID, requestUserID)
	if !isAdmin {
		return 0, errors.New("user does not have permission to delete this task")
	}

	// Fetch the task
	var task models.Task
	err := tr.DB.First(&task, taskID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("task not found")
		}
		return 0, err
	}
	// Delete the task from the database
	err = tr.DB.Delete(&task).Error
	return task.ID, err
}

func (tr *taskRepository) AssignTask(userID uint, taskID uint, requestUserID uint) error {
	if !tr.permissionRepo.CheckUserIsAdminOfTask(taskID, requestUserID) {
		return errors.New("only admins can update the project")
	}

	var task models.Task
	if err := tr.DB.Preload("Project").First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("task with ID %d does not exist", taskID)
		}
		return fmt.Errorf("failed to fetch task: %w", err)
	}

	var project models.Project
	if err := tr.DB.Preload("ProjectMembers").First(&project, task.ProjectID).Error; err != nil {
		return fmt.Errorf("failed to fetch project: %w", err)
	}

	isMember := false
	var user models.User
	for _, member := range project.ProjectMembers {
		if member.ID == userID {
			user = member
			isMember = true
			break
		}
	}

	if !isMember {
		return fmt.Errorf("user with ID %d is not a member of the project associated with task %d", userID, taskID)
	}

	task.Members = append(task.Members, user)
	if err := tr.DB.Save(&task).Error; err != nil {
		return fmt.Errorf("failed to assign task to user: %w", err)
	}

	return nil
}
