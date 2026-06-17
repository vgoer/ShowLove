// Package server implements the gRPC server for the post service.
package server

import (
	"log"

	"showlove/services/post-service/internal/service"

	"google.golang.org/grpc"
)

// PostServer wraps the post service for gRPC.
type PostServer struct {
	svc *service.PostService
}

// NewPostServer creates a new gRPC server instance.
func NewPostServer(svc *service.PostService) *PostServer {
	return &PostServer{svc: svc}
}

// RegisterGRPC registers the post service gRPC server.
// Full registration requires generated proto code; this is a skeleton.
func (s *PostServer) RegisterGRPC(server *grpc.Server) {
	// TODO: pb.RegisterPostServiceServer(server, s)
	log.Println("[post-service] gRPC server registered")
}
