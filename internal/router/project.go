package router

import handlers "mizito/internal/handlers/http"




func InitProject(r *Router) {
	pHandler := handlers.NewProjectRepository()

	projectsApp := r.app.Group("/projects")
	projectsApp.Get("/all", pHandler.GetProjectsByUser)
	projectsApp.Get("/:project_id", pHandler.GetProjectByID)
	projectsApp.Put("/:project_id", pHandler.UpdateProject)
	projectsApp.Delete("/:project_id")


	projectApp := r.app.Group("project")
	projectApp.Post("/project", pHandler.CreateProject)
}