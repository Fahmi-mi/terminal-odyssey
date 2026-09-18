package siege

import (
	"fmt"
	"math/rand"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

// Status threat levels
const (
	ThreatStatusSafe     = "Aman"
	ThreatStatusAlert    = "Waspada"
	ThreatStatusDanger   = "Bahaya"
	ThreatStatusCritical = "Kritis"
)

// SiegeResult holds the outcome of a siege confrontation
type SiegeResult struct {
	Occurred       bool   `json:"occurred"`
	RaiderID       string `json:"raider_id"`
	RaiderName     string `json:"raider_name"`
	AssaultPower   int    `json:"assault_power"`
	DefensePower   int    `json:"defense_power"`
	Victory        bool   `json:"victory"`
	GoldDifference int    `json:"gold_difference"` // positive for looted reward, negative for stolen gold
	LumberLost     int    `json:"lumber_lost"`
	StoneLost      int    `json:"stone_lost"`
	RationsLost    int    `json:"rations_lost"`
	SettlersLost   int    `json:"settlers_lost"`
	SettlersGained int    `json:"settlers_gained"`
	Log            string `json:"log"`
}

// SiegeManager handles village siege calculations, threat progression, and combat resolution
type SiegeManager struct {
	Raiders             []data.RaiderGroupDef `json:"raiders"`
	LastSiegeDay        int                   `json:"last_siege_day"`
	SiegeCooldown       int                   `json:"siege_cooldown"`
	TotalSiegesRepelled int                   `json:"total_sieges_repelled"`
	TotalSiegesFailed   int                   `json:"total_sieges_failed"`
	LastResult          *SiegeResult          `json:"last_result,omitempty"`
}

// NewSiegeManager initializes a siege coordinator with raider definitions
func NewSiegeManager() (*SiegeManager, error) {
	defs, err := data.LoadRaiderDefs()
	if err != nil {
		return nil, fmt.Errorf("gagal memuat data penyerbu: %w", err)
	}

	return &SiegeManager{
		Raiders:       defs,
		LastSiegeDay:  0,
		SiegeCooldown: 3, // Minimum 3 days peace between sieges
	}, nil
}

// CalculateThreat determines current village threat score and descriptive status
func (sm *SiegeManager) CalculateThreat(v *settlement.Settlement, dayCounter int) (int, string) {
	if v == nil {
		return 0, ThreatStatusSafe
	}

	// 1. Wealth threat from treasury
	treasuryThreat := v.Treasury / 40

	// 2. Resource stockpile threat
	storageItems := v.Lumber + v.Stone + v.Rations
	for _, qty := range v.Commodities {
		storageItems += qty
	}
	resourceThreat := storageItems / 30

	// 3. Time elapsed threat
	timeThreat := dayCounter * 2

	// 4. Village administrative scale
	levelThreat := v.Level * 8

	totalScore := treasuryThreat + resourceThreat + timeThreat + levelThreat
	if totalScore < 0 {
		totalScore = 0
	}

	var status string
	switch {
	case totalScore < 30:
		status = ThreatStatusSafe
	case totalScore < 60:
		status = ThreatStatusAlert
	case totalScore < 90:
		status = ThreatStatusDanger
	default:
		status = ThreatStatusCritical
	}

	return totalScore, status
}

// ShouldTriggerSiege evaluates if a siege raid occurs on the current day
func (sm *SiegeManager) ShouldTriggerSiege(threat int, dayCounter int) bool {
	// Respect cooldown period
	if dayCounter-sm.LastSiegeDay < sm.SiegeCooldown {
		return false
	}

	// Threat below 25 does not trigger sieges
	if threat < 25 {
		return false
	}

	// Calculate raid probability based on threat bracket
	var chance int
	switch {
	case threat < 60:
		chance = 20
	case threat < 90:
		chance = 40
	default:
		chance = 65
	}

	roll := rand.Intn(100)
	return roll < chance
}

// SelectRaider chooses an appropriate raider group matching the village threat score
func (sm *SiegeManager) SelectRaider(threat int) *data.RaiderGroupDef {
	if len(sm.Raiders) == 0 {
		return nil
	}

	var candidates []data.RaiderGroupDef
	for _, r := range sm.Raiders {
		if threat >= r.MinThreat {
			candidates = append(candidates, r)
		}
	}

	if len(candidates) == 0 {
		return &sm.Raiders[0]
	}

	// Pick randomly among qualifying raider forces
	chosen := candidates[rand.Intn(len(candidates))]
	return &chosen
}

// ResolveSiege calculates battle between the attacking force and village defenses
func (sm *SiegeManager) ResolveSiege(raider data.RaiderGroupDef, v *settlement.Settlement, dayCounter int) *SiegeResult {
	sm.LastSiegeDay = dayCounter

	variance := raider.AssaultVariance
	if variance < 1 {
		variance = 5
	}
	delta := rand.Intn(2*variance+1) - variance
	assaultPower := raider.BaseAssaultPower + delta
	if assaultPower < 10 {
		assaultPower = 10
	}

	// Defense power based on settlement defensive fortifications and militia
	v.RecalculateDefense()
	defensePower := v.DefenseVal

	res := &SiegeResult{
		Occurred:     true,
		RaiderID:     raider.ID,
		RaiderName:   raider.Name,
		AssaultPower: assaultPower,
		DefensePower: defensePower,
	}

	if defensePower >= assaultPower {
		// Village victory - repel attackers
		res.Victory = true
		sm.TotalSiegesRepelled++

		// Reward looted gold
		rewardGold := raider.GoldLootReward
		res.GoldDifference = rewardGold
		v.Treasury += rewardGold

		// Potential captive recruit
		if rand.Intn(100) < raider.CaptiveChance && v.Settlers < v.MaxSettlers() {
			res.SettlersGained = 1
			v.Settlers++
			res.Log = fmt.Sprintf("[!] PENGEPUNGAN: %s menyerbu gerbang (Kekuatan %d vs Pertahanan %d). Pasukan desa berhasil memukul mundur penyerbu! Merampas +%d Emas dan membebaskan 1 tawanan warga baru", raider.Name, assaultPower, defensePower, rewardGold)
		} else {
			res.Log = fmt.Sprintf("[!] PENGEPUNGAN: %s menyerbu gerbang (Kekuatan %d vs Pertahanan %d). Pasukan desa berhasil memukul mundur musuh dan merampas +%d Emas", raider.Name, assaultPower, defensePower, rewardGold)
		}
	} else {
		// Village defeat - wall breached
		res.Victory = false
		sm.TotalSiegesFailed++

		// Gold pillaged: 25% of treasury
		stolenGold := (v.Treasury * 25) / 100
		if stolenGold > 150 {
			stolenGold = 150
		}
		if stolenGold < 15 && v.Treasury > 0 {
			stolenGold = v.Treasury
		}
		v.Treasury -= stolenGold
		res.GoldDifference = -stolenGold

		// Rations looted: 20%
		stolenRations := (v.Rations * 20) / 100
		v.Rations -= stolenRations
		res.RationsLost = stolenRations

		// Lumber looted: 15%
		stolenLumber := (v.Lumber * 15) / 100
		v.Lumber -= stolenLumber
		res.LumberLost = stolenLumber

		// Stone looted: 10%
		stolenStone := (v.Stone * 10) / 100
		v.Stone -= stolenStone
		res.StoneLost = stolenStone

		// Population loss chance if settlement has sufficient inhabitants
		if v.Settlers > 4 && rand.Intn(100) < 35 {
			res.SettlersLost = 1
			v.Settlers--
			// Rebalance unassigned workers if needed
			if v.Workers.TotalAssigned() > v.Settlers {
				if v.Workers.Farmers > 0 {
					v.Workers.Farmers--
				} else if v.Workers.Lumberjacks > 0 {
					v.Workers.Lumberjacks--
				} else if v.Workers.Miners > 0 {
					v.Workers.Miners--
				} else if v.Workers.Militia > 0 {
					v.Workers.Militia--
				}
			}
			res.Log = fmt.Sprintf("[!] PENGEPUNGAN: %s menerobos gerbang desa (Kekuatan %d vs Pertahanan %d). Pemukiman dijarah: -%d Emas, -%d Ransum, dan 1 warga gugur dalam serbuan", raider.Name, assaultPower, defensePower, stolenGold, stolenRations)
		} else {
			res.Log = fmt.Sprintf("[!] PENGEPUNGAN: %s menerobos gerbang desa (Kekuatan %d vs Pertahanan %d). Pemukiman dijarah: -%d Emas, -%d Ransum, -%d Kayu", raider.Name, assaultPower, defensePower, stolenGold, stolenRations, stolenLumber)
		}
	}

	sm.LastResult = res
	return res
}
