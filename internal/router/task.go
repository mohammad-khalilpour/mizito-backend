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

	subtaskApp := r.App.Group("task")
	subtaskApp.Post("/:task_id", tHandler.CreateTask)
}
