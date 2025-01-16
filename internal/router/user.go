package router

import (
	"mizito/internal/database"
	"mizito/internal/handlers"
)

func InitUser(r *Router, postgreSql *database.DatabaseHandler) {
	uHandler := handlers.NewUserHandler(postgreSql)

	TaskApp := r.App.Group("/users")
	TaskApp.Get("/all", uHandler.GetUsers)
	TaskApp.Get("/:user_id", uHandler.GetUserByID)
	TaskApp.Put("/:user_id", uHandler.UpdateUser)
	TaskApp.Delete("/:user_id", uHandler.DeleteUser)

	subtaskApp := r.App.Group("/user")
	subtaskApp.Post("/", uHandler.CreateUser)
}
