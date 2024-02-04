package repository

import model "github.com/rashidalam9678/project-management-software-server/internal/models"

// import model "github.com/rashidalam9678/project-management-software-server/internal/models"

// "time"

// "github.com/rashidalam9678/project-management-software-server/internal/models"

type Database interface{
	GetUserByEmail(email string) (model.User,error)
	InsertUser(email,externalId string)(uint,error)
	GetUserByExternalID(externalId string) (model.User, error)
	// InsertReservation(res models.Reservation) (int,error)
	// InsertRoomRestriction( res models.RoomRestriction)(error)
	// SearchAvailablityByDatesByRoomId(start time.Time,end time.Time , roomId int)(bool, error)
	// SearchAvailablityForAllRooms(start, end time.Time) ([]models.Room, error)
	// GetRoomById(id int)(models.Room, error)
	// GetUserById(id int)(models.User, error)
	// UpdateUserById(u models.User)( error)
	// Authenticate(email,testPassword string)(int, string,error)
	// AllReservations() ([]models.Reservation, error)
	// AllNewReservations() ([]models.Reservation, error)
	// GetReservationById(id int) (models.Reservation, error) 
	// UpdateReservation(u models.Reservation)( error)
	// DeleteReservationById(id int)( error)
	// UpdateProcessed(id,processed int)( error)
	// InsertUser(user models.User) (int, error)

}