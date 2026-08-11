package service

import (
	"strconv"

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

func (s *EventService) ListEvents() []*domain.Event {
	return s.r.ListAll()
}

func (s *EventService) ListByID(eventID uint) *domain.Event {
	event, err := s.r.ListEventByID(eventID)
	if err != nil {
		return nil
	}

	return event
}

func (s *EventService) ListSeatsAvailables(eventID uint) []*domain.Seat {
	return s.r.ListSeatsAvailables(eventID)
}
