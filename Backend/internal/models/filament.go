package models

import "time"

// Filament tracks 3D-printer filament inventory (the owner's private
// stock — see CLAUDE.md § "Capa personal/herramientas"). Money fields
// are int64 in whole CLP (no decimals per project decision), persisted
// as numeric(12,2) in PG for compatibility with the existing schema.
//
// The PK is NOT auto-increment: filament_id is intended as a stable,
// possibly-meaningful identifier assigned by the operator. Schema note:
// the DBML declares [pk, not null] without "increment".
//
// No updated_at — same reason as Permission (immutable-ish reference data
// for cost calculations; if you start editing filaments regularly, add
// updated_at to both DBML and this struct).
type Filament struct {
	FilamentID   uint64    `gorm:"primaryKey"`
	Name         string    `gorm:"not null;size:200"`
	Slug         *string   `gorm:"uniqueIndex:idx_filament_slug;size:200"`
	CostKilogram *int64    `gorm:"type:numeric(12,2)"`
	Color        *string   `gorm:"size:100"`
	Brand        *string   `gorm:"size:200"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
}

func (Filament) TableName() string {
	return "Filament"
}
