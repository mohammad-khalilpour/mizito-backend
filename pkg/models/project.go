package models

import (
	"time"
)

type Project struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	TeamID        uint
	ProjectMember []User    `gorm:"foreignKey:ProjectID"`
	ProjectTasks  []Task    ``
	Name          string    `validate:"required;min=5"`
	Messages      []Message `gorm:"foreignKey:ProjectID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ImageUrl      string
}
