package models

import "time"

// Item3D is an internal printing configuration used to compute production
// costs for Items that the owner prints in-house. See CLAUDE.md §
// "Modelo de dominio (resumen) — Personal/herramientas".
type Item3D struct {
	Item3DID      uint64    `gorm:"primaryKey;autoIncrement;column:item3d_id"`
	Name          string    `gorm:"not null;size:200"`
	Impresora3dID *uint64   `gorm:"index:idx_items3d_printer"`
	FilamentGrams *float64
	Hours         int    `gorm:"not null;default:0"`
	Minutes       int    `gorm:"not null;default:0"`
	ExtraCost     *int64 `gorm:"type:numeric(12,2)"`
	Cost          *int64 `gorm:"type:numeric(12,2)"`
	FilamentID    *uint64
	CreatedAt     time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"not null;autoUpdateTime"`
}

func (Item3D) TableName() string {
	return "Items3D"
}
