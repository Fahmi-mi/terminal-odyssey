package dungeon

import (
	"testing"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/character"
)

func TestNewExpedition(t *testing.T) {
	p := character.NewDefaultPlayer("Hero")
	exp, err := NewExpedition(p, 3)
	if err != nil {
		t.Fatalf("failed to create expedition: %v", err)
	}

	if len(exp.Rooms) != 9 {
		t.Errorf("expected 9 rooms in DAG, got %d", len(exp.Rooms))
	}
	if exp.TotalDepths != 6 {
		t.Errorf("expected 6 total depths, got %d", exp.TotalDepths)
	}
	if exp.RoomsExploredCount != 1 {
		t.Errorf("expected 1 room explored at start, got %d", exp.RoomsExploredCount)
	}
	if exp.Torch != 100 {
		t.Errorf("expected 100 torch, got %d", exp.Torch)
	}
	if exp.Rations != 3 {
		t.Errorf("expected 3 rations, got %d", exp.Rations)
	}
	if exp.CurrentRoomIdx != 0 {
		t.Errorf("expected current room idx 0, got %d", exp.CurrentRoomIdx)
	}

	// Verify DAG topology
	// Node 0 connects to 1 and 2
	if len(exp.Rooms[0].NextRoomIndices) != 2 || exp.Rooms[0].NextRoomIndices[0] != 1 || exp.Rooms[0].NextRoomIndices[1] != 2 {
		t.Errorf("expected node 0 to connect to [1, 2], got %v", exp.Rooms[0].NextRoomIndices)
	}
	// Node 1 & 2 connect to 3 and 4
	if len(exp.Rooms[1].NextRoomIndices) != 2 || exp.Rooms[1].NextRoomIndices[0] != 3 || exp.Rooms[1].NextRoomIndices[1] != 4 {
		t.Errorf("expected node 1 to connect to [3, 4], got %v", exp.Rooms[1].NextRoomIndices)
	}
	if len(exp.Rooms[2].NextRoomIndices) != 2 || exp.Rooms[2].NextRoomIndices[0] != 3 || exp.Rooms[2].NextRoomIndices[1] != 4 {
		t.Errorf("expected node 2 to connect to [3, 4], got %v", exp.Rooms[2].NextRoomIndices)
	}
	// Node 3 & 4 connect to 5 and 6
	if len(exp.Rooms[3].NextRoomIndices) != 2 || exp.Rooms[3].NextRoomIndices[0] != 5 || exp.Rooms[3].NextRoomIndices[1] != 6 {
		t.Errorf("expected node 3 to connect to [5, 6], got %v", exp.Rooms[3].NextRoomIndices)
	}
	if len(exp.Rooms[4].NextRoomIndices) != 2 || exp.Rooms[4].NextRoomIndices[0] != 5 || exp.Rooms[4].NextRoomIndices[1] != 6 {
		t.Errorf("expected node 4 to connect to [5, 6], got %v", exp.Rooms[4].NextRoomIndices)
	}
	// Node 5 & 6 connect to 7 (Boss)
	if len(exp.Rooms[5].NextRoomIndices) != 1 || exp.Rooms[5].NextRoomIndices[0] != 7 {
		t.Errorf("expected node 5 to connect to [7], got %v", exp.Rooms[5].NextRoomIndices)
	}
	if len(exp.Rooms[6].NextRoomIndices) != 1 || exp.Rooms[6].NextRoomIndices[0] != 7 {
		t.Errorf("expected node 6 to connect to [7], got %v", exp.Rooms[6].NextRoomIndices)
	}
	// Node 7 connects to 8 (Exit)
	if len(exp.Rooms[7].NextRoomIndices) != 1 || exp.Rooms[7].NextRoomIndices[0] != 8 {
		t.Errorf("expected node 7 to connect to [8], got %v", exp.Rooms[7].NextRoomIndices)
	}
	// Node 8 has 0 outgoing
	if len(exp.Rooms[8].NextRoomIndices) != 0 {
		t.Errorf("expected node 8 to have 0 next rooms, got %v", exp.Rooms[8].NextRoomIndices)
	}

	// Verify depths
	expectedDepths := []int{1, 2, 2, 3, 3, 4, 4, 5, 6}
	for i, r := range exp.Rooms {
		if r.Depth != expectedDepths[i] {
			t.Errorf("expected room %d to have depth %d, got %d", i, expectedDepths[i], r.Depth)
		}
	}
}

func TestProceduralLifelineAndPacingRules(t *testing.T) {
	p := character.NewDefaultPlayer("Pathfinder")

	for run := 0; run < 100; run++ {
		exp, err := NewExpedition(p, 2)
		if err != nil {
			t.Fatalf("run %d: failed to create expedition: %v", run, err)
		}

		node1 := exp.Rooms[1]
		node2 := exp.Rooms[2]
		node3 := exp.Rooms[3]
		node4 := exp.Rooms[4]
		node5 := exp.Rooms[5]
		node6 := exp.Rooms[6]

		depth2HasRest := (node1.Def.Type == RoomTypeRest || node2.Def.Type == RoomTypeRest)
		depth3HasRest := (node3.Def.Type == RoomTypeRest || node4.Def.Type == RoomTypeRest)
		depth4HasRest := (node5.Def.Type == RoomTypeRest || node6.Def.Type == RoomTypeRest)

		// Rule 1: Lifeline guarantee - if no rest at Depth 2, Depth 3 MUST have at least one rest
		if !depth2HasRest && !depth3HasRest {
			t.Fatalf("run %d: lifeline failed - neither depth 2 nor depth 3 had a rest sanctuary", run)
		}

		// Rule 2: No consecutive rest rooms along any possible path
		if depth2HasRest && depth3HasRest {
			t.Fatalf("run %d: pacing violation - depth 2 and depth 3 both generated rest rooms", run)
		}
		if depth3HasRest && depth4HasRest {
			t.Fatalf("run %d: pacing violation - depth 3 and depth 4 both generated rest rooms", run)
		}

		// Rule 3: Total rest sanctuaries across all 9 nodes must not exceed 2
		totalRest := 0
		for _, r := range exp.Rooms {
			if r.Def.Type == RoomTypeRest {
				totalRest++
			}
		}
		if totalRest > 2 {
			t.Fatalf("run %d: expected at most 2 rest sanctuaries, got %d", run, totalRest)
		}
	}
}

func TestExpeditionProgressionAndTorchDecay(t *testing.T) {
	p := character.NewDefaultPlayer("Explorer")
	exp, err := NewExpedition(p, 2)
	if err != nil {
		t.Fatalf("failed to create expedition: %v", err)
	}

	// Resolve combat in room 1 if present, or mark resolved if non-combat
	if exp.ActiveCombat != nil {
		exp.ActiveCombat.IsOver = true
		exp.OnCombatWon()
	} else {
		exp.Rooms[0].IsResolved = true
	}

	initialTorch := exp.Torch
	errAdvance := exp.AdvanceRoom()
	if errAdvance != nil {
		t.Fatalf("failed to advance room: %v", errAdvance)
	}

	if exp.CurrentRoomIdx != 1 {
		t.Errorf("expected current room idx 1, got %d", exp.CurrentRoomIdx)
	}
	if exp.RoomsExploredCount != 2 {
		t.Errorf("expected 2 rooms explored, got %d", exp.RoomsExploredCount)
	}
	if exp.Torch >= initialTorch {
		t.Errorf("expected torch to decay after advancing room")
	}
}

func TestExpeditionBranchingChoicesAndValidation(t *testing.T) {
	p := character.NewDefaultPlayer("Tactician")
	exp, _ := NewExpedition(p, 2)

	// Test combat lock
	if exp.ActiveCombat != nil {
		errInCombat := exp.AdvanceToRoom(1)
		if errInCombat == nil {
			t.Errorf("expected error advancing while combat is active")
		}
		exp.ActiveCombat.IsOver = true
		exp.OnCombatWon()
	} else {
		exp.Rooms[0].IsResolved = true
	}
	choices := exp.NextRoomChoices()
	if len(choices) != 2 {
		t.Fatalf("expected 2 next room choices, got %d", len(choices))
	}
	if choices[0].GraphIdx != 1 || choices[1].GraphIdx != 2 {
		t.Errorf("expected choices with indices 1 and 2, got %d and %d", choices[0].GraphIdx, choices[1].GraphIdx)
	}

	// Error on invalid disconnected room
	errInvalid := exp.AdvanceToRoom(5)
	if errInvalid == nil {
		t.Errorf("expected error advancing to unconnected room 5 from room 0")
	}

	// Successfully choose branch 2 (Lorong Kanan)
	errBranch := exp.AdvanceToRoom(2)
	if errBranch != nil {
		t.Fatalf("unexpected error choosing branch 2: %v", errBranch)
	}
	if exp.CurrentRoomIdx != 2 {
		t.Errorf("expected current room idx 2, got %d", exp.CurrentRoomIdx)
	}
}

func TestExpeditionConsumeRation(t *testing.T) {
	p := character.NewDefaultPlayer("HungryHero")
	p.HP = 40
	p.MaxHP = 100

	exp, _ := NewExpedition(p, 2)
	healed, err := exp.ConsumeRation()
	if err != nil {
		t.Fatalf("failed to consume ration: %v", err)
	}
	if healed != 25 {
		t.Errorf("expected 25 healed HP, got %d", healed)
	}
	if exp.Rations != 1 {
		t.Errorf("expected 1 ration remaining, got %d", exp.Rations)
	}
	if p.HP != 65 {
		t.Errorf("expected 65 HP, got %d", p.HP)
	}
}

func TestExpeditionResolveTreasureAndRest(t *testing.T) {
	p := character.NewDefaultPlayer("TreasureHunter")
	p.HP = 50
	exp, _ := NewExpedition(p, 0)

	// Setup a dedicated treasure room to test resolution
	exp.Rooms[1].Def = data.RoomDef{
		Type:       RoomTypeTreasure,
		Title:      "Ruang Simpanan Penjarah",
		MinGold:    30,
		MaxGold:    50,
		MinRations: 1,
		MaxRations: 2,
		Materials: []data.RoomMaterial{
			{Type: "lumber", Amount: 10},
			{Type: "stone", Amount: 5},
		},
	}
	exp.Rooms[1].IsResolved = false
	exp.CurrentRoomIdx = 1

	logText := exp.ResolveTreasure()
	if !exp.Rooms[1].IsResolved {
		t.Errorf("expected treasure room to be marked resolved")
	}
	if exp.GoldFound < 30 {
		t.Errorf("expected gold found to be at least 30, got %d", exp.GoldFound)
	}
	if exp.LumberFound != 10 || exp.StoneFound != 5 {
		t.Errorf("expected materials gained (10 lumber, 5 stone), got %d lumber, %d stone", exp.LumberFound, exp.StoneFound)
	}
	if logText == "" {
		t.Errorf("expected non empty resolution log")
	}

	// Setup a dedicated rest room to test resolution
	exp.Rooms[3].Def = data.RoomDef{
		Type:        RoomTypeRest,
		Title:       "Suaka Pertapa Tersembunyi",
		HealPercent: 35,
		TorchBonus:  30,
	}
	exp.Rooms[3].IsResolved = false
	exp.CurrentRoomIdx = 3
	exp.Torch = 50

	exp.ResolveRest()
	if !exp.Rooms[3].IsResolved {
		t.Errorf("expected rest room to be marked resolved")
	}
	if exp.Torch != 80 {
		t.Errorf("expected torch to increase from 50 to 80, got %d", exp.Torch)
	}
}

func TestExpeditionEvacuateAndDefeat(t *testing.T) {
	p := character.NewDefaultPlayer("Survivor")
	exp, _ := NewExpedition(p, 1)

	exp.GoldFound = 100
	exp.Evacuate()
	if !exp.IsCompleted {
		t.Errorf("expected expedition to be completed")
	}
	if exp.GoldFound != 100 {
		t.Errorf("expected gold to be preserved on evacuation")
	}

	// Test defeat strips loot
	exp2, _ := NewExpedition(p, 1)
	exp2.GoldFound = 150
	exp2.HandleDefeat()

	if !exp2.IsDefeated || !exp2.IsCompleted {
		t.Errorf("expected expedition to be defeated and completed")
	}
	if exp2.GoldFound != 0 {
		t.Errorf("expected gold to be stripped on defeat, got %d", exp2.GoldFound)
	}
}

func TestEntranceRoomVariability(t *testing.T) {
	p := character.NewDefaultPlayer("Scout")
	combatCount := 0
	nonCombatCount := 0

	for i := 0; i < 300; i++ {
		exp, err := NewExpedition(p, 2)
		if err != nil {
			t.Fatalf("run %d: failed to create expedition: %v", i, err)
		}
		r0 := exp.Rooms[0]
		if r0.Def.Type == RoomTypeRest {
			t.Fatalf("run %d: entrance room must never be a rest sanctuary", i)
		}
		if r0.Def.Type == RoomTypeCombat {
			combatCount++
		} else {
			nonCombatCount++
		}
	}

	if combatCount <= nonCombatCount {
		t.Errorf("expected combat to be overwhelming majority, got %d combat vs %d non-combat", combatCount, nonCombatCount)
	}
	if nonCombatCount == 0 {
		t.Errorf("expected at least some non-combat entrance rooms across 300 runs, got %d", nonCombatCount)
	}
}

