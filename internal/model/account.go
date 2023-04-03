package model

import (
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/types"
	"time"
)

type Account struct {
	Base

	UserName      string         `gorm:"type:varchar(50)"`
	Email         string         `gorm:"type:varchar(150)"`
	Password      string         `gorm:"type:varchar(255)"`
	AuthType      types.AuthType `gorm:"type:varchar(150)"`
	EmailVerified bool           `gorm:"default:false;<-:create"`
	LastLogin     time.Time
}
