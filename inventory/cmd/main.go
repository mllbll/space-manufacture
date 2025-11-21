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
	"google.golang.org/protobuf/types/known/timestamppb"

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

		if len(filter.ManufacturerContries) > 0 {
			if part.Manufacturer == nil || !containsString(filter.ManufacturerContries, part.Manufacturer.Country) {
				continue
			}
		}

		if len(filter.Tags) > 0 && !hasOverlap(filter.Tags, part.Tags) {
			continue
		}
		result = append(result, part)
	}
	return result
}

func (s *inventoryService) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC in GetPart: %v\n", r)
			log.Printf("Request: %+v\n", req)
			panic(r)
		}
	}()

	if req == nil {
		log.Printf("GetPart: received nil request")
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	log.Printf("GetPart: received request with UUID: %q", req.Uuid)

	if req.Uuid == "" {
		log.Printf("GetPart: UUID is empty")
		return nil, status.Errorf(codes.InvalidArgument, "UUID cannot be empty")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	log.Printf("GetPart: looking up part with UUID: %q, total parts in map: %d", req.Uuid, len(s.parts))
	part, ok := s.parts[req.Uuid]
	if !ok {
		log.Printf("GetPart: part with UUID %q not found", req.Uuid)
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.Uuid)
	}

	if part == nil {
		log.Printf("GetPart: part with UUID %q is nil", req.Uuid)
		return nil, status.Errorf(codes.Internal, "part with UUID %s is nil", req.Uuid)
	}

	log.Printf("GetPart: found part with UUID %q, part details: UUID=%q, Name=%q", req.Uuid, part.Uuid, part.Name)

	// Создаем защитную копию, чтобы избежать проблем с конкурентным доступом
	partCopy := *part

	log.Printf("GetPart: creating response")
	resp := &inventoryV1.GetPartResponse{
		Parts: &partCopy,
	}
	log.Printf("GetPart: response created successfully, returning")
	return resp, nil
}

func (s *inventoryService) GetListParts(_ context.Context, req *inventoryV1.GetListPartsRequest) (*inventoryV1.GetListPartsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

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

// initTestData заполняет сервис тестовыми данными
func (s *inventoryService) initTestData() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := timestamppb.Now()

	// Двигатель
	s.parts["550e8400-e29b-41d4-a716-446655440001"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440001",
		Name:          "Ионный двигатель X-7",
		Description:   "Высокоэффективный ионный двигатель для дальних космических миссий",
		Price:         125000.50,
		StockQuantity: 15,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 250.0,
			Width:  120.0,
			Height: 180.0,
			Weight: 450.5,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "SpaceTech Industries",
			Country: "USA",
			Website: "https://spacetech.com",
		},
		Tags: []string{"engine", "ion", "premium", "long-range"},
		Metadata: map[string]*inventoryV1.Value{
			"thrust": {
				Types: &inventoryV1.Value_DoubleValue{DoubleValue: 2500.0},
			},
			"efficiency": {
				Types: &inventoryV1.Value_StringValue{StringValue: "95%"},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Топливо
	s.parts["550e8400-e29b-41d4-a716-446655440002"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440002",
		Name:          "Криогенное топливо H2-O2",
		Description:   "Высокоэнергетическое топливо для ракетных двигателей",
		Price:         8500.75,
		StockQuantity: 500,
		Category:      inventoryV1.Category_CATEGORY_FUEL,
		Dimensions: &inventoryV1.Dimensions{
			Length: 100.0,
			Width:  100.0,
			Height: 200.0,
			Weight: 150.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "FuelCorp",
			Country: "Russia",
			Website: "https://fuelcorp.ru",
		},
		Tags: []string{"fuel", "cryogenic", "high-energy"},
		Metadata: map[string]*inventoryV1.Value{
			"energy_density": {
				Types: &inventoryV1.Value_DoubleValue{DoubleValue: 12.5},
			},
			"temperature": {
				Types: &inventoryV1.Value_Int64Value{Int64Value: -253},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Иллюминатор
	s.parts["550e8400-e29b-41d4-a716-446655440003"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440003",
		Name:          "Иллюминатор премиум класса",
		Description:   "Прочный иллюминатор с многослойным защитным стеклом",
		Price:         45000.00,
		StockQuantity: 8,
		Category:      inventoryV1.Category_CATEGORY_PORTHOLE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 80.0,
			Width:  80.0,
			Height: 15.0,
			Weight: 25.5,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "GlassWorks GmbH",
			Country: "Germany",
			Website: "https://glassworks.de",
		},
		Tags: []string{"window", "premium", "safety"},
		Metadata: map[string]*inventoryV1.Value{
			"pressure_resistance": {
				Types: &inventoryV1.Value_DoubleValue{DoubleValue: 10.5},
			},
			"certified": {
				Types: &inventoryV1.Value_BoolValue{BoolValue: true},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Крыло
	s.parts["550e8400-e29b-41d4-a716-446655440004"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440004",
		Name:          "Аэродинамическое крыло Mark-III",
		Description:   "Легкое и прочное крыло для атмосферных полетов",
		Price:         32000.25,
		StockQuantity: 12,
		Category:      inventoryV1.Category_CATEGORY_WING,
		Dimensions: &inventoryV1.Dimensions{
			Length: 500.0,
			Width:  200.0,
			Height: 50.0,
			Weight: 180.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "AeroSpace Dynamics",
			Country: "USA",
			Website: "https://aerospace-dynamics.com",
		},
		Tags: []string{"wing", "aerodynamic", "lightweight"},
		Metadata: map[string]*inventoryV1.Value{
			"lift_coefficient": {
				Types: &inventoryV1.Value_DoubleValue{DoubleValue: 1.8},
			},
			"material": {
				Types: &inventoryV1.Value_StringValue{StringValue: "Carbon Fiber"},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Еще один двигатель
	s.parts["550e8400-e29b-41d4-a716-446655440005"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440005",
		Name:          "Плазменный двигатель P-42",
		Description:   "Мощный плазменный двигатель для тяжелых кораблей",
		Price:         185000.00,
		StockQuantity: 5,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 350.0,
			Width:  150.0,
			Height: 220.0,
			Weight: 680.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "PlasmaTech Systems",
			Country: "Japan",
			Website: "https://plasmatech.jp",
		},
		Tags: []string{"engine", "plasma", "heavy-duty", "premium"},
		Metadata: map[string]*inventoryV1.Value{
			"thrust": {
				Types: &inventoryV1.Value_DoubleValue{DoubleValue: 5000.0},
			},
			"power_consumption": {
				Types: &inventoryV1.Value_DoubleValue{DoubleValue: 150.0},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Еще одно топливо
	s.parts["550e8400-e29b-41d4-a716-446655440006"] = &inventoryV1.Part{
		Uuid:          "550e8400-e29b-41d4-a716-446655440006",
		Name:          "Антиматериальное топливо",
		Description:   "Экспериментальное топливо с максимальной энергоемкостью",
		Price:         250000.00,
		StockQuantity: 2,
		Category:      inventoryV1.Category_CATEGORY_FUEL,
		Dimensions: &inventoryV1.Dimensions{
			Length: 50.0,
			Width:  50.0,
			Height: 100.0,
			Weight: 5.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Quantum Fuel Labs",
			Country: "Switzerland",
			Website: "https://quantumfuel.ch",
		},
		Tags: []string{"fuel", "antimatter", "experimental", "premium"},
		Metadata: map[string]*inventoryV1.Value{
			"energy_density": {
				Types: &inventoryV1.Value_DoubleValue{DoubleValue: 1000.0},
			},
			"danger_level": {
				Types: &inventoryV1.Value_StringValue{StringValue: "EXTREME"},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	log.Printf("Initialized %d test parts", len(s.parts))
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

	// Инициализируем тестовые данные
	service.initTestData()

	inventoryV1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC recovered in gRPC server: %v\n", r)
				// Передаем панику дальше для полного стека
				panic(r)
			}
		}()
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
