package models

import "time"

// Order is a placed order. UserID is nullable so user account deletion
// preserves order history anonymously (FK constraint is SET NULL).
// ShippingAddressID is required and restricted — every order needs
// a real address and you can't drop an address that points to orders.
type Order struct {
	OrderID           uint64    `gorm:"primaryKey;autoIncrement"`
	UserID            *uint64   `gorm:"index:idx_orders_user,index:idx_orders_user_date,priority:1"`
	Total             int64     `gorm:"type:numeric(12,2);not null"`
	Status            string    `gorm:"index:idx_orders_status,index:idx_orders_status_date,priority:1;not null;size:50;default:pendiente"`
	ShippingAddressID uint64    `gorm:"index:idx_orders_shipping;not null"`
	CreatedAt         time.Time `gorm:"not null;autoCreateTime;index:idx_orders_user_date,priority:2;index:idx_orders_status_date,priority:2"`
	UpdatedAt         time.Time `gorm:"not null;autoUpdateTime"`
}

func (Order) TableName() string {
	return "Orders"
}
