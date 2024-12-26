package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/myerror"
	postgresql2 "github.com/Quizert/room-reservation-system/HotelSvc/internal/repository/postgresql"
	service2 "github.com/Quizert/room-reservation-system/HotelSvc/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	handler "github.com/Quizert/room-reservation-system/HotelSvc/api/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v\n", err)
	}
	defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}
	fmt.Println("Успешное подключение к базе данных!")

	hotelRepo := postgresql2.NewPostgresHotelRepository(db)

	hotelService := service2.NewHotelService(hotelRepo)

	roomRepo := postgresql2.NewPostgresRoomRepository(db)

	roomService := service2.NewRoomService(roomRepo)

	ownerRepo := postgresql2.NewPostgresOwnerRepository(db)

	ownerService := service2.NewOwnerService(ownerRepo)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := startHTTPServer(hotelService); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Запуск gRPC сервера в отдельной горутине
	go func() {
		defer wg.Done()
		if err := startGRPCServer(roomService, ownerService); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Ожидание завершения серверов
	wg.Wait()

	//if myerror := startHTTPServer(hotelService); myerror != nil {
	//	log.Fatalf("Failed to start HTTP server: %v", myerror)
	//}

	//if myerror := startGRPCServer(roomService); myerror != nil {
	//	log.Fatalf("Failed to start GRPC server: %v", myerror)
	//}
}

// initDB инициализирует подключение к базе данных PostgreSQL
func initDB() (*sql.DB, error) {
	dbHost := os.Getenv("HOTEL_DB_HOST")
	dbPort := os.Getenv("HOTEL_DB_PORT")
	dbUser := os.Getenv("HOTEL_DB_USER")
	dbPassword := os.Getenv("HOTEL_DB_PASSWORD")
	dbName := os.Getenv("HOTEL_DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	return sql.Open("postgres", dsn)
}

// startHTTPServer запускает HTTP сервер для обработки REST-запросов
func startHTTPServer(hotelService *service2.HotelService) error {
	mux := http.NewServeMux()
	handler.RegisterHotelRoutes(mux, hotelService)

	addr := ":" + os.Getenv("HOTEL_HTTP_PORT")
	log.Printf("Starting HTTP server on %s...", addr)
	return http.ListenAndServe(addr, mux)
}

type server struct {
	hotelpb.UnimplementedHotelServiceServer
	roomService  *service2.RoomService
	ownerService *service2.OwnerService
}

func (s *server) GetRoomsByHotelId(ctx context.Context, req *hotelpb.GetRoomsRequest) (*hotelpb.GetRoomsResponse, error) {
	hotelId := req.GetHotelId()
	rooms, err := s.roomService.GetRoomsByHotelId(int(hotelId))
	if err != nil {
		return nil, err
	}

	return &hotelpb.GetRoomsResponse{Rooms: rooms}, nil

}

func (s *server) GetOwnerIdByHotelId(ctx context.Context, req *hotelpb.GetOwnerIdRequest) (*hotelpb.GetOwnerIdResponse, error) {
	hotelId := req.GetId()
	ownerId, err := s.ownerService.GetOwnerIdByHotelId(ctx, int(hotelId))
	if err != nil {
		if errors.Is(err, myerror.ErrHotelNotFound) {
			return nil, status.Error(codes.NotFound, "hotel not found")
		}
		return nil, fmt.Errorf("in server: %w", err)
	}
	return &hotelpb.GetOwnerIdResponse{OwnerId: int32(ownerId)}, nil
}

func startGRPCServer(roomService *service2.RoomService, ownerService *service2.OwnerService) error {
	addr := ":" + os.Getenv("HOTEL_GRPC_PORT")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

	s := grpc.NewServer()
	hotelpb.RegisterHotelServiceServer(s, &server{roomService: roomService, ownerService: ownerService})

	reflection.Register(s)

	log.Printf("Starting GRPC server on %s", addr)
	return s.Serve(lis)
}
