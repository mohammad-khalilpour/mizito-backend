package models

import "time"



type Project struct{
	ID string `gorm:""`
	Team Team `validate:"dive"`
	Name string `validate:"required;min=5"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ImageUrl string
}