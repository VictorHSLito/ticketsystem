package repository

import (
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

	if err := r.db.Where("booking_id = ?", bookingID).Find(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}
