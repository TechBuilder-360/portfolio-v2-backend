package model

type Experience struct {
	Base

	UserId      string `gorm:"not null"`
	Institution string `gorm:"not null"`
	InProgress  bool   `gorm:"default:false"`
	Location    string `gorm:"type:varchar(100)"`
}
