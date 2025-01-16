package models

import "time"

type Task struct {
	ID                 uint `gorm:"primaryKey" validate:"required"`
	ProjectID          uint
	Title              string
	Description        string
	Subtasks           []Subtask `gorm:"foreignKey:TaskID"`
	TaskPriority       int       `validator:"lte=100,gt=0"`
	Members            []User    `gorm:"many2many:task_members;"`
	DueDate            time.Time
	Reports            []Report `gorm:"foreignKey:TaskID"`
	ProgressPercentage int      `validator:"gte=0;lte=100"`
}

type Report struct {
	ID        uint `gorm:"primaryKey"`
	TaskID    uint
	UserID    uint `gorm:"not null"`
	Message   string
	CreatedAt time.Time
}

type Subtask struct {
	ID          uint `gorm:"primaryKey"`
	TaskID      uint
	Title       string `validate:"required"`
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
