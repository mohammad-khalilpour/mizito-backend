package router

import (
	"mizito/internal/handlers"
)

func InitUser(r *Router) {
	uHandler := handlers.NewUserRepository()

	TaskApp := r.App.Group("/users")
	TaskApp.Get("/all", uHandler.GetUsers)
	TaskApp.Get("/:user_id", uHandler.GetUserByID)
	TaskApp.Put("/:user_id", uHandler.UpdateUser)
	TaskApp.Delete("/:user_id", uHandler.DeleteUser)

	subtaskApp := r.App.Group("user")
	subtaskApp.Post("/:user_id", uHandler.CreateUser)
}
