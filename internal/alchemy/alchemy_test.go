package alchemy

import (
	"testing"
)

func TestAlchemyManager_CanBrewAndGating(t *testing.T) {
	mgr, err := NewAlchemyManager()
	if err != nil {
		t.Fatalf("failed to initialize alchemy manager: %v", err)
	}

	commodities := map[string]int{
		"herbal_salve": 2,
		"spices":       0,
	}

	// Case 1: Level too low for Level 2 recipe (antidote)
	err = mgr.CanBrew("penawar_racun", 1, 50, commodities, 10, 10, 10)
	if err == nil {
		t.Errorf("expected error for level 1 laboratory on level 2 recipe")
	}

	// Case 2: Insufficient gold
	err = mgr.CanBrew("salep_pemulih", 1, 5, commodities, 10, 10, 10)
	if err == nil {
		t.Errorf("expected error for insufficient gold")
	}

	// Case 3: Insufficient ingredient
	err = mgr.CanBrew("tonik_penenang", 2, 50, commodities, 10, 10, 10)
	if err == nil {
		t.Errorf("expected error for missing spices ingredient")
	}

	// Case 4: Valid brew for Level 1 recipe
	err = mgr.CanBrew("salep_pemulih", 1, 50, commodities, 10, 10, 10)
	if err != nil {
		t.Errorf("expected valid brew, got error: %v", err)
	}
}

func TestAlchemyManager_BrewExecution(t *testing.T) {
	mgr, err := NewAlchemyManager()
	if err != nil {
		t.Fatalf("failed to initialize alchemy manager: %v", err)
	}

	commodities := map[string]int{
		"herbal_salve": 5,
	}

	// High ingenuity test
	res, err := mgr.Brew("salep_pemulih", 1, 50, commodities, 10, 10, 10, 18)
	if err != nil {
		t.Fatalf("failed to brew potion: %v", err)
	}
	if res.YieldAmount < 1 || res.YieldAmount > 2 {
		t.Errorf("expected yield amount 1 or 2, got %d", res.YieldAmount)
	}
	if res.Recipe.ID != "salep_pemulih" {
		t.Errorf("expected salep_pemulih recipe, got %s", res.Recipe.ID)
	}
}
