package repositories

import (
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
}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

func (ph *projectRepository) GetProjectsByUser(userID uint) ([]models.Project, error) {
	return nil, nil
}
func (ph *projectRepository) CreateProject(project *models.Project) (uint, error) {
	return 0, nil
}
func (ph *projectRepository) GetProjectByID(projectID uint) (*models.Project, error) {
	return nil, nil
}
func (ph *projectRepository) UpdateProject(projectID uint, project *models.Project) (uint, error) {
	return 0, nil
}
func (ph *projectRepository) DeleteProject(projectID uint) (uint, error) {
	return 0, nil
}
func (ph *projectRepository) GetUsersByProjectID(ProjectID uint) ([]uint, error) {
	return nil, nil
}
func (ph *projectRepository) GetProjectMembers(userID uint) ([]models.TeamMember, error) {
	return nil, nil
}
