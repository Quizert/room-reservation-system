package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/models"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/myerror"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"github.com/Quizert/room-reservation-system/Libs/metrics"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) RegisterUser(ctx context.Context, user *models.User) (int, error) {
	start := time.Now()
	status := "ok"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordDataBaseMetrics("Create booking", status, duration)
	}()

	query := `
		SELECT EXISTS (
			SELECT ID
			FROM users
			WHERE ChatID = $1
		);
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, user.ChatID).Scan(&exists)
	if err != nil {
		status = "fail"
		return 0, fmt.Errorf("myerror checking if user exists: %w", err)
	}
	if exists {
		status = "fail"
		return 0, fmt.Errorf("in register user: %w", myerror.ErrUserExists)
	}

	query = `
		INSERT INTO users (Username, ChatID, Password, IsHotelier)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var id int
	err = r.db.QueryRow(ctx, query, user.Username, user.ChatID, user.Password, user.IsHotelier).Scan(&id)
	if err != nil {
		status = "fail"
		return 0, fmt.Errorf("myerror inserting user: %w", err)
	}
	return id, nil
}

func (r *Repository) LoginUser(ctx context.Context, chatID string) (*models.User, error) {
	start := time.Now()
	status := "ok"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordDataBaseMetrics("Create booking", status, duration)
	}()

	query := `
		SELECT ID, Username, ChatID, Password, IsHotelier FROM users 
		WHERE ChatID = $1
	`

	var user models.User

	err := r.db.QueryRow(ctx, query, chatID).Scan(
		&user.ID,
		&user.Username,
		&user.ChatID,
		&user.Password,
		&user.IsHotelier,
	)
	if err != nil {
		status = "fail"
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("in login User: %w", myerror.ErrUserNotFound)
		}
		return nil, fmt.Errorf("myerror login user: %w", err)
	}
	return &user, nil
}

func (r *Repository) IsHotelier(ctx context.Context, userID int) (bool, error) {
	start := time.Now()
	status := "ok"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordDataBaseMetrics("Create booking", status, duration)
	}()

	query := `
		SELECT IsHotelier FROM users
		where UserID = $1
	`
	var isHotelier bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&isHotelier)
	if err != nil {
		status = "fail"
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("in isHotelier User: %w", myerror.ErrUserNotFound)
		}
		return false, fmt.Errorf("myerror isHotelier user: %w", err)
	}
	return isHotelier, nil
}

func (r *Repository) GetHotelierInformation(ctx context.Context, request *authpb.GetHotelierRequest) (*authpb.GetHotelierResponse, error) {
	start := time.Now()
	status := "ok"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RecordDataBaseMetrics("Create booking", status, duration)
	}()

	ownerID := request.OwnerID
	query := `
		SELECT Username, ChatID FROM users WHERE id = $1
	`
	var username string
	var chatID string
	err := r.db.QueryRow(ctx, query, ownerID).Scan(&username, &chatID)
	if err != nil {
		status = "fail"
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("in storage GetHotelierInformation: %w", myerror.ErrUserNotFound)
		}
		return nil, fmt.Errorf("in storage GetHotelierInformation: %w", err)
	}
	response := &authpb.GetHotelierResponse{Username: username, ChatID: chatID}
	return response, nil
}
