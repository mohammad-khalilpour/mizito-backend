package repositories

import (
	"mizito/pkg/models"
)


type ProjectRepository interface {
	GetProjectsByUser(userID uint) ([]models.Project, error)
	CreateProject(project models.Project) (uint, error)
	GetProjectByID(projectID uint)	(models.Project, error)
	UpdateProject(projectID uint, project models.Project) (uint, error)
	DeleteProject(projectID uint) (uint, error)
}