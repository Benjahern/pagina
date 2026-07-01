package models

import "time"

// StockMovement is the audit log for any item's stock change. Positive
// change = entry (purchase from supplier, return); negative = exit
// (sale, damage). OrderID/UserID nullable so adjustments outside sales
// can be recorded without orphan refs to past orders/users.
type StockMovement struct {
	MovementID uint64    `gorm:"primaryKey;autoIncrement"`
	ItemID     uint64    `gorm:"index:idx_stock_movement_item;not null"`
	Change     int       `gorm:"not null"`
	Reason     string    `gorm:"not null;size:50"`
	OrderID    *uint64   `gorm:"index:idx_stock_movement_order"`
	UserID     *uint64
	CreatedAt  time.Time `gorm:"not null;autoCreateTime;index:idx_stock_movement_date"`
}

func (StockMovement) TableName() string {
	return "Stock_movement"
}
