package repository

import (
	"github.com/victorhslito/ticketsystem/internal/domain"
	"gorm.io/gorm"
)

type SeatsRepository interface {
	FindAll(seatsID []uint) []*domain.Seat
	Update(seat *domain.Seat)
	UpdateAll(seats []*domain.Seat)
}

type seatsRepositoryImpl struct {
	db *gorm.DB
}

func NewSeatsRepository(db *gorm.DB) SeatsRepository {
	return &seatsRepositoryImpl{db: db}
}

func (r *seatsRepositoryImpl) FindAll(seatsID []uint) []*domain.Seat {
	var seats []*domain.Seat

	if len(seatsID) == 0 {
		return nil
	}

	r.db.Where("id IN ?", seatsID).Find(&seats)

	return seats
}

func (r *seatsRepositoryImpl) Update(seat *domain.Seat) {
	r.db.Save(seat)
}

func (r *seatsRepositoryImpl) UpdateAll(seats []*domain.Seat) {
	for _, seat := range seats {
		r.db.Save(seat)
	}
}
