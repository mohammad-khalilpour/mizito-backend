package models

import "time"


type Task struct{
	ID uint `gorm:"primaryKey" validate:"required"`
	Title string 
	Describtion string
	Subtasks	[]Subtask
	TaskPriority	int `validator:"lte=100,gt=0"`
	Members 		[]User
	DueDate			time.Time
	Reports			[]Report
	ProgressPercentage int `validator:"gte=0;lte=100"`
}


type Report struct {
	ID uint	`gorm:"primaryKey"`
	Member	User
	Message string
	CreatedAt time.Time
}



type Subtask struct {
	ID uint `gorm:"primaryKey"`
	Task Task	`validate:"dive"`
	Title string	`validate:"required"`
	IsCompleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}