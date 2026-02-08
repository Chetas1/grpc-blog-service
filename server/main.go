package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Chetas1/grpc-blog-service/config"
	"github.com/Chetas1/grpc-blog-service/internal/store"
	"github.com/Chetas1/grpc-blog-service/proto"
	"google.golang.org/grpc"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	lis, err := net.Listen(cfg.GrpcServer.Protocol, fmt.Sprintf(":%d", cfg.GrpcServer.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	proto.RegisterBlogServiceServer(s, &server{
		store: store.NewBlogStore(),
	})

	fmt.Printf("gRPC Server is running on port %d...\n", cfg.GrpcServer.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
