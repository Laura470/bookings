package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/Laura470/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	//if the room id is2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("some error with the roomid in insert reservation")
	}
	return 1, nil
}

//§InsertRoomREstriction insert a room restriction into the dataabase
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error with the room restriction")
	}
	return nil
}

//SearchAvailabilityByDatesByRoomID ritorna true se c'è disponibilità per un a particolare stanza, false se no c'è
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	// set up a test time
	//t rappresenta la data del 2049-12-31
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2060-01-01 as start
	//datetoFail rappresenta la data del 2060-01-01
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	//ritorna errore, quindi qualcosa andato male nella query
	if start == testDateToFail {
		return false, errors.New("some error")
	}

	// if the start date is after 2049-12-31, then return false,
	// indicating no availability;
	//ritorna false, quindi no availability
	if start.After(t) {
		return false, nil
	}

	// otherwise, we have availability
	return true, nil

}

//SearchAvailabilityForAllRooms return a slice of available rooms, if any, for given range date
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	//creo uno slice di Room (struct)
	var rooms []models.Room

	// set up a test time
	//t rappresenta la data del 2049-12-31
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2060-01-01 as start
	//datetoFail rappresenta la data del 2060-01-01
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	// if the start date is equal to 2060-01-01,
	//ritorna errore, quindi qualcosa andato male nella query
	if start == testDateToFail {
		return rooms, errors.New("some error")
	}

	// if the start date after 2049-12-31,
	//ritorna slice vuota, quindi no availability
	if start.After(t) {
		return rooms, nil
	}

	//roomslen >0  c'è dispobilibilità
	//creo una Room
	room := models.Room{
		ID:       1,
		RoomName: "stanza",
	}
	rooms = append(rooms, room)

	return rooms, nil
}

//già cae ci sono ritorno tutto, no solo in nome della stanza
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}
