package grpc

import (
	"context"
	"errors"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/controller"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/myerror"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
	authSvc controller.AuthService
	trace   trace.Tracer
	Addr    string
}

func NewServer(authSvc controller.AuthService, addr string, trace trace.Tracer) *Server {
	return &Server{authSvc: authSvc, Addr: addr, trace: trace}
}

func (s *Server) GetHotelierInformation(ctx context.Context, req *authpb.GetHotelierRequest) (*authpb.GetHotelierResponse, error) {
	// Начинаем новый спан, контекст автоматически будет извлечён из интерцептора
	ctx, span := s.trace.Start(ctx, "GetHotelierInformation")
	defer span.End()

	response, err := s.authSvc.GetHotelierInformation(ctx, req)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, myerror.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, myerror.ErrUserNotFound.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return response, nil
}
