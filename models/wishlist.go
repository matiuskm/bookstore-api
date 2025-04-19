package models

type Wishlist struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null"`
	BookID uint   `gorm:"not null"`
	Book   Book   `gorm:"foreignKey:BookID"`
	User   User   `gorm:"foreignKey:UserID"`
}