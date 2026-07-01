package models

import "time"

// Impresora3d is the internal table for the owner's 3D printers. Not
// exposed publicly. Used by the cost calculator to compute printing
// cost = (filament grams × price/kg) + (hours × electricity_cost/hour).
//
// Money fields are int64 in CLP (no decimals). error_margin is a
// percent (numeric(5,2)) so float64 — different domain. NOT NULL with
// default 10.00 in DBML, so the field is non-pointer.
//
// Table/column casing is preserved from DBML: lowercase "impresora3d"
// table and "impresora3d_id" PK. NamingStrategy converts the Go field
// Impresora3dID → impresora3d_id via the camelCase→snake_case rule.
//
// UsefulLifeHours int (not uint) so the field can be NULL — DB schema
// has no NOT NULL constraint and we want to allow "TBD" printers.
type Impresora3d struct {
	Impresora3dID          uint64  `gorm:"primaryKey;autoIncrement"`
	Name                   *string `gorm:"size:200"`
	ElectricityCostPerHour *int64  `gorm:"type:numeric(12,4)"`
	CostReparation         *int64  `gorm:"type:numeric(12,2)"`
	ErrorMargin            float64 `gorm:"type:numeric(5,2);not null;default:10.00"`
	UsefulLifeHours        *int
	IsActive               bool      `gorm:"not null;default:true"`
	CreatedAt              time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt              time.Time `gorm:"not null;autoUpdateTime"`
}

func (Impresora3d) TableName() string {
	return "impresora3d"
}
