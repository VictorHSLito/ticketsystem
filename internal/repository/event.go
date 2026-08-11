package repository

import (
	"errors"

	"github.com/victorhslito/ticketsystem/internal/domain"
	"gorm.io/gorm"
)

type EventRepository interface {
	ListEventByID(eventID uint) (*domain.Event, error)
	ListAll() []*domain.Event
	Create(event *domain.Event) (*domain.Event, error)
	ListSeatsAvailables(eventID uint) []*domain.Seat
}

type eventRepositoryImpl struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepositoryImpl{db: db}
}

func (r *eventRepositoryImpl) ListEventByID(eventID uint) (*domain.Event, error) {
	var event domain.Event

	err := r.db.Where("id = ?", eventID).First(&event).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &event, nil
}

func (r *eventRepositoryImpl) ListAll() []*domain.Event {
	var events []*domain.Event

	r.db.Find(&events)
	return events
}

func (r *eventRepositoryImpl) Create(event *domain.Event) (*domain.Event, error) {
	err := r.db.Create(event).Error

	if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, err
	}

	return event, nil
}

func (r *eventRepositoryImpl) ListSeatsAvailables(eventID uint) []*domain.Seat {
	var seats []*domain.Seat
	r.db.Where("event_id = ? and seat_status = ?", eventID, domain.StatusAvailable).Find(&seats)

	return seats
}
