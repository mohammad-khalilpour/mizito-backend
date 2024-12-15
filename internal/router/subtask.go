package router



import "mizito/internal/handlers"

func InitSubtask(r *Router) {
	sHandler := handlers.NewSubtaskRepository()

	SubtasksApp := r.app.Group("/subtasks")
	SubtasksApp.Get("/all", sHandler.GetSubtasksByTask)
	SubtasksApp.Get("/:subtask_id", sHandler.GetSubtaskByID)
	SubtasksApp.Put("/:subtask_id", sHandler.UpdateSubtask)
	SubtasksApp.Delete("/:subtask_id")


	subtaskApp := r.app.Group("project")
	subtaskApp.Post("/subtask:subtask_id", sHandler.CreateSubtask)
}