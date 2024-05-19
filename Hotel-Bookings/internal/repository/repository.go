package repository

import "github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/internal/models"

type DatabaseRepo interface {
	AllUser() bool

	InsertReservation(res models.Reservation) error
}
