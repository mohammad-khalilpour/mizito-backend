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
	AddUsersToTeam(userIDs []uint, teamID uint, role models.Role) (uint, error)
	DeleteUsersFromTeam(userIDs []uint, teamID uint) (uint, error)
	CreateTeam(team *models.Team) (uint, error)
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

// GetTeams retrieves all teams that a specific user is a member of, including their projects and members.
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

// GetTeamByID retrieves a team by its ID, including its projects and members.
// Returns nil if the team is not found.
func (tr *teamRepository) GetTeamByID(teamID uint) (*models.Team, error) {
	var team models.Team
	err := tr.db.DB.
		Preload("Projects").
		Preload("Members").
		First(&team, teamID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Team not found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get team %d: %w", teamID, err)
	}
	return &team, nil
}

// GetProjectsByTeam retrieves all projects associated with a specific team.
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

// AddUsersToTeam adds multiple users to a team with a specified role.
// Uses a transaction to ensure atomicity.
func (tr *teamRepository) AddUsersToTeam(userIDs []uint, teamID uint, role models.Role) (uint, error) {
	if len(userIDs) == 0 {
		return 0, errors.New("no user IDs provided")
	}

	// Begin a transaction
	tx := tr.db.DB.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the team exists
	var team models.Team
	if err := tx.First(&team, teamID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("team %d not found", teamID)
		}
		return 0, fmt.Errorf("failed to retrieve team %d: %w", teamID, err)
	}

	// Iterate over user IDs and add them to the team
	addedCount := uint(0)
	for _, uid := range userIDs {
		// Check if the user exists
		var user models.User
		if err := tx.First(&user, uid).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, fmt.Errorf("user %d not found", uid)
			}
			return 0, fmt.Errorf("failed to retrieve user %d: %w", uid, err)
		}

		// Check if the user is already a member of the team
		var existingMember models.TeamMember
		err := tx.Where("user_id = ? AND team_id = ?", uid, teamID).First(&existingMember).Error
		if err == nil {
			// User is already a member; optionally, update their role
			existingMember.Role = role
			if err := tx.Save(&existingMember).Error; err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("failed to update role for user %d in team %d: %w", uid, teamID, err)
			}
			addedCount++
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return 0, fmt.Errorf("failed to check existing membership for user %d: %w", uid, err)
		}

		// Add the user to the team
		tm := models.TeamMember{
			UserID: uid,
			TeamID: teamID,
			Role:   role,
		}
		if err := tx.Create(&tm).Error; err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to add user %d to team %d: %w", uid, teamID, err)
		}
		addedCount++
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return addedCount, nil
}

// DeleteUsersFromTeam removes multiple users from a team.
// Uses a transaction to ensure atomicity.
func (tr *teamRepository) DeleteUsersFromTeam(userIDs []uint, teamID uint) (uint, error) {
	if len(userIDs) == 0 {
		return 0, errors.New("no user IDs provided")
	}

	// Begin a transaction
	tx := tr.db.DB.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the team exists
	var team models.Team
	if err := tx.First(&team, teamID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("team %d not found", teamID)
		}
		return 0, fmt.Errorf("failed to retrieve team %d: %w", teamID, err)
	}

	// Delete team members
	result := tx.Where("user_id IN ? AND team_id = ?", userIDs, teamID).Delete(&models.TeamMember{})
	if result.Error != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to delete users from team %d: %w", teamID, result.Error)
	}

	// Check how many records were deleted
	deletedCount := result.RowsAffected
	if deletedCount == 0 {
		tx.Rollback()
		return 0, fmt.Errorf("no users were deleted from team %d", teamID)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return uint(deletedCount), nil
}

// CreateTeam creates a new team, optionally adding initial members.
// Uses a transaction to ensure atomicity.
func (tr *teamRepository) CreateTeam(team *models.Team) (uint, error) {
	if team == nil {
		return 0, errors.New("team cannot be nil")
	}

	// Begin a transaction
	tx := tr.db.DB.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the team
	if err := tx.Create(team).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create team: %w", err)
	}

	// Optionally, add initial members if provided
	if len(team.Members) > 0 {
		for _, member := range team.Members {
			// Ensure the TeamID is set correctly
			member.TeamID = team.ID

			// Check if the user exists
			var user models.User
			if err := tx.First(&user, member.UserID).Error; err != nil {
				tx.Rollback()
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return 0, fmt.Errorf("user %d not found", member.UserID)
				}
				return 0, fmt.Errorf("failed to retrieve user %d: %w", member.UserID, err)
			}

			// Check if the user is already a member of the team
			var existingMember models.TeamMember
			err := tx.Where("user_id = ? AND team_id = ?", member.UserID, team.ID).First(&existingMember).Error
			if err == nil {
				// User is already a member; optionally, update their role
				existingMember.Role = member.Role
				if err := tx.Save(&existingMember).Error; err != nil {
					tx.Rollback()
					return 0, fmt.Errorf("failed to update role for user %d in team %d: %w", member.UserID, team.ID, err)
				}
				continue
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return 0, fmt.Errorf("failed to check existing membership for user %d: %w", member.UserID, err)
			}

			// Add the user to the team
			if err := tx.Create(&member).Error; err != nil {
				tx.Rollback()
				return 0, fmt.Errorf("failed to add user %d to team %d: %w", member.UserID, team.ID, err)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return team.ID, nil
}

// UpdateTeam updates an existing team.
// It updates only the provided fields and optionally manages members.
func (tr *teamRepository) UpdateTeam(team *models.Team) (uint, error) {
	if team == nil {
		return 0, errors.New("team cannot be nil")
	}

	// Begin a transaction
	tx := tr.db.DB.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the team exists
	var existingTeam models.Team
	if err := tx.First(&existingTeam, team.ID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("team %d not found", team.ID)
		}
		return 0, fmt.Errorf("failed to retrieve team %d: %w", team.ID, err)
	}

	// Update the team fields
	// Only update non-zero fields to prevent overwriting with zero values
	if err := tx.Model(&existingTeam).Updates(team).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to update team %d: %w", team.ID, err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return team.ID, nil
}

// DeleteTeam deletes a team and all associated tasks and team members.
// Ensures referential integrity by leveraging GORM's `OnDelete:CASCADE` constraints.
func (tr *teamRepository) DeleteTeam(teamID uint) (uint, error) {
	// Begin a transaction
	tx := tr.db.DB.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if the team exists
	var team models.Team
	if err := tx.First(&team, teamID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("team %d not found", teamID)
		}
		return 0, fmt.Errorf("failed to retrieve team %d: %w", teamID, err)
	}

	// Delete associated tasks
	if err := tx.Where("team_id = ?", teamID).Delete(&models.Task{}).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to delete tasks for team %d: %w", teamID, err)
	}

	// Delete team members
	if err := tx.Where("team_id = ?", teamID).Delete(&models.TeamMember{}).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to delete team members for team %d: %w", teamID, err)
	}

	// Delete the team
	if err := tx.Delete(&team).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to delete team %d: %w", teamID, err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return teamID, nil
}

// DeleteTasks deletes all tasks associated with a specific team.
// This method is retained from your original interface.
func (tr *teamRepository) DeleteTasks(teamID uint) (uint, error) {
	result := tr.db.DB.Where("team_id = ?", teamID).Delete(&models.Task{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete tasks for team %d: %w", teamID, result.Error)
	}
	return uint(result.RowsAffected), nil
}
