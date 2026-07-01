package models

import "time"

// Cart is one shopping cart per user (FK to Users, nullable because ON
// DELETE SET NULL — user deletion leaves the cart orphaned).
type Cart struct {
	CartID    uint64    `gorm:"primaryKey;autoIncrement"`
	UserID    *uint64   `gorm:"uniqueIndex:uq_cart_user"`
	Status    string    `gorm:"not null;size:50;default:activo"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

func (Cart) TableName() string {
	return "Cart"
}
