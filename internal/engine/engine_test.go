package engine

import (
	"testing"

	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

func TestEngineInitialState(t *testing.T) {
	eng := NewGame("Petualang", "Oakhaven")

	if eng.DayCounter != 1 {
		t.Fatalf("expected day 1, got %d", eng.DayCounter)
	}
	if eng.CurrentSeason != settlement.SeasonSpring {
		t.Fatalf("expected Spring, got %v", eng.CurrentSeason)
	}
	if eng.CurrentState != StateTownMenu {
		t.Fatalf("expected StateTownMenu, got %v", eng.CurrentState)
	}
}

func TestEnginePassDayAndSeasonTransition(t *testing.T) {
	eng := NewGame("Petualang", "Oakhaven")
	eng.DaysPerSeason = 5 // Fast seasons for test

	// Pass 4 days -> still Spring
	for i := 0; i < 4; i++ {
		eng.PassDay()
	}
	if eng.CurrentSeason != settlement.SeasonSpring {
		t.Fatalf("expected Spring on day 5, got %v", eng.CurrentSeason)
	}

	// 5th day pass -> day 6 -> transition to Summer
	result := eng.PassDay()
	if !result.SeasonChanged {
		t.Fatalf("expected season change on day 6")
	}
	if eng.CurrentSeason != settlement.SeasonSummer {
		t.Fatalf("expected Summer on day 6, got %v", eng.CurrentSeason)
	}
}

func TestEngineStateSwitch(t *testing.T) {
	eng := NewGame("Petualang", "Oakhaven")
	eng.SetAlert("Test Alert")

	eng.SwitchState(StateWorkerAssign)
	if eng.CurrentState != StateWorkerAssign {
		t.Fatalf("expected StateWorkerAssign, got %v", eng.CurrentState)
	}
	if eng.PreviousState != StateTownMenu {
		t.Fatalf("expected PreviousState StateTownMenu, got %v", eng.PreviousState)
	}
	if eng.StatusAlert != "" {
		t.Fatalf("expected alert to be cleared on state switch")
	}
}

func TestEngineExpeditionFlow(t *testing.T) {
	eng := NewGame("Petualang", "Oakhaven")
	eng.Village.Rations = 10
	initialTreasury := eng.Village.Treasury

	// Test insufficient rations error
	errTooMany := eng.StartExpedition(20)
	if errTooMany == nil {
		t.Errorf("expected error when taking more rations than village has")
	}

	// Test valid expedition start
	errStart := eng.StartExpedition(3)
	if errStart != nil {
		t.Fatalf("unexpected start expedition error: %v", errStart)
	}
	if eng.Village.Rations != 7 {
		t.Errorf("expected 7 village rations remaining, got %d", eng.Village.Rations)
	}
	if eng.ActiveExpedition == nil {
		t.Fatalf("expected ActiveExpedition to be initialized")
	}

	// Simulate finding loot and evacuating
	eng.ActiveExpedition.GoldFound = 50
	eng.ActiveExpedition.LumberFound = 20
	eng.FinishExpedition(true)

	if eng.CurrentState != StateExpeditionSummary {
		t.Errorf("expected StateExpeditionSummary, got %v", eng.CurrentState)
	}
	if eng.Village.Treasury != initialTreasury+50 {
		t.Errorf("expected treasury %d, got %d", initialTreasury+50, eng.Village.Treasury)
	}

	// Test defeat flow
	eng.Village.Rations = 5
	initialDay := eng.DayCounter
	_ = eng.StartExpedition(2)
	eng.ActiveExpedition.GoldFound = 100
	eng.ActiveExpedition.IsDefeated = true
	eng.FinishExpedition(false)

	if eng.DayCounter != initialDay+1 {
		t.Errorf("expected day to advance after medical rescue, got %d", eng.DayCounter)
	}
	if eng.Player.HP != 25 {
		t.Errorf("expected player HP 25 after rescue, got %d", eng.Player.HP)
	}

	// Test PassDay recovers +20 HP
	eng.Village.Rations = 20
	eng.PassDay()
	if eng.Player.HP != 45 {
		t.Errorf("expected player HP 45 after daily rest (25+20), got %d", eng.Player.HP)
	}
}

func TestEngineBlacksmithCraftAndRepair(t *testing.T) {
	eng := NewGame("CraftTester", "Oakhaven")

	// Error if Blacksmith not built
	errNoBS := eng.CraftWeapon("iron_broadsword")
	if errNoBS == nil {
		t.Errorf("expected error crafting weapon without Blacksmith")
	}

	// Build Blacksmith Lvl 1
	eng.Village.Buildings[settlement.BuildingBlacksmith] = 1
	eng.Village.Lumber = 100
	eng.Village.Stone = 100
	eng.Village.Treasury = 200

	// Craft iron_broadsword
	errCraft := eng.CraftWeapon("iron_broadsword")
	if errCraft != nil {
		t.Fatalf("unexpected craft error: %v", errCraft)
	}
	if eng.Player.EquippedWeapon.ID != "iron_broadsword" {
		t.Errorf("expected equipped weapon iron_broadsword, got %s", eng.Player.EquippedWeapon.ID)
	}
	if eng.Village.Lumber != 90 || eng.Village.Stone != 85 || eng.Village.Treasury != 160 {
		t.Errorf("resources not deducted correctly after crafting")
	}

	// Damage weapon and test repair
	eng.Player.EquippedWeapon.Durability = 20
	errRepair := eng.RepairEquippedWeapon()
	if errRepair != nil {
		t.Fatalf("unexpected repair error: %v", errRepair)
	}
	if eng.Player.EquippedWeapon.Durability != eng.Player.EquippedWeapon.MaxDura {
		t.Errorf("expected durability to be fully restored, got %d", eng.Player.EquippedWeapon.Durability)
	}

	// Repair when already full
	errFullRepair := eng.RepairEquippedWeapon()
	if errFullRepair == nil {
		t.Errorf("expected error repairing fully intact weapon")
	}

	// Switch back to starter weapon (rusty_sword)
	errSwitch := eng.SwitchWeapon("rusty_sword")
	if errSwitch != nil {
		t.Fatalf("unexpected error switching weapon: %v", errSwitch)
	}
	if eng.Player.EquippedWeapon.ID != "rusty_sword" {
		t.Errorf("expected equipped weapon rusty_sword, got %s", eng.Player.EquippedWeapon.ID)
	}

	// Test repairing unequipped weapon in storage
	for i := range eng.Player.OwnedWeapons {
		if eng.Player.OwnedWeapons[i].ID == "iron_broadsword" {
			eng.Player.OwnedWeapons[i].Durability = 30
		}
	}
	errRepairStored := eng.RepairWeaponByID("iron_broadsword")
	if errRepairStored != nil {
		t.Fatalf("unexpected error repairing stored weapon: %v", errRepairStored)
	}
}

func TestEngineTrainStat(t *testing.T) {
	eng := NewGame("TrainTester", "Oakhaven")

	// Error if Training Grounds not built
	errNoTG := eng.TrainStat("Might")
	if errNoTG == nil {
		t.Errorf("expected error training stat without Training Grounds")
	}

	// Build Training Grounds Lvl 1
	eng.Village.Buildings[settlement.BuildingTrainingGround] = 1
	eng.Village.Treasury = 300
	eng.Village.Rations = 50

	initialMight := eng.Player.Stats.Might
	errTrain := eng.TrainStat("Might")
	if errTrain != nil {
		t.Fatalf("unexpected train error: %v", errTrain)
	}
	if eng.Player.Stats.Might != initialMight+1 {
		t.Errorf("expected might %d, got %d", initialMight+1, eng.Player.Stats.Might)
	}

	// Train Resolve and verify MaxHP scales
	initialHP := eng.Player.MaxHP
	errTrainRes := eng.TrainStat("Resolve")
	if errTrainRes != nil {
		t.Fatalf("unexpected train resolve error: %v", errTrainRes)
	}
	if eng.Player.MaxHP != initialHP+5 {
		t.Errorf("expected MaxHP %d, got %d", initialHP+5, eng.Player.MaxHP)
	}
}



