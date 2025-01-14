package router

import (
	"github.com/gofiber/fiber/v2"
	"mizito/internal/database"
	"mizito/internal/handlers"
)

func RegisterRoutes(app *fiber.App, db *database.DatabaseHandler) {
	routes := app.Group("/teams")

	th := handlers.NewTeamHandler(db)

	routes.Get("/", th.GetTeams)
	routes.Get("/:id", th.GetTeamByID)
	routes.Get("/:id/projects", th.GetProjectsByTeam)
	routes.Post("/add-users", th.AddUsersToTeam)
	routes.Delete("/remove-users", th.DeleteUsersFromTeam)
	routes.Post("/create", th.CreateTeam)
	routes.Put("/update", th.UpdateTeam)
	routes.Delete("/delete-task", th.DeleteTask)
}
