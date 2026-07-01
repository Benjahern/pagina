package models

import "time"

// Item is a product listed in the public storefront. See CLAUDE.md §
// "Modelo de dominio — Capa pública".
//
// Status holds the value column-wise (varchar(50)); the valid set lives
// in item_status.go. Service layer is responsible for not inserting "".
type Item struct {
	ItemID          uint64    `gorm:"primaryKey;autoIncrement"`
	Sku             string    `gorm:"uniqueIndex:idx_items_sku;not null;size:100"`
	Slug            string    `gorm:"uniqueIndex:idx_items_slug;not null;size:500"`
	Name            string    `gorm:"not null;size:200"`
	Description     *string
	Price           int64     `gorm:"type:numeric(12,2);not null"`
	Cost            int64     `gorm:"type:numeric(12,2);not null"`
	Stock           int       `gorm:"not null;default:0"`
	Backorder       bool      `gorm:"not null;default:false"`
	Status          string    `gorm:"index:idx_items_status,index:idx_items_status_category,priority:1;not null;size:50;default:activo"`
	CategoryID      uint64    `gorm:"index:idx_items_category,index:idx_items_status_category,priority:2;not null"`
	Brand           *string   `gorm:"size:200"`
	Color           *string   `gorm:"size:100"`
	ImageURL        *string   `gorm:"size:500"`
	Items3DID       *uint64   `gorm:"column:items3d_id"`
	MetaTitle       *string   `gorm:"size:200"`
	MetaDescription *string
	ViewCount       int       `gorm:"not null;default:0"`
	IsFeatured      bool      `gorm:"not null;default:false;index:idx_items_featured"`
	CreatedAt       time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"not null;autoUpdateTime"`
}

func (Item) TableName() string {
	return "Items"
}
