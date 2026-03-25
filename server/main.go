package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Chetas1/grpc-blog-service/config"
	"github.com/Chetas1/grpc-blog-service/internal/store"
	"github.com/Chetas1/grpc-blog-service/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runGRPCServer(cfg config.Config) {
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

func runGatewayServer(cfg config.Config) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	grpcServerEndpoint := fmt.Sprintf("localhost:%d", cfg.GrpcServer.Port)
	err := proto.RegisterBlogServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	fmt.Printf("Gateway Server is running on port %d...\n", cfg.GatewayServer.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GatewayServer.Port), mux); err != nil {
		log.Fatalf("failed to serve gateway: %v", err)
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	go runGRPCServer(cfg)

	runGatewayServer(cfg)
}
