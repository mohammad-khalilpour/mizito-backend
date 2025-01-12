package router

import (
	"mizito/internal/database"
	"mizito/internal/handlers"
)

func InitProject(r *Router, postgreSql *database.DatabaseHandler) {

	pHandler := handlers.NewProjectHandler(postgreSql)

	projectsApp := r.App.Group("/projects")
	projectsApp.Get("/all", pHandler.GetProjectsByUser)
	projectsApp.Get("/:project_id", pHandler.GetProjectByID)
	projectsApp.Put("/:project_id", pHandler.UpdateProject)
	projectsApp.Delete("/:project_id", pHandler.DeleteProject)
	projectsApp.Put("/:project_id/users", pHandler.GetUsersByProjectID)
	projectsApp.Put("/assign_task", pHandler.AssignTask)

	projectApp := r.App.Group("/project")
	projectApp.Post("", pHandler.CreateProject)

}
