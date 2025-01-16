package repositories

import (
	"errors"
	"gorm.io/gorm"
	"mizito/internal/database"
	"mizito/pkg/models"
)

type DashboardRepository interface {
	GetDashboardDetails(requestUserID uint) (string, string, int, []TeamMemberWithRole, []SubtaskDetail, []ProjectDetail, error)
}

type dashboardRepository struct {
	DB *gorm.DB
}

func NewDashboardRepository(postgreSql *database.DatabaseHandler) DashboardRepository {
	return &dashboardRepository{DB: postgreSql.DB}
}

func (dr *dashboardRepository) GetDashboardDetails(requestUserID uint) (string, string, int, []TeamMemberWithRole, []SubtaskDetail, []ProjectDetail, error) {
	var user models.User
	var coworkers []TeamMemberWithRole
	var todoList []SubtaskDetail
	var projectList []ProjectDetail

	// Logic for fetching username and profile picture
	if err := dr.DB.First(&user, requestUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", 0, nil, nil, nil, errors.New("user not found")
		}
		return "", "", 0, nil, nil, nil, err
	}

	// Logic for fetching coworkers
	var teamMember models.TeamMember
	if err := dr.DB.Where("user_id = ? AND team_id IN (?)", requestUserID, dr.DB.Model(&models.TeamMember{}).Select("team_id")).First(&teamMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", 0, nil, nil, nil, errors.New("team not found for user")
		}
		return "", "", 0, nil, nil, nil, err
	}

	// Fetching all coworkers in the same team excluding the request user
	if err := dr.DB.Model(&models.TeamMember{}).
		Where("team_id = ? AND user_id != ?", teamMember.TeamID, requestUserID).
		Joins("JOIN users ON users.id = team_members.user_id").
		Select("users.username, team_members.role").
		Find(&coworkers).Error; err != nil {
		return "", "", 0, nil, nil, nil, err
	}

	// Logic for fetching the todoList (unfinished subtasks)
	var tasks []models.Task
	// Fetch all tasks the user is assigned to
	if err := dr.DB.Model(&models.Task{}).
		Joins("JOIN task_members ON task_members.task_id = tasks.id").
		Where("task_members.user_id = ?", requestUserID).
		Find(&tasks).Error; err != nil {
		return "", "", 0, nil, nil, nil, err
	}

	// Collect subtasks from these tasks that are not finished
	var unfinishedSubtasks []models.Subtask
	for _, task := range tasks {
		var taskSubtasks []models.Subtask
		if err := dr.DB.Where("task_id = ? AND is_completed = ?", task.ID, false).Find(&taskSubtasks).Error; err != nil {
			return "", "", 0, nil, nil, nil, err
		}
		unfinishedSubtasks = append(unfinishedSubtasks, taskSubtasks...)
	}

	// Logic for fetching projectList (project ID, name, and remaining tasks percentage)
	var projects []models.Project
	// Fetch all projects the user is involved in
	if err := dr.DB.Model(&models.Project{}).
		Joins("JOIN tasks ON tasks.project_id = projects.id").
		Joins("JOIN task_members ON task_members.task_id = tasks.id").
		Where("task_members.user_id = ?", requestUserID).
		Find(&projects).Error; err != nil {
		return "", "", 0, nil, nil, nil, err
	}

	// For each project, calculate the average remaining percentage of all tasks
	var projectDetails []ProjectDetail
	for _, project := range projects {
		var tasks []models.Task
		if err := dr.DB.Where("project_id = ?", project.ID).Find(&tasks).Error; err != nil {
			return "", "", 0, nil, nil, nil, err
		}

		var totalRemainingPct int
		var taskCount int
		for _, task := range tasks {
			totalRemainingPct += 100 - task.ProgressPercentage
			taskCount++
		}

		// Calculate the average remaining percentage
		remainingPct := 0
		if taskCount > 0 {
			remainingPct = totalRemainingPct / taskCount
		}

		projectDetails = append(projectDetails, ProjectDetail{
			ProjectName:   project.Name,
			CompletionPct: remainingPct,
		})
	}

	return user.Username, "", len(todoList), coworkers, todoList, projectList, nil
}

type TeamMemberWithRole struct {
	Name string
	Role string
}

type SubtaskDetail struct {
	SubtaskName string
	Detail      string
}

type ProjectDetail struct {
	ProjectName   string
	CompletionPct int
}
