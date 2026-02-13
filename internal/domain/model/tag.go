package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;type:varchar(255);not null"`
	Color     string `gorm:"type:varchar(7)"`
	UserID    uint   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User User   `gorm:"foreignKey:UserID"`
	Todo []Todo `gorm:"many2many:todo_tags;"`
}

func (Tag) TableName() string {
	return "tags"
}
