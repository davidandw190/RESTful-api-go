package models

import "time"

type Order struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ProductRefer uint      `json:"product_id"`
	Product      Product   `gorm:"foreignKey:ProductReferer"`
	UserReferer  uint      `json:"user_id"`
	User         User      `gorm:"foreignKey:UserReference"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
