package models

import (
	"time"
)

type Project struct {
	ID             uint `gorm:"primaryKey;autoIncrement"`
	TeamID         uint
	ProjectMembers []User    `gorm:"many2many:users_projects"`
	ProjectTasks   []Task    `gorm:"foreignKey:ProjectID"`
	Name           string    `validate:"required;min=5"`
	Messages       []Message `gorm:"foreignKey:ProjectID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ImageUrl       string
}
