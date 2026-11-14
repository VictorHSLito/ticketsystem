package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/repository"
)

type BookingService struct {
	bookingRepo repository.BookingRepository
	userRepo    repository.UserRepository
	eventRepo   repository.EventRepository
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository,
) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		userRepo:    userRepo,
		eventRepo:   eventRepo,
	}
}

func (s *BookingService) FindByID(bookingID uint) (*domain.Booking, error) {
	return s.bookingRepo.FindByID(bookingID)
}

func (s *BookingService) CreateBooking(userID uint, eventID uint, seatCodes []string) (*domain.Booking, error) {
	// Valida usuário
	user, err := s.userRepo.ListByID(userID)
	if err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	// Valida evento
	event, err := s.eventRepo.ListEventByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	// Valida e mapeia os assentos
	bookedSeats, err := verifySeats(event.Seats, seatCodes)
	if err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	if len(bookedSeats) == 0 {
		return nil, domain.ErrSeatNotValidOrNotAvailable
	}

	booking := &domain.Booking{
		UserID:        userID,
		BookingStatus: domain.Pending,
	}

	tickets := make([]domain.Ticket, 0, len(bookedSeats))
	var seatIDs []uint
	for seatID := range bookedSeats {
		seatIDs = append(seatIDs, seatID)
		tickets = append(tickets, domain.Ticket{
			ExternalID: uuid.NewString(),
			SeatID:     seatID,
		})
	}

	if err := s.bookingRepo.CreateBookingTx(booking, tickets, seatIDs); err != nil {
		return nil, err
	}

	return booking, nil
}

func verifySeats(seats []domain.Seat, seatCodes []string) (map[uint]string, error) {
	seatMap := make(map[string]domain.Seat)

	for _, seat := range seats {
		seatMap[seat.SeatCode] = seat
	}

	bookingSeats := make(map[uint]string)
	for _, code := range seatCodes {
		seat, exists := seatMap[code]
		if !exists {
			return nil, fmt.Errorf("seat %q: %w", code, domain.ErrSeatNotFound)
		}

		if seat.SeatStatus != domain.StatusAvailable {
			return nil, fmt.Errorf("seat %q: %w", code, domain.ErrSeatAlreadyBooked)
		}

		bookingSeats[seat.ID] = code
	}
	return bookingSeats, nil
}
