package router


import handlers "mizito/internal/handlers/http"

func InitTask(r *Router) {
	tHandler := handlers.NewTaskRepository()

	TaskApp := r.app.Group("/tasks")
	TaskApp.Get("/all", tHandler.GetTasksByProject)
	TaskApp.Get("/:task_id", tHandler.GetTaskByID)
	TaskApp.Put("/:task_id", tHandler.UpdateTask)
	TaskApp.Delete("/:task_id")


	subtaskApp := r.app.Group("task")
	subtaskApp.Post("/:task_id", tHandler.CreateTask)
}