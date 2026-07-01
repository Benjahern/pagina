package models

import "time"

type Admins struct {
	AdminId   uint64    `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"uniqueIndex:idx_admins_email;not null;size:200"`
	Name      string    `gorm:"not null;size:200"`
	Pass      string    `gorm:"not null;size:255"`
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

func (Admins) TableName() string {
	return "Admins"
}
