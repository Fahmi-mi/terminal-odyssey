package tavern

import (
	"testing"
)

func TestTavernManager_HireAndDismiss(t *testing.T) {
	tm, err := NewTavernManager()
	if err != nil {
		t.Fatalf("failed to initialize tavern manager: %v", err)
	}

	// Case 1: Tavern level 0
	err = tm.CanHire("valen_rogue", 0, 100, 0)
	if err == nil {
		t.Errorf("expected error when tavern is level 0")
	}

	// Case 2: Insufficient gold
	err = tm.CanHire("valen_rogue", 1, 10, 0)
	if err == nil {
		t.Errorf("expected error for insufficient gold")
	}

	// Case 3: Party limit reached
	err = tm.CanHire("valen_rogue", 1, 100, 1)
	if err == nil {
		t.Errorf("expected error when party limit is reached")
	}

	// Case 4: Successful hire
	merc, err := tm.Hire("valen_rogue", 1, 100, 0)
	if err != nil {
		t.Fatalf("unexpected error hiring mercenary: %v", err)
	}
	if !merc.IsHired {
		t.Errorf("expected mercenary to be marked as hired")
	}

	// Case 5: Cannot hire twice
	err = tm.CanHire("valen_rogue", 2, 100, 1)
	if err == nil {
		t.Errorf("expected error when trying to hire already hired mercenary")
	}

	// Case 6: Dismissal
	err = tm.Dismiss("valen_rogue")
	if err != nil {
		t.Fatalf("unexpected error dismissing mercenary: %v", err)
	}
	if merc.IsHired {
		t.Errorf("expected mercenary to be unhired after dismissal")
	}
}

func TestTavernManager_RestAndRumors(t *testing.T) {
	tm, err := NewTavernManager()
	if err != nil {
		t.Fatalf("failed to initialize tavern manager: %v", err)
	}

	err = tm.CanRest(5, 1)
	if err == nil {
		t.Errorf("expected error with insufficient gold for rest")
	}

	err = tm.CanRest(15, 0)
	if err == nil {
		t.Errorf("expected error with insufficient rations for rest")
	}

	err = tm.CanRest(15, 2)
	if err != nil {
		t.Errorf("expected valid rest with sufficient resources: %v", err)
	}

	rumor := tm.GetRandomRumor()
	if rumor == "" {
		t.Errorf("expected non-empty rumor")
	}
}
