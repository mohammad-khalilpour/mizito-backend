package models


type Team struct{
	ID uint `gorm:"primaryKey"`
	Members	[]TeamMember `gorm:"manyymany:users_teams;"`
}

type TeamMember struct {
	User User
	Role Role
}


type Role string 


const (
	Admin Role = "admin"
	Member Role = "member"
)