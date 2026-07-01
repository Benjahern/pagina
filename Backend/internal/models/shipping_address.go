package models

import "time"

// ShippingAddress is a delivery address saved per user. Multiple per
// user; one can be marked is_default. The "only one default per user"
// rule is enforced via a PARTIAL UNIQUE INDEX in raw SQL after the
// DBML export — see database.sql Note in this table.
type ShippingAddress struct {
	ShippingAddressID uint64    `gorm:"primaryKey;autoIncrement"`
	UserID            *uint64   `gorm:"index:idx_shipping_user"`
	AddressLine       string    `gorm:"not null;size:500"`
	City              string    `gorm:"not null;size:200"`
	PostalCode        string    `gorm:"not null;size:50"`
	Commune           string    `gorm:"not null;size:200"`
	IsDefault         bool      `gorm:"not null;default:false"`
	CreatedAt         time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt         time.Time `gorm:"not null;autoUpdateTime"`
}

func (ShippingAddress) TableName() string {
	return "shipping_address"
}
