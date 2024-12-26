package postgresql

import (
	"database/sql"

	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
)

type PostgresRoomRepository struct {
	db *sql.DB
}

func NewPostgresRoomRepository(db *sql.DB) *PostgresRoomRepository {
	return &PostgresRoomRepository{db: db}
}

func (repo *PostgresRoomRepository) GetRoomsByHotelId(id int) ([]*hotelpb.Room, error) {
	rows, err := repo.db.Query("SELECT r.id AS RoomId, r.HotelId, r.RoomTypeId, r.Number, rt.BasePrice AS Cost FROM rooms r JOIN room_type rt ON r.RoomTypeId = rt.id WHERE r.HotelId = $1;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*hotelpb.Room
	for rows.Next() {
		var room hotelpb.Room
		if err := rows.Scan(&room.Id, &room.HotelId, &room.RoomTypeId, &room.Number, &room.Cost); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	return rooms, rows.Err()

}
