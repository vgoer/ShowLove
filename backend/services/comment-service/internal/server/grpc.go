package server

import (
	"log"

	"showlove/services/comment-service/internal/service"

	"google.golang.org/grpc"
)

type CommentServer struct{ svc *service.CommentService }

func NewCommentServer(svc *service.CommentService) *CommentServer {
	return &CommentServer{svc: svc}
}

func (s *CommentServer) RegisterGRPC(server *grpc.Server) {
	log.Println("[comment-service] gRPC server registered")
}
