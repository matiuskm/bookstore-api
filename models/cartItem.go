package models

import "time"

type CartItem struct {
	ID			uint     `gorm:"primaryKey"`
	UserID		uint     
	User		User
	BookID		uint
	Book		Book
	Quantity	int
	CreatedAt	time.Time
}