package entities

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	Name        string
	Description string
	Price       float64
	SKU         string `gorm:"uniqueIndex"`
	CategoryID  uint
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Category    Category       `gorm:"foreignKey:CategoryID"`
	Images      []Image        `gorm:"foreignKey:ProductID"`
}
