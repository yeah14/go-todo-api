package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	Username     string  `gorm:"type:varchar(50);uniqueIndex;not null"`
	Email        string  `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string  `gorm:"type:varchar(255);not null"`
	AvatarURL    *string `gorm:"type:varchar(255)"`
	Status       uint8   `gorm:"type:tinyint;default:1;comment:0-禁用,1-正常"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
