package models

import "time"

// Payment records money moving for an order. transaction_id is the
// external gateway identifier; UNIQUE + NOT NULL gives idempotency
// (a duplicated webhook with the same transaction_id will be rejected
// by the DB before the service sees it). paid_at nullable — not
// settled until the gateway confirms.
type Payment struct {
	PaymentID     uint64     `gorm:"primaryKey;autoIncrement"`
	OrderID       uint64     `gorm:"index:idx_payment_order;not null"`
	Method        string     `gorm:"not null;size:50"`
	Amount        int64      `gorm:"type:numeric(12,2);not null"`
	Status        string     `gorm:"not null;size:50;default:pendiente"`
	TransactionID string     `gorm:"uniqueIndex:uq_payment_transaction;not null;size:200"`
	PaidAt        *time.Time
	CreatedAt     time.Time  `gorm:"not null;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"not null;autoUpdateTime"`
}

func (Payment) TableName() string {
	return "Payment"
}
