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
