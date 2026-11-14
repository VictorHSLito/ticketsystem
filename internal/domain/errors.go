package domain

import "errors"

var (
	ErrUserNotFound               = errors.New("user not found")
	ErrUserAlreadyExists          = errors.New("user already exists")
	ErrBookingNotFound            = errors.New("booking not found")
	ErrEventAlreadyExists         = errors.New("event with this id already exists")
	ErrEventNotFound              = errors.New("event not found")
	ErrTicketNotFound             = errors.New("ticket not found")
	ErrSeatNotFound               = errors.New("seat not found")
	ErrSeatAlreadyBooked          = errors.New("seat already booked")
	ErrSeatNotValidOrNotAvailable = errors.New("no valid or available seats")
	ErrBookingInvalidState        = errors.New("invalid booking state")
	ErrSeatsNotAvailable          = errors.New("one seat or more are not available")
	ErrAdminSetupDisabled         = errors.New("admin already exists and another one cannot be created")
)
