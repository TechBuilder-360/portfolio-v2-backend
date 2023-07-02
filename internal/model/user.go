package model

type User struct {
	Base

	AccountID   string  `gorm:"not null"`
	FirstName   string  `gorm:"type:varchar(50)"`
	LastName    string  `gorm:"type:varchar(50)"`
	MiddleName  string  `gorm:"type:varchar(50)"`
	Email       string  `gorm:"type:varchar(150);not null;<-:create"`
	Gender      string  `gorm:"type:varchar(6)"`
	PhoneNumber string  `gorm:"type:varchar(20)"`
	Bio         string  `gorm:"type:text"`
	DateOfBirth string  `gorm:"type:varchar(10)"`
	ProfilePix  string  `gorm:"type:varchar(255)"`
	Profession  string  `gorm:"type:varchar(50)"`
	HasProfile  bool    `gorm:"default:false"`
	Account     Account `gorm:"constraint:OnDelete:Set null;"`
}
