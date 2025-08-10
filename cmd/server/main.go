package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthV1 "google.golang.org/grpc/health/grpc_health_v1"

	newsv1 "github.com/sanjevscet/news-grpc/api/news/v1"

	inGrpc "github.com/sanjevscet/news-grpc/internal/grpc"
	"github.com/sanjevscet/news-grpc/internal/memstore"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	stores := memstore.NewStore()
	newsv1.RegisterNewsServiceServer(srv, inGrpc.NewServer(
		stores,
	))

	healthSrv := health.NewServer()

	healthV1.RegisterHealthServer(srv, healthSrv)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
