package entities

import (
	"time"
)

type Image struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	ProductID uint
	URL       string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Product   Product   `gorm:"foreignKey:ProductID"`
}
