package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `validate:"required" gorm:"unique"`
	Password  string
	Reports   []Report `gorm:"foreignKey:UserID"`
	Email     string   `validate:"required, endswith=@gmail.com" gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
