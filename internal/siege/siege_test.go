package siege

import (
	"testing"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

func TestNewSiegeManager(t *testing.T) {
	sm, err := NewSiegeManager()
	if err != nil {
		t.Fatalf("failed to create siege manager: %v", err)
	}

	if len(sm.Raiders) == 0 {
		t.Errorf("expected raiders to be populated")
	}
	if sm.SiegeCooldown != 3 {
		t.Errorf("expected siege cooldown 3, got %d", sm.SiegeCooldown)
	}
}

func TestCalculateThreat(t *testing.T) {
	sm, err := NewSiegeManager()
	if err != nil {
		t.Fatalf("failed to create siege manager: %v", err)
	}

	v := settlement.NewSettlement("Desa Damai")
	v.Treasury = 100
	v.Lumber = 20
	v.Stone = 10
	v.Rations = 15
	v.Level = 1

	score, status := sm.CalculateThreat(v, 1)
	if score <= 0 {
		t.Errorf("expected positive threat score, got %d", score)
	}
	if status != ThreatStatusSafe {
		t.Errorf("expected status %s for low threat, got %s", ThreatStatusSafe, status)
	}

	// High threat setup
	v.Treasury = 2000
	v.Lumber = 300
	v.Stone = 200
	v.Rations = 300
	v.Level = 4

	highScore, highStatus := sm.CalculateThreat(v, 40)
	if highScore < 90 {
		t.Errorf("expected threat >= 90, got %d", highScore)
	}
	if highStatus != ThreatStatusCritical {
		t.Errorf("expected %s status, got %s", ThreatStatusCritical, highStatus)
	}
}

func TestShouldTriggerSiegeCooldown(t *testing.T) {
	sm, _ := NewSiegeManager()
	sm.LastSiegeDay = 5

	// Day 6 is within cooldown (diff = 1 < 3)
	if sm.ShouldTriggerSiege(100, 6) {
		t.Errorf("expected siege to NOT trigger within cooldown")
	}

	// Threat below 25 should never trigger
	if sm.ShouldTriggerSiege(10, 15) {
		t.Errorf("expected low threat to never trigger siege")
	}
}

func TestResolveSiegeVictory(t *testing.T) {
	sm, _ := NewSiegeManager()
	v := settlement.NewSettlement("Benteng Kuat")
	v.DefenseVal = 300 // overwhelming defense
	v.Treasury = 100
	initialTreasury := v.Treasury

	raider := data.RaiderGroupDef{
		ID:               "test_raider",
		Name:             "Perampok Uji",
		BaseAssaultPower: 20,
		AssaultVariance:  2,
		GoldLootReward:   50,
		CaptiveChance:    100,
	}

	res := sm.ResolveSiege(raider, v, 10)
	if !res.Victory {
		t.Fatalf("expected victory with overwhelming defense")
	}
	if sm.TotalSiegesRepelled != 1 {
		t.Errorf("expected 1 siege repelled, got %d", sm.TotalSiegesRepelled)
	}
	if v.Treasury != initialTreasury+50 {
		t.Errorf("expected treasury %d, got %d", initialTreasury+50, v.Treasury)
	}
	if res.Log == "" {
		t.Errorf("expected log string to be populated")
	}
}

func TestResolveSiegeDefeat(t *testing.T) {
	sm, _ := NewSiegeManager()
	v := settlement.NewSettlement("Desa Rentan")
	v.DefenseVal = 5 // very low defense
	v.Treasury = 200
	v.Rations = 100
	v.Lumber = 100
	v.Stone = 100
	v.Settlers = 8

	raider := data.RaiderGroupDef{
		ID:               "strong_raider",
		Name:             "Penyerbu Perkasa",
		BaseAssaultPower: 150,
		AssaultVariance:  5,
	}

	res := sm.ResolveSiege(raider, v, 12)
	if res.Victory {
		t.Fatalf("expected defeat with low defense")
	}
	if sm.TotalSiegesFailed != 1 {
		t.Errorf("expected 1 siege failed, got %d", sm.TotalSiegesFailed)
	}
	if v.Treasury >= 200 {
		t.Errorf("expected gold loss after defeat, treasury is %d", v.Treasury)
	}
	if res.GoldDifference >= 0 {
		t.Errorf("expected negative gold difference, got %d", res.GoldDifference)
	}
}
