// Package server implements the gRPC server for the user service.
package server

import (
	"context"
	"log"

	"showlove/services/user-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// UserServer wraps the user service for gRPC.
type UserServer struct {
	svc *service.UserService
	db  *gorm.DB
}

// RegisterGRPC registers the user service gRPC server.
// Note: Full gRPC registration requires generated proto code.
// This placeholder will be replaced when protoc is available.
func (s *UserServer) RegisterGRPC(server *grpc.Server) {
	// TODO: Register generated gRPC service implementation
	// pb.RegisterUserServiceServer(server, s)
	log.Println("[user-service] gRPC server registered")
}

// Register handles user registration.
func (s *UserServer) Register(ctx context.Context, req interface{}) (interface{}, error) {
	// Placeholder: will use generated Request/Response types
	return nil, status.Error(codes.Unimplemented, "gRPC stubs not yet generated")
}

// Login handles user login.
func (s *UserServer) Login(ctx context.Context, req interface{}) (interface{}, error) {
	return nil, status.Error(codes.Unimplemented, "gRPC stubs not yet generated")
}

// NewUserServer creates a new gRPC server instance.
func NewUserServer(db *gorm.DB, svc *service.UserService) *UserServer {
	return &UserServer{db: db, svc: svc}
}
