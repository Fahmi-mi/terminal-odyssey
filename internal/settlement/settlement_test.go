package settlement

import (
	"testing"
)

func TestSettlementInitialState(t *testing.T) {
	s := NewSettlement("Oakhaven")

	if s.Name != "Oakhaven" {
		t.Fatalf("expected name Oakhaven, got %s", s.Name)
	}
	if s.Settlers != 6 {
		t.Fatalf("expected 6 settlers, got %d", s.Settlers)
	}
	if s.Workers.TotalAssigned() != 6 {
		t.Fatalf("expected 6 assigned workers, got %d", s.Workers.TotalAssigned())
	}
	if s.UnassignedSettlers() != 0 {
		t.Fatalf("expected 0 unassigned, got %d", s.UnassignedSettlers())
	}

	maxL, maxS, maxR := s.StorageCap()
	if maxL != 100 || maxS != 50 || maxR != 60 {
		t.Fatalf("unexpected storage caps: %d, %d, %d", maxL, maxS, maxR)
	}
}

func TestWorkerAssignment(t *testing.T) {
	s := NewSettlement("Oakhaven")
	s.Settlers = 7 // 1 unassigned

	if s.UnassignedSettlers() != 1 {
		t.Fatalf("expected 1 unassigned settler, got %d", s.UnassignedSettlers())
	}

	// Blacksmith fails if building level is 0
	err := s.AssignWorker(RoleBlacksmith)
	if err == nil {
		t.Fatalf("expected error assigning blacksmith without building, got nil")
	}

	// Assign to Farmer
	err = s.AssignWorker(RoleFarmer)
	if err != nil {
		t.Fatalf("failed to assign farmer: %v", err)
	}
	if s.Workers.Farmers != 3 {
		t.Fatalf("expected 3 farmers, got %d", s.Workers.Farmers)
	}

	// Now unassigned should be 0
	err = s.AssignWorker(RoleMiner)
	if err != ErrNoUnassignedSettlers {
		t.Fatalf("expected ErrNoUnassignedSettlers, got %v", err)
	}

	// Unassign Farmer
	err = s.UnassignWorker(RoleFarmer)
	if err != nil {
		t.Fatalf("failed to unassign farmer: %v", err)
	}
	if s.Workers.Farmers != 2 {
		t.Fatalf("expected 2 farmers, got %d", s.Workers.Farmers)
	}
}

func TestSimulateDay(t *testing.T) {
	s := NewSettlement("Oakhaven")
	s.Rations = 20

	result := s.SimulateDay(1, SeasonSpring)
	if result.FoodConsumed != 6 {
		t.Fatalf("expected 6 food consumed, got %d", result.FoodConsumed)
	}
	if len(result.Logs) == 0 {
		t.Fatalf("expected daily logs to be generated")
	}
}

func TestBuildingUpgrade(t *testing.T) {
	s := NewSettlement("Oakhaven")
	s.Lumber = 100
	s.Stone = 100
	s.Treasury = 300

	err := s.UpgradeBuilding(BuildingBlacksmith)
	if err != nil {
		t.Fatalf("failed to upgrade blacksmith: %v", err)
	}

	if s.Buildings[BuildingBlacksmith] != 1 {
		t.Fatalf("expected blacksmith level 1, got %d", s.Buildings[BuildingBlacksmith])
	}

	// Now assigning a blacksmith should succeed if there is an unassigned settler
	s.Settlers++
	err = s.AssignWorker(RoleBlacksmith)
	if err != nil {
		t.Fatalf("expected blacksmith assignment to succeed after building, got %v", err)
	}
}
