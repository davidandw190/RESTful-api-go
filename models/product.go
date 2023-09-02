package models

import "time"

type Product struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	SerialNumber string    `json:"serial_num" gorm:"unique"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
