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


