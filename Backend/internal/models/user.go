package models

import "time"

// User represents a customer of the storefront. Column names come from
// GORM's NamingStrategy (CamelCase → snake_case): UserID → user_id,
// CreatedAt → created_at. The lowercase "pass" column is preserved by
// the same rule, no explicit column tag needed.
//
// Pass holds the bcrypt hash; validation (length, complexity) lives in
// the service layer, not on the struct. Phone is a pointer so an empty
// form field becomes SQL NULL instead of an empty string.
//
// Email has a unique index; duplicates surface as repository.ErrEmailTaken
// at the service layer.
type User struct {
	UserID    uint64    `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"uniqueIndex:idx_users_email;not null;size:255"`
	Pass      string    `gorm:"not null;size:255"`
	Name      string    `gorm:"not null;size:200"`
	Phone     *string   `gorm:"size:50"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName returns the exact DBML table name. We override GORM's default
// plural-lowercase ("users") to preserve the original "Users" casing from
// the schema. SingularTable/NamingStrategy stay untouched for column names.
func (User) TableName() string {
	return "Users"
}
