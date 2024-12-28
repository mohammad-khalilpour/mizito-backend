package repositories


import "mizito/pkg/models"


type TeamRepository interface {
	GetTeams() ([]models.Team, error)
	GetTeamByID(TeamID uint) (*models.Team, error)
	GetProjectsByTeam(TeamID uint) (*models.Project, error)
	AddUsersToTeam(userID []uint, teamID uint) (uint, error)
	DeleteUsersFromTeam(userID []uint, teamID uint) (uint, error)
	CreateTeam(Team *models.Team) (uint, error)
	UpdateTeam(Team *models.Team) (uint, error)
	DeleteTask(TeamID uint) (uint, error)
}


type teamRepository struct {

}

func NewTeamRepository () TeamRepository{
	return &teamRepository{}
}

func (tr *teamRepository) GetTeams() ([]models.Team, error) {
	return nil, nil
}
func (tr *teamRepository) GetTeamByID(TeamID uint) (*models.Team, error) {
	return nil, nil
}
func (tr *teamRepository) GetProjectsByTeam(TeamID uint) (*models.Project, error) {
	return nil, nil
}
func (tr *teamRepository) CreateTeam(Team *models.Team) (uint, error) {
	return 0, nil
}
func (tr *teamRepository) UpdateTeam(Team *models.Team) (uint, error) {
	return 0, nil
}
func (tr *teamRepository) DeleteTask(TeamID uint) (uint, error) {
	return 0, nil
}
func (tr *teamRepository) AddUsersToTeam(userID []uint, teamID uint) (uint, error) {
	return 0, nil
}
func (tr *teamRepository) DeleteUsersFromTeam(userID []uint, teamID uint) (uint, error) {
	return 0, nil
}




