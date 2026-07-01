package models

import "time"

// User represents a customer of the storefront. Column names come from
// GORM's NamingStrategy (CamelCase → snake_case): UserID → user_id,
// CreatedAt → created_at. The lowercase "pass" column is preserved by
// the same rule, no explicit column tag needed.
//
// Pass holds the bcrypt hash; validation (length, complexity) lives in
// the service layer, not on the struct. Phone is required at the DB
// level (NOT NULL) — an empty form is rejected by the binding tag,
// not silently stored as NULL.
//
// Rut is the Chilean national ID with check digit (e.g. "12345678-9"
// or "12345678-K"), unique per person. Format normalization and check
// digit verification live in the service layer — the struct only
// guarantees storage. UniqueIndex on rut so duplicate registrations
// surface as a repository conflict error at the service layer.
//
// Email has a unique index; duplicates surface as repository.ErrEmailTaken
// at the service layer.
type User struct {
	UserID    uint64    `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"uniqueIndex:idx_users_email;not null;size:255"`
	Pass      string    `gorm:"not null;size:255"`
	Name      string    `gorm:"not null;size:200"`
	Phone     string    `gorm:"not null;size:20"`
	Rut       string    `gorm:"uniqueIndex:idx_users_rut;not null;size:12"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName returns the exact DBML table name. We override GORM's default
// plural-lowercase ("users") to preserve the original "Users" casing from
// the schema. SingularTable/NamingStrategy stay untouched for column names.
func (User) TableName() string {
	return "Users"
}
