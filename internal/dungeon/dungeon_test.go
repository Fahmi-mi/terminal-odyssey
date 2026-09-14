package dungeon

import (
	"testing"

	"github.com/Fahmi-mi/terminal-odyssey/internal/character"
)

func TestNewExpedition(t *testing.T) {
	p := character.NewDefaultPlayer("Hero")
	exp, err := NewExpedition(p, 3)
	if err != nil {
		t.Fatalf("failed to create expedition: %v", err)
	}

	if len(exp.Rooms) != 6 {
		t.Errorf("expected 6 rooms, got %d", len(exp.Rooms))
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
}

func TestExpeditionProgressionAndTorchDecay(t *testing.T) {
	p := character.NewDefaultPlayer("Explorer")
	exp, err := NewExpedition(p, 2)
	if err != nil {
		t.Fatalf("failed to create expedition: %v", err)
	}

	// Resolve combat in room 1 if present
	if exp.ActiveCombat != nil {
		exp.ActiveCombat.IsOver = true
		exp.OnCombatWon()
	}

	initialTorch := exp.Torch
	errAdvance := exp.AdvanceRoom()
	if errAdvance != nil {
		t.Fatalf("failed to advance room: %v", errAdvance)
	}

	if exp.CurrentRoomIdx != 1 {
		t.Errorf("expected current room idx 1, got %d", exp.CurrentRoomIdx)
	}
	if exp.Torch >= initialTorch {
		t.Errorf("expected torch to decay after advancing room")
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

	// Test Treasure room
	treasureRoom := exp.Rooms[1]
	exp.CurrentRoomIdx = 1

	if treasureRoom.Def.Type == RoomTypeTreasure {
		logText := exp.ResolveTreasure()
		if !treasureRoom.IsResolved {
			t.Errorf("expected treasure room to be marked resolved")
		}
		if exp.GoldFound <= 0 {
			t.Errorf("expected gold found to increase, got %d", exp.GoldFound)
		}
		if logText == "" {
			t.Errorf("expected non empty resolution log")
		}
	}

	// Test Rest room
	restRoom := exp.Rooms[3]
	exp.CurrentRoomIdx = 3
	exp.Torch = 50

	if restRoom.Def.Type == RoomTypeRest {
		exp.ResolveRest()
		if !restRoom.IsResolved {
			t.Errorf("expected rest room to be marked resolved")
		}
		if exp.Torch <= 50 {
			t.Errorf("expected torch to increase at sanctuary, got %d", exp.Torch)
		}
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
