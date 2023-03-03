package model

type Position struct {
	Base

	ExperienceId   string
	Name           string
	StartDate      string
	EndDate        string
	EmploymentType string
	InProgress     bool `gorm:"default:false"`
}
