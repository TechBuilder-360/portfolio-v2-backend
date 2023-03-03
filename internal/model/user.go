package model

type User struct {
	Base

	AccountId   string `gorm:"not null"`
	FirstName   string `gorm:"type:varchar(50)"`
	LastName    string `gorm:"type:varchar(50)"`
	MiddleName  string `gorm:"type:varchar(50)"`
	Email       string `gorm:"type:varchar(150)"`
	Gender      string `gorm:"type:varchar(6)"`
	PhoneNumber string `gorm:"type:varchar(20)"`
	Bio         string `gorm:"type:text"`
	UserName    string `gorm:"type:varchar(50)"`
	DateOfBirth string `gorm:"type:varchar(10)"`
	ProfilePix  string `gorm:"type:varchar(255)"`
	Profession  string `gorm:"type:varchar(50)"`
}
