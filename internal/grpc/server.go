package grpc

import newsv1 "github.com/sanjevscet/news-grpc/api/news/v1"

// NewsServiceServer is the server API for NewsService service.

type Server struct {
	newsv1.UnimplementedNewsServiceServer
}

func NewServer() *Server {
	return &Server{}
}
