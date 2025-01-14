package models

type Role string

const (
	Admin  Role = "admin"
	Member Role = "member"
)

type Team struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Projects []Project    `gorm:"foreignKey:TeamID"`
	Members  []TeamMember `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE;"`
}

type TeamMember struct {
	UserID uint `gorm:"primaryKey"`
	TeamID uint `gorm:"primaryKey"`
	Role   Role `gorm:"not null"`
}
