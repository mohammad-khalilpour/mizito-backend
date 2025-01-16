package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"mizito/internal/database"
	"mizito/internal/repositories/utils"
	"mizito/pkg/models"
)

type ProjectCrudRepo interface {
	CreateProject(project *models.Project, requestUserID uint) (uint, error)
	UpdateProject(projectID uint, project *models.Project, requestUserID uint) (uint, error)
	DeleteProject(projectID uint, requestUserID uint) (uint, error)
	GetProjectByID(projectID uint, requestUserID uint) (*models.Project, error)
}

type ProjectDetailRepo interface {
	GetProjectsByUser(userID uint) ([]models.Project, error)
	// GetProjectMembers will use redis to cache the members for clients
	// one problem is that each message might require a query lookup for db to find the corresponding users
	// associated with ProjectID which is derived from event
	// the method can leverage redis handler and main db
	// one pattern might be cache aside pattern
	GetProjectMembers(ProjectID uint) ([]models.TeamMember, error)
	AddUserToProject(ProjectID uint, userID uint, requestUserID uint) error
	GetUsersByProjectID(ProjectID uint, requestUserID uint) ([]uint, error)
}

type ProjectRepository interface {
	ProjectDetailRepo
	ProjectCrudRepo
}

type projectRepository struct {
	permissionRepo utils.PermissionRepository
	DB             *gorm.DB
}

func NewProjectRepository(postgreSql *database.DatabaseHandler) ProjectRepository {
	permissionRepo := utils.NewPermissionRepository(postgreSql)
	return &projectRepository{DB: postgreSql.DB, permissionRepo: permissionRepo}
}

func (th *projectRepository) GetProjectsByUser(userID uint) ([]models.Project, error) {
	var projects []models.Project

	var teamMembers []models.TeamMember
	if err := th.DB.Where("user_id = ?", userID).Find(&teamMembers).Error; err != nil {
		return nil, err
	}

	var teams []models.Team
	var teamIDs []uint
	for _, tm := range teamMembers {
		teamIDs = append(teamIDs, tm.TeamID)
	}

	if err := th.DB.Where("id IN ?", teamIDs).Find(&teams).Error; err != nil {
		return nil, err
	}

	if err := th.DB.Where("team_id IN ?", teamIDs).Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (th *projectRepository) CreateProject(project *models.Project, requestUserID uint) (uint, error) {
	if !th.permissionRepo.CheckUserIsAdminOfTeam(requestUserID, project.TeamID) {
		return 0, errors.New("you are not an admin of the team")
	}

	tx := th.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the project within the transaction
	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Add the user to the project as a member within the transaction
	if err := tx.Model(project).Association("ProjectMembers").Append(&models.User{ID: requestUserID}); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to add user to project: %w", err)
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return project.ID, nil
}

func (th *projectRepository) GetProjectByID(projectID uint, requestUserID uint) (*models.Project, error) {
	if !th.permissionRepo.CheckUserHasAccessToProject(projectID, requestUserID) {
		return nil, errors.New("you don't have access to the project")
	}

	var project models.Project
	if err := th.DB.First(&project, projectID).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (th *projectRepository) UpdateProject(projectID uint, project *models.Project, requestUserID uint) (uint, error) {
	if !th.permissionRepo.CheckUserIsAdminOfProject(projectID, requestUserID) {
		return 0, errors.New("only admins can update the project")
	}

	var existingProject models.Project
	if err := th.DB.First(&existingProject, projectID).Error; err != nil {
		return 0, err
	}

	if err := th.DB.Model(&existingProject).Updates(project).Error; err != nil {
		return 0, err
	}
	return existingProject.ID, nil
}

func (th *projectRepository) DeleteProject(projectID uint, requestUserID uint) (uint, error) {
	if !th.permissionRepo.CheckUserIsAdminOfProject(projectID, requestUserID) {
		return 0, errors.New("only admins can update the project")
	}

	if err := th.DB.Delete(&models.Project{}, projectID).Error; err != nil {
		return 0, err
	}
	return projectID, nil
}

func (th *projectRepository) GetUsersByProjectID(projectID uint, requestUserID uint) ([]uint, error) {
	if !th.permissionRepo.CheckUserHasAccessToProject(projectID, requestUserID) {
		return nil, errors.New("you don't have access to the project")
	}

	var userIDs []uint
	if err := th.DB.Model(&models.TeamMember{}).Where("project_id = ?", projectID).Pluck("user_id", &userIDs).Error; err != nil {
		return nil, err
	}
	return userIDs, nil
}

func (th *projectRepository) GetProjectMembers(userID uint) ([]models.TeamMember, error) {
	var teamMembers []models.TeamMember
	if err := th.DB.Where("user_id = ?", userID).Find(&teamMembers).Error; err != nil {
		return nil, err
	}
	return teamMembers, nil
}

func (th *projectRepository) AddUserToProject(projectID uint, userID uint, requestUserID uint) error {
	if !th.permissionRepo.CheckUserIsAdminOfProject(projectID, requestUserID) {
		return errors.New("only admins can update the project")
	}
	var project models.Project
	if err := th.DB.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("project with ID %d does not exist", projectID)
		}
		return fmt.Errorf("failed to fetch project: %w", err)
	}

	var user models.User
	if err := th.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %d does not exist", userID)
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	if err := th.DB.Model(&project).Association("ProjectMembers").Append(&user); err != nil {
		return fmt.Errorf("failed to add user to project: %w", err)
	}

	return nil
}
