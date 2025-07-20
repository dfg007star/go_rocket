package part

import (
	"crypto/rand"
	"math/big"
	"sync"
	"time"

	"github.com/google/uuid"

	def "github.com/dfg007star/go_rocket/inventory/internal/repository"
	repoModel "github.com/dfg007star/go_rocket/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data []repoModel.Part
}

func NewRepository() *repository {
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
		Length: float64(randomInt(10, 100)),
		Width:  float64(randomInt(5, 50)),
		Height: float64(randomInt(5, 50)),
		Weight: float64(randomInt(50, 1000)),
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
			Int64Value: int64Ptr(int64(randomInt(1, 11))),
		},
		"efficiency": {
			DoubleValue: float64Ptr(float64(randomInt(1, 10)) / 10),
		},
		"certified": {
			BoolValue: boolPtr(randomInt(0, 2) == 1),
		},
	}

	return repoModel.Part{
		Uuid:          uuid.New().String(),
		Name:          name,
		Description:   description,
		Price:         float64(randomInt(1000, 10000)),
		StockQuantity: int64(randomInt(1, 101)),
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
	return countries[randomInt(0, len(countries))]
}

func randomTag() string {
	tags := []string{"reliable", "durable", "efficient", "lightweight", "heavy-duty", "advanced", "next-gen"}
	return tags[randomInt(0, len(tags))]
}

func randomMaterial() string {
	materials := []string{"aluminum", "titanium", "carbon fiber", "steel", "composite", "ceramic"}
	return materials[randomInt(0, len(materials))]
}

func randomInt(min, max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		return min
	}
	return int(n.Int64()) + min
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
