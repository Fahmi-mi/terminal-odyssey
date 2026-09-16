package settlement

import (
	"fmt"

	"github.com/Fahmi-mi/terminal-odyssey/data"
)

// Season represents the four distinct seasons of the year
type Season int

const (
	SeasonSpring Season = iota
	SeasonSummer
	SeasonAutumn
	SeasonWinter
)

func (s Season) String() string {
	switch s {
	case SeasonSpring:
		return "Musim Semi (Spring)"
	case SeasonSummer:
		return "Musim Panas (Summer)"
	case SeasonAutumn:
		return "Musim Gugur (Autumn)"
	case SeasonWinter:
		return "Musim Dingin (Winter)"
	default:
		return "Tidak Diketahui"
	}
}

// BuildingType identifiers
const (
	BuildingTownHall       = "Balai Desa"
	BuildingStorehouse     = "Gudang Logistik"
	BuildingBlacksmith     = "Bengkel Pandai Besi"
	BuildingTrainingGround = "Pusat Latihan"
	BuildingApothecary     = "Laboratorium Alkimia"
	BuildingTavern         = "Kedai Minum"
	BuildingCaravanPost    = "Pos Kafilah"
)

// BuildingInfo holds details and upgrade costs for a building
type BuildingInfo struct {
	Name        string
	Description string
	Level       int
	MaxLevel    int
	WoodCost    int
	StoneCost   int
	GoldCost    int
}

// Workers holds the allocation of population across roles
type Workers struct {
	Farmers     int
	Lumberjacks int
	Miners      int
	Blacksmiths int
	Militia     int
}

// TotalAssigned returns the count of all assigned workers
func (w Workers) TotalAssigned() int {
	return w.Farmers + w.Lumberjacks + w.Miners + w.Blacksmiths + w.Militia
}

// Settlement represents the village's resources, infrastructure, and inhabitants
type Settlement struct {
	Name       string
	Level      int
	Lumber     int
	Stone      int
	Rations    int
	Treasury   int // Gold
	Settlers   int // Total population
	DefenseVal int

	Workers     Workers
	Buildings   map[string]int // Building Name -> Level
	Commodities map[string]int // Commodity ID -> Quantity
}

// NewSettlementFromScenario creates a settlement from scenario data
func NewSettlementFromScenario(v data.ScenarioVillage) *Settlement {
	s := &Settlement{
		Name:       v.Name,
		Level:      v.Level,
		Lumber:     v.Lumber,
		Stone:      v.Stone,
		Rations:    v.Rations,
		Treasury:   v.Treasury,
		Settlers:   v.Settlers,
		DefenseVal: 20,
		Workers: Workers{
			Farmers:     v.Workers.Farmers,
			Lumberjacks: v.Workers.Lumberjacks,
			Miners:      v.Workers.Miners,
			Blacksmiths: v.Workers.Blacksmiths,
			Militia:     v.Workers.Militia,
		},
		Buildings:   make(map[string]int),
		Commodities: make(map[string]int),
	}

	for k, val := range v.Buildings {
		s.Buildings[k] = val
	}

	s.RecalculateDefense()
	return s
}

// NewSettlement creates an initial village state using the default scenario
func NewSettlement(name string) *Settlement {
	sc, err := data.LoadDefaultScenario()
	if err == nil {
		s := NewSettlementFromScenario(sc.Village)
		if name != "" {
			s.Name = name
		}
		return s
	}

	// Fallback in case of error
	s := &Settlement{
		Name:       name,
		Level:      1,
		Lumber:     60,
		Stone:      30,
		Rations:    40,
		Treasury:   200,
		Settlers:   6,
		DefenseVal: 20,
		Workers: Workers{
			Farmers:     2,
			Lumberjacks: 2,
			Miners:      1,
			Blacksmiths: 0,
			Militia:     1,
		},
		Buildings:   make(map[string]int),
		Commodities: make(map[string]int),
	}

	s.Buildings[BuildingTownHall] = 1
	s.Buildings[BuildingStorehouse] = 1
	s.RecalculateDefense()
	return s
}

// MaxSettlers returns maximum population supported by the Town Hall
func (s *Settlement) MaxSettlers() int {
	townHallLvl := s.Buildings[BuildingTownHall]
	if townHallLvl < 1 {
		townHallLvl = 1
	}
	return 5 + (townHallLvl * 5) // Lvl 1: 10, Lvl 2: 15, Lvl 3: 20, dll
}

// StorageCap returns maximum storage capacity for Lumber, Stone, and Rations based on Storehouse level
func (s *Settlement) StorageCap() (maxLumber, maxStone, maxRations int) {
	storeLvl := s.Buildings[BuildingStorehouse]
	if storeLvl < 1 {
		storeLvl = 1
	}
	maxLumber = 50 + (storeLvl * 50)  // Lvl 1: 100, Lvl 2: 150
	maxStone = 25 + (storeLvl * 25)   // Lvl 1: 50, Lvl 2: 75
	maxRations = 30 + (storeLvl * 30) // Lvl 1: 60, Lvl 2: 90
	return
}

// UnassignedSettlers returns how many settlers have not been assigned to a job
func (s *Settlement) UnassignedSettlers() int {
	unassigned := s.Settlers - s.Workers.TotalAssigned()
	if unassigned < 0 {
		return 0
	}
	return unassigned
}

// RecalculateDefense updates village defense value based on militia and buildings
func (s *Settlement) RecalculateDefense() {
	base := 10
	militiaBonus := s.Workers.Militia * 15
	townHallBonus := s.Buildings[BuildingTownHall] * 10
	s.DefenseVal = base + militiaBonus + townHallBonus
}

// GetCommodityStock returns current count of a commodity in village store
func (s *Settlement) GetCommodityStock(commodityID string) int {
	switch commodityID {
	case "lumber":
		return s.Lumber
	case "stone":
		return s.Stone
	case "rations":
		return s.Rations
	default:
		if s.Commodities == nil {
			s.Commodities = make(map[string]int)
		}
		return s.Commodities[commodityID]
	}
}

// AddCommodity adds stock of a commodity to village store
func (s *Settlement) AddCommodity(commodityID string, amount int) {
	switch commodityID {
	case "lumber":
		s.Lumber += amount
		maxL, _, _ := s.StorageCap()
		if s.Lumber > maxL {
			s.Lumber = maxL
		}
	case "stone":
		s.Stone += amount
		_, maxS, _ := s.StorageCap()
		if s.Stone > maxS {
			s.Stone = maxS
		}
	case "rations":
		s.Rations += amount
		_, _, maxR := s.StorageCap()
		if s.Rations > maxR {
			s.Rations = maxR
		}
	default:
		if s.Commodities == nil {
			s.Commodities = make(map[string]int)
		}
		s.Commodities[commodityID] += amount
	}
}

// DeductCommodity removes stock of a commodity from village store
func (s *Settlement) DeductCommodity(commodityID string, amount int) error {
	switch commodityID {
	case "lumber":
		if s.Lumber < amount {
			return fmt.Errorf("kayu tidak mencukupi")
		}
		s.Lumber -= amount
	case "stone":
		if s.Stone < amount {
			return fmt.Errorf("batu tidak mencukupi")
		}
		s.Stone -= amount
	case "rations":
		if s.Rations < amount {
			return fmt.Errorf("ransum tidak mencukupi")
		}
		s.Rations -= amount
	default:
		if s.Commodities == nil {
			s.Commodities = make(map[string]int)
		}
		if s.Commodities[commodityID] < amount {
			return fmt.Errorf("komoditas tidak mencukupi")
		}
		s.Commodities[commodityID] -= amount
	}
	return nil
}

