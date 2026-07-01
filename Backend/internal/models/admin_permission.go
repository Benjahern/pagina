package models

import "time"

// AdminPermission is the join table between Admins and Permission. CASCADE on
// both FKs: deleting an admin strips their perms; deleting a permission
// removes references.
type AdminPermission struct {
	AdminPermissionID uint64    `gorm:"primaryKey;autoIncrement"`
	AdminID           uint64    `gorm:"not null;uniqueIndex:uq_admin_permission,priority:1"`
	PermissionID      uint64    `gorm:"not null;uniqueIndex:uq_admin_permission,priority:2"`
	CreatedAt         time.Time `gorm:"not null;autoCreateTime"`
}

func (AdminPermission) TableName() string {
	return "Admin_Permission"
}
