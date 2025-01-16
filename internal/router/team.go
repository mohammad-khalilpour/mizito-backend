package router

import (
	"mizito/internal/database"
	"mizito/internal/handlers"
)

func InitTeam(r *Router, db *database.DatabaseHandler) {
	routes := r.App.Group("/teams")

	th := handlers.NewTeamHandler(db)

	routes.Get("/", th.GetTeams)
	routes.Get("/:id", th.GetTeamByID)
	routes.Get("/:id/projects", th.GetProjectsByTeam)
	routes.Post("/add-users", th.AddUsersToTeam)
	routes.Delete("/remove-users", th.DeleteUsersFromTeam)
	routes.Post("/create", th.CreateTeam)
	routes.Put("/update", th.UpdateTeam)
	//routes.Delete("/delete-task", th.DeleteTask)
}
