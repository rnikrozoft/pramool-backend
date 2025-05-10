package entity

import "time"

type User struct {
	UserId    string `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
