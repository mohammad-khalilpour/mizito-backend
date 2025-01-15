package utils

import (
	"mizito/internal/database"
	"mizito/pkg/models"
)

type TeamPermissionHandler interface {
	CheckUserHasAccessToTeam(userId uint, teamId uint) bool
	CheckUserIsAdminOfTeam(userId uint, teamId uint) bool
}

type TaskPermissionHandler interface {
	CheckUserHasAccessToTask(taskId uint, userId uint) bool
	CheckUserIsAdminOfTask(taskId uint, userId uint) bool
}

type ProjectPermissionHandler interface {
	CheckUserHasAccessToProject(projectId uint, userId uint) bool
	CheckUserIsAdminOfProject(projectId uint, userId uint) bool
}

type PermissionHandler interface {
	TeamPermissionHandler
	TaskPermissionHandler
	ProjectPermissionHandler
}

type permissionHandler struct {
	db *database.DatabaseHandler
}

func NewPermissionHandler() PermissionHandler {
	return &permissionHandler{}

}

func (ph *permissionHandler) CheckUserHasAccessToTeam(userId uint, teamId uint) bool {
	var count int64
	err := ph.db.DB.Model(&models.TeamMember{}).
		Where("user_id = ? AND team_id = ?", userId, teamId).
		Count(&count).Error
	if err != nil {
		// Handle error appropriately, possibly logging it
		return false
	}
	return count > 0
}

func (ph *permissionHandler) CheckUserIsAdminOfTeam(userId uint, teamId uint) bool {
	var count int64
	err := ph.db.DB.Model(&models.TeamMember{}).
		Where("user_id = ? AND team_id = ? AND role = ?", userId, teamId, models.Admin).
		Count(&count).Error
	if err != nil {
		// Handle error appropriately
		return false
	}
	return count > 0
}

func (ph *permissionHandler) CheckUserHasAccessToTask(taskId uint, userId uint) bool {
	var task models.Task
	err := ph.db.DB.First(&task, taskId).Error
	if err != nil {
		// Task not found or other error
		return false
	}


	// Otherwise, check if the user is a member of the task
	var count int64
	err = ph.db.DB.Model(&models.Task{}).
		Where("id = ?", taskId).
		Joins("JOIN task_members ON task_members.task_id = tasks.id").
		Where("task_members.user_id = ?", userId).
		Count(&count).Error
	if err != nil {
		// Handle error appropriately, possibly logging it
		return false
	}

	return count > 0
}

func (ph *permissionHandler) CheckUserIsAdminOfTask(taskId uint, userId uint) bool {
	var task models.Task
	err := ph.db.DB.First(&task, taskId).Error
	if err != nil {
		// Task not found or other error
		return false
	}

	// Check if user is an admin of the project
	var projectAdmin bool = ph.CheckUserIsAdminOfProject(task.ProjectID, userId)
	if projectAdmin {
		return true
	}

	// Optionally, check for task-specific admin roles if applicable
	// For example, if there's a TaskMember model with roles
	// Adjust according to your actual data model

	return false
}

func (ph *permissionHandler) CheckUserHasAccessToProject(projectId uint, userId uint) bool {
	var project models.Project
	err := ph.db.DB.First(&project, projectId).Error
	if err != nil {
		return false
	}
	var count int64
	err = ph.db.DB.Model(&models.Project{}).
		Where("id = ?", projectId).
		Joins("JOIN users_projects ON users_projects.project_id = projects.id").
		Where("users_projects.user_id = ?", userId).
		Count(&count).Error
	if err != nil {
		return false
	}

	return count > 0
}

func (ph *permissionHandler) CheckUserIsAdminOfProject(projectId uint, userId uint) bool {
	var project models.Project
	err := ph.db.DB.First(&project, projectId).Error
	if err != nil {
		return false
	}

	// Check if user is an admin of the team
	var teamMember models.TeamMember
	err = ph.db.DB.Where("team_id = ? AND user_id = ? AND role = ?", project.TeamID, userId, models.Admin).First(&teamMember).Error
	if err == nil {
		// User is a team admin
		return true
	}

	return false
}
