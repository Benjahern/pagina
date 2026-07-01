package models

import "time"

// Review is a user's rating + comment on an item. Composite unique
// (user_id, item_id) makes it impossible for the same user to review
// the same item twice. CHECK constraint chk_review_rating_range
// enforces rating 1..5 at DB level.
type Review struct {
	ReviewID  uint64    `gorm:"primaryKey;autoIncrement"`
	UserID    *uint64   `gorm:"uniqueIndex:uq_review_user_item,priority:1"`
	ItemID    uint64    `gorm:"not null;index:idx_review_item;index:idx_review_item_approved,priority:1;uniqueIndex:uq_review_user_item,priority:2"`
	Rating    int       `gorm:"not null"`
	Comment   *string
	Approved  bool      `gorm:"not null;default:false;index:idx_review_item_approved,priority:2"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (Review) TableName() string {
	return "Review"
}
