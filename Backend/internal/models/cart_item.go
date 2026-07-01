package models

import "time"

// CartItem is one line in a cart. Composite unique (cart_id, item_id)
// enforces "same item only once per cart" — UPSERT bumps quantity.
type CartItem struct {
	CartItemID uint64    `gorm:"primaryKey;autoIncrement"`
	CartID     uint64    `gorm:"not null;uniqueIndex:uq_cart_item,priority:1"`
	ItemID     uint64    `gorm:"not null;uniqueIndex:uq_cart_item,priority:2"`
	Quantity   int       `gorm:"not null;default:1"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"not null;autoUpdateTime"`
}

func (CartItem) TableName() string {
	return "Cart_Item"
}
