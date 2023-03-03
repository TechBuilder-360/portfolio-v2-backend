package model

type Education struct {
	Base

	Name       string
	Degree     string
	Course     string
	StartDate  string
	EndDate    string
	Location   string
	InProgress bool `gorm:"default:false"`
}
