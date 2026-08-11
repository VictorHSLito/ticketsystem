package repository

import (
	"github.com/victorhslito/ticketsystem/internal/domain"
	"gorm.io/gorm"
)

type BookingRepository interface {
	FindByID(bookingID uint) (*domain.Booking, error)
	CreateBookingTx(booking *domain.Booking, tickets []domain.Ticket, seatIDs []uint) error
}

type bookingRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepositoryImpl{db: db}
}

func (r *bookingRepositoryImpl) FindByID(bookingID uint) (*domain.Booking, error) {
	var booking domain.Booking
	if err := r.db.First(&booking, bookingID).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepositoryImpl) CreateBookingTx(booking *domain.Booking, tickets []domain.Ticket, seatIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		for i := range tickets {
			tickets[i].BookingID = booking.ID
		}

		if err := tx.Create(&tickets).Error; err != nil {
			return err
		}

		// Atualiza o status dos assentos
		if err := tx.Model(&domain.Seat{}).Where("id IN ?", seatIDs).Update("seat_status", domain.StatusBooked).Error; err != nil {
			return err
		}

		return nil
	})
}
