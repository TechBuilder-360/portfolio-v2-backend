package model

type Social struct {
	Base

	UserId string `gorm:"not null"`
	Name   string `gorm:"type:varchar(20)"`
	URL    string `gorm:"type:varchar(255)"`
}
