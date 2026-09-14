package service

import (
	"errors"
	"log"

	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/repository"
)

type ConfirmBookingService struct {
	br repository.BookingRepository
}

func NewConfirmBookingService(bookingRepository repository.BookingRepository) *ConfirmBookingService {
	return &ConfirmBookingService{br: bookingRepository}
}

func (s *ConfirmBookingService) ConfirmBooking(bookingID uint) error {
	booking, err := s.br.FindByID(bookingID)
	if err != nil {
		log.Println("Booking not found!")
		return err
	}

	if booking.BookingStatus != domain.Pending {
		log.Println("[ERROR] - Tried to execute an unsupported operation")
		return errors.ErrUnsupported
	}

	booking.BookingStatus = domain.Confirmed
	s.br.Update(booking)
	return nil
}
