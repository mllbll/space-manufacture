package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	inventoryV1 "github.com/mllbll/space-manufacture/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50052

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

// ищет элемент в срезе строк
func containsString(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// ищет элемент типа Category в срезе таких значений
func containsCategory(list []inventoryV1.Category, cat inventoryV1.Category) bool {
	for _, v := range list {
		if v == cat {
			return true
		}
	}
	return false
}

// проверяет есть ли хотя бы одно совпадение в двух списках
func hasOverlap(a, b []string) bool {
	m := make(map[string]struct{}, len(a))
	for _, v := range a {
		m[v] = struct{}{}
	}
	for _, v := range b {
		if _, ok := m[v]; ok {
			return true
		}
	}
	return false
}

func ListParts(parts []inventoryV1.Part, filter inventoryV1.PartsFilter) []inventoryV1.Part {
	var result []inventoryV1.Part

	for _, part := range parts {
		if len(filter.Uuids) > 0 && !containsString(filter.Uuids, part.Uuid) {
			continue
		}

		if len(filter.Names) > 0 && !containsString(filter.Names, part.Name) {
			continue
		}

		if len(filter.Categories) > 0 && !containsCategory(filter.Categories, part.Category) {
			continue
		}

		if len(filter.ManufacturerContries) > 0 && !containsString(filter.ManufacturerContries, part.Manufacturer.Country) {
			continue
		}

		if len(filter.Tags) > 0 && !hasOverlap(filter.Tags, part.Tags) {
			continue
		}
		result = append(result, part)
	}
	return result
}

func (s *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.Unlock()

	part, ok := s.parts[req.Uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.Uuid)
	}

	return &inventoryV1.GetPartResponse{
		Parts: part,
	}, nil
}

func (s *inventoryService) GetListParts(_ context.Context, req *inventoryV1.GetListPartsRequest) (*inventoryV1.GetListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.Unlock()

	var parts []*inventoryV1.Part
	for _, part := range s.parts {
		p := part
		parts = append(parts, p)
	}

	// Применяем фильтрацию, если она указана
	if req.Filter != nil {
		var filteredParts []inventoryV1.Part
		for _, part := range parts {
			filteredParts = append(filteredParts, *part)
		}
		filtered := ListParts(filteredParts, *req.Filter)
		parts = make([]*inventoryV1.Part, len(filtered))
		for i := range filtered {
			parts[i] = &filtered[i]
		}
	}

	resp := &inventoryV1.GetListPartsResponse{
		Parts: parts,
	}

	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server ...")
	s.GracefulStop()
	log.Println("Server Stopped")

}
