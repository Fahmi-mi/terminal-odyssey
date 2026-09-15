package data

import "testing"

func TestLoadDefaultScenario(t *testing.T) {
	sc, err := LoadDefaultScenario()
	if err != nil {
		t.Fatalf("failed to load default scenario: %v", err)
	}

	if sc.ScenarioID != "default" {
		t.Fatalf("expected default scenario_id, got %s", sc.ScenarioID)
	}
	if sc.Village.Settlers != 6 {
		t.Fatalf("expected 6 settlers, got %d", sc.Village.Settlers)
	}
	if sc.Player.HP != 100 {
		t.Fatalf("expected 100 HP, got %d", sc.Player.HP)
	}
}

func TestLoadBuildingDefs(t *testing.T) {
	defs, err := LoadBuildingDefs()
	if err != nil {
		t.Fatalf("failed to load building defs: %v", err)
	}

	if len(defs) < 7 {
		t.Fatalf("expected at least 7 building defs, got %d", len(defs))
	}
}

func TestLoadAtmosphereEvents(t *testing.T) {
	events, err := LoadAtmosphereEvents()
	if err != nil {
		t.Fatalf("failed to load atmosphere events: %v", err)
	}

	if len(events) == 0 {
		t.Fatalf("expected non-empty atmosphere events")
	}
}

func TestLoadDungeonData(t *testing.T) {
	enemies, err := LoadEnemyDefs()
	if err != nil {
		t.Fatalf("failed to load enemy defs: %v", err)
	}
	if len(enemies) < 5 {
		t.Fatalf("expected at least 5 enemies, got %d", len(enemies))
	}

	rooms, err := LoadRoomDefs()
	if err != nil {
		t.Fatalf("failed to load room defs: %v", err)
	}
	if len(rooms) < 5 {
		t.Fatalf("expected at least 5 rooms, got %d", len(rooms))
	}

	items, err := LoadItemDefs()
	if err != nil {
		t.Fatalf("failed to load item defs: %v", err)
	}
	if len(items) < 3 {
		t.Fatalf("expected at least 3 items, got %d", len(items))
	}
}

func TestLoadRecipeDefs(t *testing.T) {
	recipes, err := LoadRecipeDefs()
	if err != nil {
		t.Fatalf("failed to load recipe defs: %v", err)
	}
	if len(recipes) < 4 {
		t.Fatalf("expected at least 4 recipe defs, got %d", len(recipes))
	}
}


