package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	UserID      uint       `gorm:"index;not null"`
	Title       string     `gorm:"type:varchar(255);not null"`
	Description *string    `gorm:"type:text"`
	Status      uint8      `gorm:"type:tinyint;default:0;comment:0-待办,1-进行中,2-已完成;index"`
	Priority    uint8      `gorm:"type:tinyint;default:1;comment:1-低,2-中,3-高,4-紧急;index"`
	DueDate     *time.Time `grom:"index"`
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	User User  `gorm:"foreignKey:UserID"`
	Tags []Tag `gorm:"many2many:todo_tags;"`
}

func (Todo) TableName() string {
	return "todos"
}
