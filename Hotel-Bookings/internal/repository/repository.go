package repository

import (
	"time"

	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/internal/models"
)

type DatabaseRepo interface {
	AllUser() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
}
