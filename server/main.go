// Command server starts both the gRPC server (configured port) and the
// gRPC-Gateway REST proxy (configured port), shares an in-memory BlogStore,
// and shuts both down gracefully on SIGINT/SIGTERM.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chetas1/grpc-blog-service/config"
	"github.com/Chetas1/grpc-blog-service/internal/store"
	"github.com/Chetas1/grpc-blog-service/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const shutdownTimeout = 10 * time.Second

func runGRPCServer(ctx context.Context, cfg *config.Config) error {
	lis, err := net.Listen(cfg.GrpcServer.Protocol, fmt.Sprintf(":%d", cfg.GrpcServer.Port))
	if err != nil {
		return fmt.Errorf("listen grpc: %w", err)
	}

	s := grpc.NewServer()
	proto.RegisterBlogServiceServer(s, &server{
		store: store.NewBlogStore(),
	})

	go func() {
		<-ctx.Done()
		log.Printf("gRPC server: graceful shutdown initiated")
		stopped := make(chan struct{})
		go func() {
			s.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
		case <-time.After(shutdownTimeout):
			log.Printf("gRPC server: graceful shutdown timed out, forcing stop")
			s.Stop()
		}
	}()

	log.Printf("gRPC Server listening on :%d", cfg.GrpcServer.Port)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("serve grpc: %w", err)
	}
	return nil
}

func runGatewayServer(ctx context.Context, cfg *config.Config) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	endpoint := fmt.Sprintf("localhost:%d", cfg.GrpcServer.Port)
	if err := proto.RegisterBlogServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return fmt.Errorf("register gateway handler: %w", err)
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.GatewayServer.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Shutdown goroutine intentionally uses a fresh context.Background so
	// srv.Shutdown gets the full timeout to drain in-flight requests after
	// the parent context has already been cancelled.
	go func() { // #nosec G118
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		log.Printf("gateway server: graceful shutdown initiated")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("gateway server: shutdown error: %v", err)
		}
	}()

	log.Printf("Gateway Server listening on :%d", cfg.GatewayServer.Port)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve gateway: %w", err)
	}
	return nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 2)

	go func() { errCh <- runGRPCServer(ctx, &cfg) }()
	go func() { errCh <- runGatewayServer(ctx, &cfg) }()

	select {
	case <-ctx.Done():
		log.Printf("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			log.Printf("server error: %v", err)
		}
		cancel()
	}

	// Drain the second goroutine so we exit cleanly.
	for i := 0; i < 1; i++ {
		if err := <-errCh; err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("server error during shutdown: %v", err)
		}
	}
}
