package models

import "time"

// OrderItem is one line on an order. Snapshot of price at purchase
// time (unit_price may differ from Items.price today). Subtotal is
// precomputed — service layer is responsible for unit_price * quantity
// integrity. CHECK constraint chk_order_item_qty_positive enforces
// quantity > 0 at DB level.
type OrderItem struct {
	OrderItemID uint64    `gorm:"primaryKey;autoIncrement"`
	OrderID     uint64    `gorm:"not null;uniqueIndex:uq_order_item,priority:1"`
	ItemID      uint64    `gorm:"not null;uniqueIndex:uq_order_item,priority:2;index:idx_order_item_item"`
	Quantity    int       `gorm:"not null"`
	UnitPrice   int64     `gorm:"type:numeric(12,2);not null"`
	Subtotal    int64     `gorm:"type:numeric(12,2);not null"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
}

func (OrderItem) TableName() string {
	return "Order_item"
}
