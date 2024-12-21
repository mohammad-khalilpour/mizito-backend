package repositories

import (
	"mizito/pkg/models"
)


type ProjectRepository interface {
	GetProjectsByUser(userID uint) ([]models.Project, error)
	CreateProject(project *models.Project) (uint, error)
	GetProjectByID(projectID uint)	(*models.Project, error)
	UpdateProject(projectID uint, project *models.Project) (uint, error)
	DeleteProject(projectID uint) (uint, error)
}


type projectRepository struct {

}


func NewProjectRepository() ProjectRepository{
	return &projectRepository{}
}



func (ph *projectRepository) GetProjectsByUser(userID uint) ([]models.Project, error) {
	return nil, nil
}
func (ph *projectRepository) CreateProject(project *models.Project) (uint, error) {
	return 0, nil
}
func (ph *projectRepository) GetProjectByID(projectID uint)	(*models.Project, error) {
	return nil, nil
}
func (ph *projectRepository) UpdateProject(projectID uint, project *models.Project) (uint, error) {
	return 0, nil
}
func (ph *projectRepository) DeleteProject(projectID uint) (uint, error) {
	return 0, nil
}



