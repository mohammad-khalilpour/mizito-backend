package router

import (
	"mizito/internal/database"
	"mizito/internal/handlers"
)

func InitSubtask(r *Router, postgreSql *database.DatabaseHandler) {
	sHandler := handlers.NewSubtaskHandler(postgreSql)

	SubtasksApp := r.App.Group("/subtasks")
	SubtasksApp.Get("/all", sHandler.GetSubtasksByTask)
	SubtasksApp.Get("/:subtask_id", sHandler.GetSubtaskByID)
	SubtasksApp.Put("/:subtask_id", sHandler.UpdateSubtask)
	SubtasksApp.Delete("/:subtask_id", sHandler.DeleteSubtask)

	subtaskApp := r.App.Group("/subtask")
	subtaskApp.Post("", sHandler.CreateSubtask)
}
