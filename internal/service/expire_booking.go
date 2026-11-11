package service

import (
	"log/slog"

	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/repository"
)

type ExpireBookingService struct {
	br     repository.BookingRepository
	sr     repository.SeatsRepository
	tr     repository.TicketRepository
	logger *slog.Logger
}

func NewExpireBookingService(
	bookingRepository repository.BookingRepository,
	seatRepository repository.SeatsRepository,
	ticketRepository repository.TicketRepository,
	logger *slog.Logger,
) *ExpireBookingService {
	return &ExpireBookingService{br: bookingRepository, sr: seatRepository, tr: ticketRepository, logger: logger.With("component", "expire_booking_service")}
}

func (s *ExpireBookingService) ExpireBookings(id int) {
	booking, err := s.findBooking(uint(id))
	if err != nil {
		return
	}

	s.expireBooking(booking)

	s.turnSeatsAvailable(uint(id))
}

func (s *ExpireBookingService) findBooking(bookingID uint) (*domain.Booking, error) {
	var booking *domain.Booking

	booking, err := s.br.FindByID(bookingID)
	if err != nil {
		s.logger.Info("Booking ID not foun!")
		return nil, err
	}

	return booking, nil
}

func (s *ExpireBookingService) expireBooking(booking *domain.Booking) {
	if booking.BookingStatus == domain.Pending {
		booking.BookingStatus = domain.Expired
	}

	s.br.Update(booking)
}

func (s *ExpireBookingService) turnSeatsAvailable(bookingID uint) {
	seatsIDs := s.findSeats(bookingID)

	if len(seatsIDs) == 0 {
		return
	}

	seats := s.sr.FindAll(seatsIDs)

	for _, seat := range seats {
		if seat.SeatStatus == domain.StatusBooked {
			seat.SeatStatus = domain.StatusAvailable
		}
	}

	s.sr.UpdateAll(seats)
}

func (s *ExpireBookingService) findSeats(bookingID uint) []uint {
	tickets, err := s.tr.FindTicketsByBookingID(bookingID)
	if err != nil {
		s.logger.Error("Couldn't proccess tickets")
		return nil
	}

	seatsIDs := make([]uint, 0)
	for _, ticket := range tickets {
		seatsIDs = append(seatsIDs, ticket.SeatID)
	}

	return seatsIDs
}
