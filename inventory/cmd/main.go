package main

import (
	"context"
	"fmt"
	partV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/part/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
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

	// Регистрируем наш сервис
	service := &InventoryService{
		parts: make(map[string]*partV1.Part),
	}

	partV1.RegisterPartServiceServer(s, service)

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

// InventoryService реализует gRPC сервис для работы с деталями корабля
type InventoryService struct {
	partV1.UnimplementedPartServiceServer

	mu    sync.RWMutex
	parts map[string]*partV1.Part
}

// GetPart возвращает деталь
func (s *InventoryService) GetPart(_ context.Context, req *partV1.GetPartRequest) (*partV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, grpc.Errorf(codes.NotFound, "Part: not found", req.GetUuid())
	}

	return &partV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *InventoryService) ListParts(_ context.Context, req *partV1.ListPartsRequest) (*partV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parts := make([]*partV1.Part, 0, len(s.parts))
	filter := req.Filter

	if filter == nil || isEmptyFilter(filter) {
		return &partV1.ListPartsResponse{Parts: parts}, nil
	}

	// создаем map для более быстрого поиска значения :)
	uuidSet := makeStringSet(filter.Uuids)
	nameSet := makeStringSet(filter.Names)
	categorySet := makeCategorySet(filter.Categories)
	countrySet := makeStringSet(filter.ManufacturerCountries)
	tagSet := makeStringSet(filter.Tags)

	for _, part := range s.parts {
		if isMatchAnyFilter(part, uuidSet, nameSet, categorySet, countrySet, tagSet) {
			parts = append(parts, part)
		}
	}

	return &partV1.ListPartsResponse{Parts: parts}, nil
}

func isMatchAnyFilter(part *partV1.Part,
	uuidSet map[string]struct{},
	nameSet map[string]struct{},
	categorySet map[partV1.Category]struct{},
	countrySet map[string]struct{},
	tagSet map[string]struct{}) bool {

	if len(uuidSet) > 0 {
		if _, exists := uuidSet[part.Uuid]; !exists {
			return false
		}
	}

	if len(nameSet) > 0 {
		if _, exists := nameSet[strings.ToLower(part.Name)]; !exists {
			return false
		}
	}

	if len(categorySet) > 0 {
		if _, exists := categorySet[part.Category]; !exists {
			return false
		}
	}

	if len(countrySet) > 0 {
		if _, exists := countrySet[strings.ToLower(part.Manufacturer.Country)]; !exists {
			return false
		}
	}

	if len(tagSet) > 0 {
		hasMatchingTag := false
		for _, tag := range part.Tags {
			if _, exists := tagSet[strings.ToLower(tag)]; exists {
				hasMatchingTag = true
				break
			}
		}
		if !hasMatchingTag {
			return false
		}
	}

	return true
}

func isEmptyFilter(filter *partV1.PartsFilter) bool {
	return len(filter.Uuids) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}

func makeStringSet(items []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, item := range items {
		set[strings.ToLower(item)] = struct{}{}
	}

	return set
}

func makeCategorySet(categories []partV1.Category) map[partV1.Category]struct{} {
	set := make(map[partV1.Category]struct{})
	for _, cat := range categories {
		set[cat] = struct{}{}
	}

	return set
}
