package router

import (
	"mizito/internal/database"
	"mizito/internal/handlers"
)

func InitDashboard(r *Router, postgreSql *database.DatabaseHandler) {

	pHandler := handlers.NewDashboardHandler(postgreSql)

	projectsApp := r.App.Group("/projects")
	projectsApp.Get("", pHandler.GetDashboardDetails)

}
