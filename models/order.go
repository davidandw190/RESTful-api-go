package models

import "time"

type Order struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	UserRefer uint64    `json:"user_id"`
	User      User      `gorm:"foreignKey:UserReference"`
	Products  []Product `gorm:"many2many:order_products;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type OrderProduct struct {
	OrderID   uint64
	ProductID uint64
}
