package main

import (
	"context"
	"log"
	"time"

	"github.com/Chetas1/grpc-blog-service/config"
	"github.com/Chetas1/grpc-blog-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	conn, err := grpc.NewClient(cfg.GrpcClient.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewBlogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	runTest(ctx, client)
}
