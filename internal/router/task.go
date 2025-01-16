package router

import (
	"mizito/internal/database"
	"mizito/internal/handlers"
)

func InitTask(r *Router, postgreSql *database.DatabaseHandler) {
	tHandler := handlers.NewTaskHandler(postgreSql)

	TaskApp := r.App.Group("/tasks")
	TaskApp.Get("/all", tHandler.GetTasksByProject)
	TaskApp.Get("/:task_id", tHandler.GetTaskByID)
	TaskApp.Put("/:task_id", tHandler.UpdateTask)
	TaskApp.Delete("/:task_id", tHandler.DeleteTask)
	TaskApp.Put("/assign_task", tHandler.AssignTask)

	subtaskApp := r.App.Group("/task")
	subtaskApp.Post("/", tHandler.CreateTask)
}
