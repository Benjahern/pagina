package models

import "time"

// Permission represents a granular access key (e.g. "products.create",
// "users.delete") referenced by Admin_Permission. Renamed from "acces"
// (typo) in the original schema; the column stays "access".
//
// Permission rows are effectively immutable: there is no updated_at
// column in the DBML (line 70-75), so this struct deliberately has
// no UpdatedAt. GORM's autoUpdateTime would break selects/updates.
type Permission struct {
	PermissionID uint64    `gorm:"primaryKey;autoIncrement"`
	Access       string    `gorm:"uniqueIndex:idx_permission_access;not null;size:200"`
	Description  *string
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
}

func (Permission) TableName() string {
	return "Permission"
}
