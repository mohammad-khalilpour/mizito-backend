package router

import (
	"mizito/internal/handlers"
)

func InitProject(r *Router) {
	pHandler := handlers.NewProjectRepository()

	projectsApp := r.App.Group("/projects")
	projectsApp.Get("/all", pHandler.GetProjectsByUser)
	projectsApp.Get("/:project_id", pHandler.GetProjectByID)
	projectsApp.Put("/:project_id", pHandler.UpdateProject)
	projectsApp.Delete("/:project_id", pHandler.DeleteProject)

	projectApp := r.App.Group("project")
	projectApp.Post("/project", pHandler.CreateProject)
}
