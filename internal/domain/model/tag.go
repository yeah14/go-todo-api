package model

import (
	"time"
)

type Tag struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;type:varchar(255);not null"`
	Color     string `gorm:"type:varchar(7)"`
	UserID    uint   `gorm:"not null;index:idx_user_tag,unique;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_id"`
	CreatedAt time.Time
	UpdatedAt time.Time

	User User   `gorm:"foreignKey:UserID"`
	Todo []Todo `gorm:"many2many:todo_tags;"`
}

func (Tag) TableName() string {
	return "tags"
}
