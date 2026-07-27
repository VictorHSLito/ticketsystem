package domain

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Name        string
	Description string
}
