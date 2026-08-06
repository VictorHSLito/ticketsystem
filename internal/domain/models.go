package domain

import (
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string
	Description string
	Seats       []Seat
}

type (
	SeatStatus    string
	BookingStatus string
)

const (
	StatusAvailable SeatStatus    = "AVAILABLE"
	StatusBooked    SeatStatus    = "BOOKED"
	Pending         BookingStatus = "PENDING"
	Confirmed       BookingStatus = "CONFIRMED"
	Expired         BookingStatus = "EXPIRED"
	Canceled        BookingStatus = "CANCELED"
)

type Seat struct {
	ID         uint
	EventID    uint
	SeatCode   string
	SeatStatus SeatStatus
}

type User struct {
	gorm.Model
	Username string
	Email    string `gorm:"unique"`
	Password string
	Bookings []Booking
}

type Booking struct {
	gorm.Model
	UserID        uint
	BookingStatus BookingStatus
}

type Ticket struct {
	ID         uint
	ExternalID string
	BookingID  uint
	SeatID     uint
}
