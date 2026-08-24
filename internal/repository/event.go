package repository

import (
	"errors"

	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"gorm.io/gorm"
)

type EventRepository interface {
	ListEventByID(eventID uint) (*domain.Event, error)
	ListAll(page, pageSize int) ([]*domain.Event, int64)
	Create(event *domain.Event) (*domain.Event, error)
	ListSeatsByStatus(eventID uint, seatStatus domain.SeatStatus) []*domain.Seat
	ListAllSeats(eventID uint, page, pageSize int) ([]*domain.Seat, int64)
}

type eventRepositoryImpl struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepositoryImpl{db: db}
}

func (r *eventRepositoryImpl) ListEventByID(eventID uint) (*domain.Event, error) {
	var event domain.Event

	err := r.db.Preload("Seats").Where("id = ?", eventID).First(&event).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &event, nil
}

func (r *eventRepositoryImpl) ListAll(page, pageSize int) ([]*domain.Event, int64) {
	var events []*domain.Event
	var totalElements int64

	r.db.Model(&domain.Event{}).Count(&totalElements)
	r.db.Scopes(database.Paginate(page, pageSize)).Find(&events)

	return events, totalElements
}

func (r *eventRepositoryImpl) Create(event *domain.Event) (*domain.Event, error) {
	err := r.db.Create(event).Error

	if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, err
	}

	return event, nil
}

func (r *eventRepositoryImpl) ListSeatsByStatus(eventID uint, seatStatus domain.SeatStatus) []*domain.Seat {
	var seats []*domain.Seat
	r.db.Where("event_id = ? and seat_status = ?", eventID, seatStatus).Find(&seats)

	return seats
}

func (r *eventRepositoryImpl) ListAllSeats(eventID uint, page, pageSize int) ([]*domain.Seat, int64) {
	var seats []*domain.Seat
	var totalElements int64

	r.db.Model(&domain.Seat{}).Count(&totalElements)
	r.db.Where("event_id = ?", eventID).Preload("Seats").Scopes(database.Paginate(page, pageSize)).Find(&seats)

	return seats, totalElements
}
