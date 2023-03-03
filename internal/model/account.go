package model

type Account struct {
	Base

	Email     string `gorm:"type:varchar(150)"`
	Password  string `gorm:"type:varchar(255)"`
	LastLogin string `gorm:"type:varchar(50)"`
	AuthType  string `gorm:"type:varchar(150)"` //
}
