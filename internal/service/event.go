package service

import (
	"strconv"
	"strings"

	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/repository"
)

type EventService struct {
	r repository.EventRepository
}

func NewEventService(r repository.EventRepository) *EventService {
	return &EventService{r: r}
}

func (s *EventService) CreateEvent(name, description string, quantityOfSeats int) uint {
	seats := s.createSeats(quantityOfSeats)

	event := &domain.Event{
		Name:        name,
		Description: description,
		Seats:       seats,
	}

	event, err := s.r.Create(event)
	if err != nil {
		return 0
	}

	return event.ID
}

func (s *EventService) createSeats(quantity int) []domain.Seat {
	seats := make([]domain.Seat, 0, quantity)

	for i := range quantity {
		seats = append(seats, domain.Seat{
			SeatCode:   "SA" + strconv.Itoa(i),
			SeatStatus: domain.StatusAvailable,
		})
	}
	return seats
}

func (s *EventService) ListEvents(page, pageSize int) ([]*domain.Event, int64) {
	return s.r.ListAll(page, pageSize)
}

func (s *EventService) ListByID(eventID uint) *domain.Event {
	event, err := s.r.ListEventByID(eventID)
	if err != nil {
		return nil
	}

	return event
}

func (s *EventService) ListAllSeats(eventID uint, page, pageSize int) ([]*domain.Seat, int64) {
	return s.r.ListAllSeats(eventID, page, pageSize)
}

func (s *EventService) ListSeatsByStatus(eventID uint, seatStatus string) []*domain.Seat {
	seatStatus = strings.ToUpper(seatStatus)

	if seatStatus == string(domain.StatusBooked) {
		return s.r.ListSeatsByStatus(eventID, domain.StatusBooked)
	} else {
		return s.r.ListSeatsByStatus(eventID, domain.StatusAvailable)
	}
}
