package model

import (
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/util"
	"gorm.io/gorm"
	"time"
)

// Base definition
type Base struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = util.GenerateUUID()
	}
	return
}
