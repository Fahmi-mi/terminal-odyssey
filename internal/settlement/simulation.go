package settlement

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/Fahmi-mi/terminal-odyssey/data"
)

var atmosphereEvents []string

func init() {
	events, err := data.LoadAtmosphereEvents()
	if err == nil {
		atmosphereEvents = events
	}
}

// DailyResult summarizes what happened during a day transition
type DailyResult struct {
	DayPassed       int
	SeasonChanged   bool
	OldSeason       Season
	NewSeason       Season
	FoodProduced    int
	FoodConsumed    int
	LumberProduced  int
	LumberConsumed  int
	StoneProduced   int
	NewSettlers     int
	StarvationEvent bool
	Logs            []string
}

// SimulateDay runs 1 day of simulation for the settlement
func (s *Settlement) SimulateDay(day int, currentSeason Season) DailyResult {
	result := DailyResult{
		DayPassed: day,
		Logs:      make([]string, 0),
	}

	maxLumber, maxStone, maxRations := s.StorageCap()

	// 1. Food Production & Season Modifiers
	baseFoodPerFarmer := 4
	foodMod := 1.0
	switch currentSeason {
	case SeasonSpring:
		foodMod = 1.2 // +20%
	case SeasonSummer:
		foodMod = 1.35 // +35%
	case SeasonAutumn:
		foodMod = 1.0
	case SeasonWinter:
		foodMod = 0.0 // winter frozen fields
	}

	foodGained := int(float64(s.Workers.Farmers*baseFoodPerFarmer) * foodMod)
	result.FoodProduced = foodGained
	s.Rations += foodGained

	// 2. Lumber Production & Season Modifiers
	baseLumberPerWorker := 3
	lumberMod := 1.0
	if currentSeason == SeasonAutumn {
		lumberMod = 1.25 // Autumn +25% lumber
	}
	lumberGained := int(float64(s.Workers.Lumberjacks*baseLumberPerWorker) * lumberMod)
	result.LumberProduced = lumberGained
	s.Lumber += lumberGained

	// 3. Stone Production
	stoneGained := s.Workers.Miners * 2
	result.StoneProduced = stoneGained
	s.Stone += stoneGained

	// 4. Food Consumption (1 ration per settler)
	foodNeeded := s.Settlers
	result.FoodConsumed = foodNeeded
	if s.Rations >= foodNeeded {
		s.Rations -= foodNeeded
	} else {
		// Deficit / Kelaparan
		shortage := foodNeeded - s.Rations
		s.Rations = 0
		result.StarvationEvent = true
		result.Logs = append(result.Logs, fmt.Sprintf("[!] KELAPARAN: Defisit %d ransum! Warga menderita kekurangan pangan", shortage))
		// Risk of settler death if starvation is acute
		if shortage > 2 && s.Settlers > 3 && rand.Float64() < 0.4 {
			s.Settlers--
			result.Logs = append(result.Logs, "[-] Seorang warga gugur akibat penyakit dan kelaparan")
			s.RebalanceWorkers()
		}
	}

	// 5. Winter Wood Consumption for Heating
	if currentSeason == SeasonWinter {
		woodNeeded := 2 + (s.Settlers / 3)
		result.LumberConsumed = woodNeeded
		if s.Lumber >= woodNeeded {
			s.Lumber -= woodNeeded
			result.Logs = append(result.Logs, fmt.Sprintf("[-] Desa menghabiskan %d kayu bakar untuk bertahan dari badai es", woodNeeded))
		} else {
			s.Lumber = 0
			result.Logs = append(result.Logs, "[!] Kayu bakar habis! Warga menggigil kedinginan di tengah musim dingin")
		}
	}

	// Cap resources to Storehouse limits
	if s.Lumber > maxLumber {
		s.Lumber = maxLumber
	}
	if s.Stone > maxStone {
		s.Stone = maxStone
	}
	if s.Rations > maxRations {
		s.Rations = maxRations
	}

	// 6. Population Migration (Peluang warga baru datang bila ransum aman dan ada kuota balai desa)
	if s.Settlers < s.MaxSettlers() && s.Rations > 15 {
		migrantChance := 0.25
		if currentSeason == SeasonSpring {
			migrantChance = 0.50 // Musim semi kelahiran & migrasi naik
		}
		if rand.Float64() < migrantChance {
			s.Settlers++
			result.NewSettlers++
			result.Logs = append(result.Logs, "[+] Seorang pengembara terkesan dengan stabilitas desa dan memutuskan menetap")
		}
	}

	// 7. Consolidated Daily Production Log
	var prodParts []string
	if currentSeason == SeasonWinter && foodGained == 0 {
		prodParts = append(prodParts, "0 Ransum (Beku)")
	} else {
		prodParts = append(prodParts, fmt.Sprintf("+%d Ransum", foodGained))
	}
	prodParts = append(prodParts, fmt.Sprintf("+%d Kayu", lumberGained))
	prodParts = append(prodParts, fmt.Sprintf("+%d Batu", stoneGained))

	result.Logs = append(result.Logs, fmt.Sprintf("[+] Produksi Harian : %s", strings.Join(prodParts, " | ")))

	// 8. Random Atmosphere / Threat Warning Logs (hanya jika log belum melampaui kuota)
	s.RecalculateDefense()
	if rand.Float64() < 0.25 && len(atmosphereEvents) > 0 && len(result.Logs) < 3 {
		result.Logs = append(result.Logs, atmosphereEvents[rand.Intn(len(atmosphereEvents))])
	}

	return result
}

// RebalanceWorkers reduces assigned workers if population fell below assigned total
func (s *Settlement) RebalanceWorkers() {
	for s.Workers.TotalAssigned() > s.Settlers {
		if s.Workers.Militia > 0 {
			s.Workers.Militia--
		} else if s.Workers.Blacksmiths > 0 {
			s.Workers.Blacksmiths--
		} else if s.Workers.Miners > 0 {
			s.Workers.Miners--
		} else if s.Workers.Lumberjacks > 0 {
			s.Workers.Lumberjacks--
		} else if s.Workers.Farmers > 0 {
			s.Workers.Farmers--
		} else {
			break
		}
	}
	s.RecalculateDefense()
}
