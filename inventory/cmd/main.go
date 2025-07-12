package main

import (
	"fmt"
	inventoryAPI "github.com/dfg007star/go_rocket/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/dfg007star/go_rocket/inventory/internal/repository/part"
	inventoryService "github.com/dfg007star/go_rocket/inventory/internal/service/part"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))

	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if err := lis.Close(); err != nil {
			log.Printf("failed to close listener: %v\n", err)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	repo := inventoryRepository.NewRepository()
	service := inventoryService.NewService(repo)
	api := inventoryAPI.NewApi(service)
	inventoryV1.RegisterInventoryServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
