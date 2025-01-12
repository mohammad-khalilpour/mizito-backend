package repositories

import (
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
	GetUsersByProjectID(ProjectID uint) ([]uint, error)
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
