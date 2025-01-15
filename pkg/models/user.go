package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `validate:"required"`
	Password  string
	Email     string `validate:"required, endswith=@gmail.com"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
