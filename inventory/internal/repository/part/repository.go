package part

import (
	"math/rand"
	"time"

	"github.com/google/uuid"

	def "github.com/dfg007star/go_rocket/inventory/internal/repository"
	repoModel "github.com/dfg007star/go_rocket/inventory/internal/repository/model"
	"sync"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data []repoModel.Part
}

func NewRepository() *repository {
	rand.Seed(time.Now().UnixNano())

	parts := []repoModel.Part{
		generateRandomPart("Engine XJ-2000", "High-performance rocket engine", repoModel.ENGINE),
		generateRandomPart("Fuel Tank T-500", "Large capacity fuel tank", repoModel.FUEL),
		generateRandomPart("Porthole P-100", "Reinforced glass porthole", repoModel.PORTHOLE),
	}

	return &repository{
		data: parts,
	}
}

// generateRandomPart for test purpose
func generateRandomPart(name, description string, category repoModel.Category) repoModel.Part {
	now := time.Now()

	dimensions := repoModel.Dimensions{
		Length: 10.0 + rand.Float64()*90.0,
		Width:  5.0 + rand.Float64()*45.0,
		Height: 5.0 + rand.Float64()*45.0,
		Weight: 50.0 + rand.Float64()*950.0,
	}

	manufacturer := repoModel.Manufacturer{
		Name:    "Rocket Parts Inc.",
		Country: randomCountry(),
		Website: "https://rocketparts.example.com",
	}

	tags := []string{
		"rocket",
		"space",
		name,
		randomTag(),
	}

	metadata := map[string]repoModel.Value{
		"material": {
			StringValue: stringPtr(randomMaterial()),
		},
		"durability": {
			Int64Value: int64Ptr(rand.Int63n(10) + 1),
		},
		"efficiency": {
			DoubleValue: float64Ptr(0.1 + rand.Float64()*0.9),
		},
		"certified": {
			BoolValue: boolPtr(rand.Intn(2) == 1),
		},
	}

	return repoModel.Part{
		Uuid:          uuid.New().String(),
		Name:          name,
		Description:   description,
		Price:         1000.0 + rand.Float64()*9000.0,
		StockQuantity: rand.Int63n(100) + 1,
		Category:      category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          tags,
		Metadata:      metadata,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func randomCountry() string {
	countries := []string{"USA", "Russia", "China", "Japan", "Germany", "France", "UK", "Canada"}
	return countries[rand.Intn(len(countries))]
}

func randomTag() string {
	tags := []string{"reliable", "durable", "efficient", "lightweight", "heavy-duty", "advanced", "next-gen"}
	return tags[rand.Intn(len(tags))]
}

func randomMaterial() string {
	materials := []string{"aluminum", "titanium", "carbon fiber", "steel", "composite", "ceramic"}
	return materials[rand.Intn(len(materials))]
}

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func boolPtr(b bool) *bool {
	return &b
}
