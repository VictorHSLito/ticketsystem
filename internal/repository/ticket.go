package repository

import (
	"errors"
	"fmt"

	"github.com/victorhslito/ticketsystem/internal/domain"

	"gorm.io/gorm"
)

type TicketRepository interface {
	FindTicketsByBookingID(bookingID uint) ([]*domain.Ticket, error)
}

type ticketRepositoryImpl struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepositoryImpl{db: db}
}

func (r *ticketRepositoryImpl) FindTicketsByBookingID(bookingID uint) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket

	err := r.db.Where("booking_id = ?", bookingID).Find(&tickets).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrTicketNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find tickets by booking id: %w", err)
	}
	return tickets, nil
}
