package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"mizito/internal/database"
	"mizito/pkg/models"
)

type ProjectCrudRepo interface {
	CreateProject(project *models.Project) (uint, error)
	UpdateProject(projectID uint, project *models.Project) (uint, error)
	DeleteProject(projectID uint) (uint, error)
	GetProjectByID(projectID uint) (*models.Project, error)
}

type ProjectDetailRepo interface {
	GetProjectsByUser(userID uint) ([]models.Project, error)
	// GetProjectMembers will use redis to cache the members for clients
	// one problem is that each message might require a query lookup for db to find the corresponding users
	// associated with ProjectID which is derived from event
	// the method can leverage redis handler and main db
	// one pattern might be cache aside pattern
	GetProjectMembers(userID uint) ([]models.TeamMember, error)
	AddUserToProject(ProjectID uint, userID uint) error
	GetUsersByProjectID(ProjectID uint) ([]uint, error)
	AssignTask(UserID uint, TaskTitle uint) error
}

type ProjectRepository interface {
	ProjectDetailRepo
	ProjectCrudRepo
}

type projectRepository struct {
	DB *gorm.DB
}

func NewProjectRepository(postgreSql *database.DatabaseHandler) ProjectRepository {
	return &projectRepository{DB: postgreSql.DB}
}

func (ph *projectRepository) GetProjectsByUser(userID uint) ([]models.Project, error) {
	var projects []models.Project

	var teamMembers []models.TeamMember
	if err := ph.DB.Where("user_id = ?", userID).Find(&teamMembers).Error; err != nil {
		return nil, err
	}

	var teams []models.Team
	var teamIDs []uint
	for _, tm := range teamMembers {
		teamIDs = append(teamIDs, tm.TeamID)
	}

	if err := ph.DB.Where("id IN ?", teamIDs).Find(&teams).Error; err != nil {
		return nil, err
	}

	if err := ph.DB.Where("team_id IN ?", teamIDs).Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (ph *projectRepository) CreateProject(project *models.Project) (uint, error) {
	if err := ph.DB.Create(project).Error; err != nil {
		return 0, err
	}
	return project.ID, nil
}

func (ph *projectRepository) GetProjectByID(projectID uint) (*models.Project, error) {
	var project models.Project
	if err := ph.DB.First(&project, projectID).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (ph *projectRepository) UpdateProject(projectID uint, project *models.Project) (uint, error) {
	var existingProject models.Project
	if err := ph.DB.First(&existingProject, projectID).Error; err != nil {
		return 0, err
	}

	if err := ph.DB.Model(&existingProject).Updates(project).Error; err != nil {
		return 0, err
	}
	return existingProject.ID, nil
}

func (ph *projectRepository) DeleteProject(projectID uint) (uint, error) {
	if err := ph.DB.Delete(&models.Project{}, projectID).Error; err != nil {
		return 0, err
	}
	return projectID, nil
}

func (ph *projectRepository) GetUsersByProjectID(projectID uint) ([]uint, error) {
	var userIDs []uint
	if err := ph.DB.Model(&models.TeamMember{}).Where("project_id = ?", projectID).Pluck("user_id", &userIDs).Error; err != nil {
		return nil, err
	}
	return userIDs, nil
}

func (ph *projectRepository) GetProjectMembers(userID uint) ([]models.TeamMember, error) {
	var teamMembers []models.TeamMember
	if err := ph.DB.Where("user_id = ?", userID).Find(&teamMembers).Error; err != nil {
		return nil, err
	}
	return teamMembers, nil
}

func (ph *projectRepository) AddUserToProject(projectID uint, userID uint) error {
	var project models.Project
	if err := ph.DB.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("project with ID %d does not exist", projectID)
		}
		return fmt.Errorf("failed to fetch project: %w", err)
	}

	var user models.User
	if err := ph.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %d does not exist", userID)
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	if err := ph.DB.Model(&project).Association("ProjectMembers").Append(&user); err != nil {
		return fmt.Errorf("failed to add user to project: %w", err)
	}

	return nil
}

func (ph *projectRepository) AssignTask(userID uint, taskID uint) error {
	var task models.Task
	if err := ph.DB.Preload("Project").First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("task with ID %d does not exist", taskID)
		}
		return fmt.Errorf("failed to fetch task: %w", err)
	}

	var project models.Project
	if err := ph.DB.Preload("ProjectMembers").First(&project, task.ProjectID).Error; err != nil {
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
	if err := ph.DB.Save(&task).Error; err != nil {
		return fmt.Errorf("failed to assign task to user: %w", err)
	}

	return nil
}
