package repositories

import (
	"errors"
	"fmt"

	"mizito/internal/database"
	"mizito/pkg/models"

	"gorm.io/gorm"
)

type TeamRepository interface {
	GetTeams(userID uint) ([]models.Team, error)
	GetTeamByID(teamID uint) (*models.Team, error)
	GetProjectsByTeam(teamID uint) ([]*models.Project, error)
	AddUsersToTeam(usernames []string, teamID uint, role models.Role) (uint, error)
	DeleteUsersFromTeam(userIDs []uint, teamID uint) (uint, error)
	CreateTeam(team *models.Team, requestUserID uint) (uint, error)
	UpdateTeam(team *models.Team) (uint, error)
	DeleteTeam(teamID uint) (uint, error)
	DeleteTasks(teamID uint) (uint, error)
}

type teamRepository struct {
	db *database.DatabaseHandler
}

func NewTeamRepository(db *database.DatabaseHandler) TeamRepository {
	return &teamRepository{
		db: db,
	}
}

func (tr *teamRepository) GetTeams(userID uint) ([]models.Team, error) {
	var teams []models.Team
	err := tr.db.DB.
		Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ?", userID).
		Preload("Projects").
		Preload("Members").
		Find(&teams).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get teams for user %d: %w", userID, err)
	}
	return teams, nil
}

func (tr *teamRepository) GetTeamByID(teamID uint) (*models.Team, error) {
	var team models.Team
	err := tr.db.DB.
		Preload("Projects").
		Preload("Members").
		First(&team, teamID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get team %d: %w", teamID, err)
	}
	return &team, nil
}

func (tr *teamRepository) GetProjectsByTeam(teamID uint) ([]*models.Project, error) {
	var projects []*models.Project
	err := tr.db.DB.
		Where("team_id = ?", teamID).
		Find(&projects).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get projects for team %d: %w", teamID, err)
	}
	return projects, nil
}

func (tr *teamRepository) AddUsersToTeam(usernames []string, teamID uint, role models.Role) (uint, error) {
	if len(usernames) == 0 {
		return 0, errors.New("no user IDs provided")
	}

	tx := tr.db.DB.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var team models.Team
	if err := tx.First(&team, teamID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("team %d not found", teamID)
		}
		return 0, fmt.Errorf("failed to retrieve team %d: %w", teamID, err)
	}

	addedCount := uint(0)
	for _, username := range usernames {
		var user models.User
		if err := tx.Where("username = ?", username).First(&user).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, fmt.Errorf("user %s not found", username)
			}
			return 0, fmt.Errorf("failed to retrieve user %s: %w", username, err)
		}

		var existingMember models.TeamMember
		err := tx.Where("user_id = ? AND team_id = ?", user.ID, teamID).First(&existingMember).Error
		if err == nil {
			existingMember.Role = role
			if err := tx.Save(&existingMember).Error; err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("failed to update role for user %d in team %d: %w", username, teamID, err)
			}
			addedCount++
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return 0, fmt.Errorf("failed to check existing membership for user %d: %w", username, err)
		}

		tm := models.TeamMember{
			UserID: user.ID,
			TeamID: teamID,
			Role:   role,
		}
		if err := tx.Create(&tm).Error; err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to add user %d to team %d: %w", username, teamID, err)
		}
		addedCount++
	}

	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return addedCount, nil
}

func (tr *teamRepository) DeleteUsersFromTeam(userIDs []uint, teamID uint) (uint, error) {
	if len(userIDs) == 0 {
		return 0, errors.New("no user IDs provided")
	}

	tx := tr.db.DB.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var team models.Team
	if err := tx.First(&team, teamID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("team %d not found", teamID)
		}
		return 0, fmt.Errorf("failed to retrieve team %d: %w", teamID, err)
	}

	result := tx.Where("user_id IN ? AND team_id = ?", userIDs, teamID).Delete(&models.TeamMember{})
	if result.Error != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to delete users from team %d: %w", teamID, result.Error)
	}

	deletedCount := result.RowsAffected
	if deletedCount == 0 {
		tx.Rollback()
		return 0, fmt.Errorf("no users were deleted from team %d", teamID)
	}

	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return uint(deletedCount), nil
}

func (tr *teamRepository) CreateTeam(team *models.Team, requestUserID uint) (uint, error) {
	if team == nil {
		return 0, errors.New("team cannot be nil")
	}

	tx := tr.db.DB.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	if err := tx.Create(team).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create team: %w", err)
	}

	var teamMember = models.TeamMember{TeamID: team.ID, UserID: requestUserID, Role: models.Admin}

	if err := tx.Create(teamMember).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create team: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return team.ID, nil
}

func (tr *teamRepository) UpdateTeam(team *models.Team) (uint, error) {
	if team == nil {
		return 0, errors.New("team cannot be nil")
	}

	tx := tr.db.DB.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingTeam models.Team
	if err := tx.First(&existingTeam, team.ID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("team %d not found", team.ID)
		}
		return 0, fmt.Errorf("failed to retrieve team %d: %w", team.ID, err)
	}

	if err := tx.Model(&existingTeam).Updates(team).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to update team %d: %w", team.ID, err)
	}

	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return team.ID, nil
}

func (tr *teamRepository) DeleteTeam(teamID uint) (uint, error) {
	tx := tr.db.DB.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var team models.Team
	if err := tx.First(&team, teamID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("team %d not found", teamID)
		}
		return 0, fmt.Errorf("failed to retrieve team %d: %w", teamID, err)
	}

	if err := tx.Where("team_id = ?", teamID).Delete(&models.Task{}).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to delete tasks for team %d: %w", teamID, err)
	}

	if err := tx.Where("team_id = ?", teamID).Delete(&models.TeamMember{}).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to delete team members for team %d: %w", teamID, err)
	}

	if err := tx.Delete(&team).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to delete team %d: %w", teamID, err)
	}

	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return teamID, nil
}

func (tr *teamRepository) DeleteTasks(teamID uint) (uint, error) {
	result := tr.db.DB.Where("team_id = ?", teamID).Delete(&models.Task{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete tasks for team %d: %w", teamID, result.Error)
	}
	return uint(result.RowsAffected), nil
}
