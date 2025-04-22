package entities

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"uniqueIndex"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Products  []Product      `gorm:"foreignKey:CategoryID"`
}
