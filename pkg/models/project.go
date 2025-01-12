package models

import (
	"time"
)

type Project struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	TeamID        uint
	ProjectMember []User    `gorm:"many2many:users_projects"`
	ProjectTasks  []Task    `gorm:"many2many:projects_tasks"`
	Name          string    `validate:"required;min=5"`
	Messages      []Message `gorm:"foreignKey:ProjectID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ImageUrl      string
}
