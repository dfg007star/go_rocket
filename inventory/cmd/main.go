package main

import (
	"context"
	"fmt"
	inventoryV1 "github.com/dfg007star/go_rocket/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	// мокаем метку времени
	now := timestamppb.Now()

	// Регистрируем наш сервис
	service := &InventoryService{
		parts: map[string]*inventoryV1.Part{
			"123e4567-e89b-12d3-a456-426614174000": {
				Uuid:          "123e4567-e89b-12d3-a456-426614174000",
				Name:          "Turbo Engine V2",
				Description:   "High-performance aircraft engine",
				Price:         12599.99,
				StockQuantity: 15,
				Category:      inventoryV1.Category_CATEGORY_ENGINE,
				Dimensions: &inventoryV1.Dimensions{
					Length: 250.5,
					Width:  120.3,
					Height: 95.7,
					Weight: 425.8,
				},
				Manufacturer: &inventoryV1.Manufacturer{
					Name:    "AeroTech",
					Country: "Germany",
					Website: "https://aerotech.de",
				},
				Tags: []string{"engine", "turbo", "premium"},
				Metadata: map[string]*inventoryV1.Value{
					"warranty": {
						Value: &inventoryV1.Value_StringValue{
							StringValue: "5 years",
						},
					},
					"certified": {
						Value: &inventoryV1.Value_BoolValue{
							BoolValue: true,
						},
					},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			"550e8400-e29b-41d4-a716-446655440000": {
				Uuid:          "550e8400-e29b-41d4-a716-446655440000",
				Name:          "Aircraft Wing",
				Description:   "Composite material wing assembly",
				Price:         8750.50,
				StockQuantity: 8,
				Category:      inventoryV1.Category_CATEGORY_WING,
				Dimensions: &inventoryV1.Dimensions{
					Length: 600.0,
					Width:  280.0,
					Height: 45.2,
					Weight: 320.0,
				},
				Manufacturer: &inventoryV1.Manufacturer{
					Name:    "SkyComposites",
					Country: "USA",
					Website: "https://skycomposites.com",
				},
				Tags: []string{"wing", "composite", "assembly"},
				Metadata: map[string]*inventoryV1.Value{
					"material": {
						Value: &inventoryV1.Value_StringValue{
							StringValue: "carbon-fiber",
						},
					},
					"max_load": {
						Value: &inventoryV1.Value_DoubleValue{
							DoubleValue: 1250.75,
						},
					},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

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
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

// GetPart возвращает деталь
func (s *InventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, grpc.Errorf(codes.NotFound, "Part: not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *InventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parts := make([]*inventoryV1.Part, 0, len(s.parts))
	filter := req.Filter

	if filter == nil || isEmptyFilter(filter) {
		for _, part := range s.parts {
			parts = append(parts, part)
		}
		return &inventoryV1.ListPartsResponse{Parts: parts}, nil
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

	return &inventoryV1.ListPartsResponse{Parts: parts}, nil
}

func isMatchAnyFilter(part *inventoryV1.Part,
	uuidSet map[string]struct{},
	nameSet map[string]struct{},
	categorySet map[inventoryV1.Category]struct{},
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

func isEmptyFilter(filter *inventoryV1.PartsFilter) bool {
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

func makeCategorySet(categories []inventoryV1.Category) map[inventoryV1.Category]struct{} {
	set := make(map[inventoryV1.Category]struct{})
	for _, cat := range categories {
		set[cat] = struct{}{}
	}

	return set
}
