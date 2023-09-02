package models

import "time"

type Order struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	ProductRefer uint64    `json:"product_id"`
	Product      Product   `gorm:"foreignKey:ProductReferer"`
	UserReferer  uint64    `json:"user_id"`
	User         User      `gorm:"foreignKey:UserReference"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
