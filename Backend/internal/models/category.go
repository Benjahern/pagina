package models

import "time"

// Category represents a product category. Self-referencing via ParentID
// (top-level categories have NULL parent_id; subcategories point at their
// parent's category_id).
//
// Slug is unique and used in URLs (/categoria/<slug>); the auto-generation
// hook (BeforeCreate) is pending — for now the service layer must populate
// it explicitly. See Backend/CLAUDE.md § Pendientes.
type Category struct {
	CategoryID      uint64    `gorm:"primaryKey;autoIncrement"`
	Name            string    `gorm:"not null;size:200"`
	Slug            string    `gorm:"uniqueIndex:idx_category_slug;not null;size:200"`
	Description     *string
	MetaTitle       *string   `gorm:"size:200"`
	MetaDescription *string
	ParentID        *uint64   `gorm:"index:idx_category_parent"`
	CreatedAt       time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName returns the exact DBML table name (singular: "Category").
func (Category) TableName() string {
	return "Category"
}
