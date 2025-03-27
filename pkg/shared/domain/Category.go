package domain

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
