package part

import (
	"time"

	repoModel "github.com/mllbll/space-manufacture/inventory/internal/repository/model"
)

func (r *repository) InitTestData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// Двигатель
	r.data["550e8400-e29b-41d4-a716-446655440001"] = repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440001",
		Name:          "Ионный двигатель X-7",
		Description:   "Высокоэффективный ионный двигатель для дальних космических миссий",
		Price:         125000.50,
		Stock_quantity: 15,
		Category:      repoModel.CATEGORY_ENGINE,
		Dimensions: repoModel.Dimensions{
			Length: 250.0,
			Width:  120.0,
			Height: 180.0,
			Weight: 450.5,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "SpaceTech Industries",
			Country: "USA",
			Website: "https://spacetech.com",
		},
		Tags: []string{"engine", "ion", "premium", "long-range"},
		Metadata: map[string]repoModel.Value{
			"thrust":     repoModel.DoubleValue{Value: 2500.0},
			"efficiency": repoModel.StringValue{Value: "95%"},
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Топливо
	r.data["550e8400-e29b-41d4-a716-446655440002"] = repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440002",
		Name:          "Криогенное топливо H2-O2",
		Description:   "Высокоэнергетическое топливо для ракетных двигателей",
		Price:         8500.75,
		Stock_quantity: 500,
		Category:      repoModel.CATEGORY_FUEL,
		Dimensions: repoModel.Dimensions{
			Length: 100.0,
			Width:  100.0,
			Height: 200.0,
			Weight: 150.0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "FuelCorp",
			Country: "Russia",
			Website: "https://fuelcorp.ru",
		},
		Tags: []string{"fuel", "cryogenic", "high-energy"},
		Metadata: map[string]repoModel.Value{
			"energy_density": repoModel.DoubleValue{Value: 12.5},
			"temperature":     repoModel.Int64Value{Value: -253},
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Иллюминатор
	r.data["550e8400-e29b-41d4-a716-446655440003"] = repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440003",
		Name:          "Иллюминатор премиум класса",
		Description:   "Прочный иллюминатор с многослойным защитным стеклом",
		Price:         45000.00,
		Stock_quantity: 8,
		Category:      repoModel.CATEGORY_PORTHOLE,
		Dimensions: repoModel.Dimensions{
			Length: 80.0,
			Width:  80.0,
			Height: 15.0,
			Weight: 25.5,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "GlassWorks GmbH",
			Country: "Germany",
			Website: "https://glassworks.de",
		},
		Tags: []string{"window", "premium", "safety"},
		Metadata: map[string]repoModel.Value{
			"pressure_resistance": repoModel.DoubleValue{Value: 10.5},
			"certified":           repoModel.BoolValue{Value: true},
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Крыло
	r.data["550e8400-e29b-41d4-a716-446655440004"] = repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440004",
		Name:          "Аэродинамическое крыло Mark-III",
		Description:   "Легкое и прочное крыло для атмосферных полетов",
		Price:         32000.25,
		Stock_quantity: 12,
		Category:      repoModel.CATEGORY_WING,
		Dimensions: repoModel.Dimensions{
			Length: 500.0,
			Width:  200.0,
			Height: 50.0,
			Weight: 180.0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "AeroSpace Dynamics",
			Country: "USA",
			Website: "https://aerospace-dynamics.com",
		},
		Tags: []string{"wing", "aerodynamic", "lightweight"},
		Metadata: map[string]repoModel.Value{
			"lift_coefficient": repoModel.DoubleValue{Value: 1.8},
			"material":         repoModel.StringValue{Value: "Carbon Fiber"},
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Еще один двигатель
	r.data["550e8400-e29b-41d4-a716-446655440005"] = repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440005",
		Name:          "Плазменный двигатель P-42",
		Description:   "Мощный плазменный двигатель для тяжелых кораблей",
		Price:         185000.00,
		Stock_quantity: 5,
		Category:      repoModel.CATEGORY_ENGINE,
		Dimensions: repoModel.Dimensions{
			Length: 350.0,
			Width:  150.0,
			Height: 220.0,
			Weight: 680.0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "PlasmaTech Systems",
			Country: "Japan",
			Website: "https://plasmatech.jp",
		},
		Tags: []string{"engine", "plasma", "heavy-duty", "premium"},
		Metadata: map[string]repoModel.Value{
			"thrust":           repoModel.DoubleValue{Value: 5000.0},
			"power_consumption": repoModel.DoubleValue{Value: 150.0},
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Еще одно топливо
	r.data["550e8400-e29b-41d4-a716-446655440006"] = repoModel.Part{
		UUID:          "550e8400-e29b-41d4-a716-446655440006",
		Name:          "Антиматериальное топливо",
		Description:   "Экспериментальное топливо с максимальной энергоемкостью",
		Price:         250000.00,
		Stock_quantity: 2,
		Category:      repoModel.CATEGORY_FUEL,
		Dimensions: repoModel.Dimensions{
			Length: 50.0,
			Width:  50.0,
			Height: 100.0,
			Weight: 5.0,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "Quantum Fuel Labs",
			Country: "Switzerland",
			Website: "https://quantumfuel.ch",
		},
		Tags: []string{"fuel", "antimatter", "experimental", "premium"},
		Metadata: map[string]repoModel.Value{
			"energy_density": repoModel.DoubleValue{Value: 1000.0},
			"danger_level":   repoModel.StringValue{Value: "EXTREME"},
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}
