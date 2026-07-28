package domain

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Name        string
	Description string
	Seats       []Seat
}

type SeatStatus string

const (
	StatusAvailable SeatStatus = "AVAILABLE"
	StatusBooked    SeatStatus = "BOOKED"
)

type Seat struct {
	ID         uint
	EventID    uint
	SeatCode   string
	SeatStatus SeatStatus
}
