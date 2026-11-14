package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/repository"
)

type ConfirmBookingService struct {
	br     repository.BookingRepository
	logger *slog.Logger
}

func NewConfirmBookingService(bookingRepository repository.BookingRepository, logger *slog.Logger) *ConfirmBookingService {
	return &ConfirmBookingService{br: bookingRepository, logger: logger.With("component", "confirm_booking_service")}
}

func (s *ConfirmBookingService) ConfirmBooking(bookingID uint) error {
	booking, err := s.br.FindByID(bookingID)
	if err != nil {
		return fmt.Errorf("confirm booking: %w", err)
	}

	if booking.BookingStatus != domain.Pending {
		s.logger.Error("Tried to execute an unsupported operation")
		return errors.ErrUnsupported
	}

	booking.BookingStatus = domain.Confirmed
	s.br.Update(booking)
	return nil
}
