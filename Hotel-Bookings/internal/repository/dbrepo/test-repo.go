package dbrepo

import (
	"errors"
	"time"

	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/internal/models"
)

func (m *testDBRepo) AllUser() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("some errors")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	// if the room id is 2, then fail; otherwise, pass
	if r.RoomID == 1000 {
		return errors.New("some errors")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability exists for roomID
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	// return false if roomID = 1
	if roomID == 1 {
		return false, nil
	}
	if roomID == 1000 {
		return false, errors.New("some errors")
	}
	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of available room, if any, for give date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	testtime, _ := time.Parse("2006-01-02", "2050-01-30")
	if start == testtime {
		return nil, errors.New("some errors!")
	}
	var rooms []models.Room
	testtime, _ = time.Parse("2006-01-02", "2050-01-01")
	if start == testtime {
		rooms = append(rooms, models.Room{})
	}
	return rooms, nil
}

// GetRoomByID gets a room by ID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("Some error")
	}
	return room, nil
}

// GetUserByID returns a user by ID
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}

// UpdateUser updates user in database
func (m *testDBRepo) UpdateUser(user models.User) error {
	return nil
}

// Authenticate authenticates a user
func (m *testDBRepo) Authenticate(email, password string) (int, string, error) {
	return 1, "abc", nil
}

// AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	return []models.Reservation{}, nil
}

// AllNewReservations returns a slice of all new reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	return []models.Reservation{}, nil
}

// GetReservationByID returns a reservation by ID
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	return models.Reservation{}, nil
}

// UpdateReservation updates reservation in database
func (m *testDBRepo) UpdateReservation(res models.Reservation) error {
	return nil
}

// DeleteReservation deletes a reservation by ID
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// UpdateProcessedForReservation updates processed of reservation by ID
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}

// AllRooms returns all rooms
func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRestrictionsForRoomByDate returns all restrictions for a room by date range
func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction
	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction
func (m *testDBRepo) InsertBlockForRoom(id int, date time.Time) error {
	return nil
}

// DeleteBlockByID deletes a room restriction by ID
func (m *testDBRepo) DeleteBlockByID(id int) error {
	return nil
}
