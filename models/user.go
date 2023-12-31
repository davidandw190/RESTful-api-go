package models

import "time"

type User struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" gorm:"unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
