package save

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

func TestSaveAndLoadRoundtrip(t *testing.T) {
	tempDir := t.TempDir()

	eng := engine.NewGame("Pahlawan Uji", "Desa Megah")
	eng.DayCounter = 14
	eng.CurrentSeason = settlement.SeasonAutumn
	eng.Village.Treasury = 850
	eng.Village.Lumber = 120
	eng.Village.Stone = 95
	eng.Village.Rations = 80
	eng.Village.Buildings[settlement.BuildingTownHall] = 2
	eng.Village.Buildings[settlement.BuildingFortification] = 1
	eng.Village.RecalculateDefense()
	eng.Player.HP = 88
	eng.Player.Stats.Might = 14
	eng.BossDefeated = true

	// Save to Slot 1
	err := SaveGame(eng, Slot1, tempDir)
	if err != nil {
		t.Fatalf("failed to save game: %v", err)
	}

	// Verify file exists
	expectedPath := filepath.Join(tempDir, Slot1+".json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("save file does not exist at %s", expectedPath)
	}

	// Load back
	loadedEng, err := LoadGame(Slot1, tempDir)
	if err != nil {
		t.Fatalf("failed to load game: %v", err)
	}

	if loadedEng.DayCounter != 14 {
		t.Errorf("expected DayCounter 14, got %d", loadedEng.DayCounter)
	}
	if loadedEng.CurrentSeason != settlement.SeasonAutumn {
		t.Errorf("expected season Autumn, got %v", loadedEng.CurrentSeason)
	}
	if loadedEng.Village.Name != "Desa Megah" {
		t.Errorf("expected village name Desa Megah, got %s", loadedEng.Village.Name)
	}
	if loadedEng.Village.Treasury != 850 {
		t.Errorf("expected treasury 850, got %d", loadedEng.Village.Treasury)
	}
	if loadedEng.Player.HP != 88 {
		t.Errorf("expected player HP 88, got %d", loadedEng.Player.HP)
	}
	if loadedEng.Player.Stats.Might != 14 {
		t.Errorf("expected player Might 14, got %d", loadedEng.Player.Stats.Might)
	}
	if !loadedEng.BossDefeated {
		t.Errorf("expected BossDefeated to be true")
	}
	if loadedEng.Village.DefenseVal <= 0 {
		t.Errorf("expected recalculated positive defense")
	}
}

func TestListSaveSlots(t *testing.T) {
	tempDir := t.TempDir()

	// Initially all slots are empty
	slots, err := ListSaveSlots(tempDir)
	if err != nil {
		t.Fatalf("failed to list save slots: %v", err)
	}
	if len(slots) != 4 {
		t.Fatalf("expected 4 standard slots, got %d", len(slots))
	}
	for _, s := range slots {
		if s.Exists {
			t.Errorf("expected slot %s to not exist initially", s.SlotID)
		}
	}

	// Create a save in Slot 2
	eng := engine.NewGame("Ksatria", "Benteng Barat")
	eng.DayCounter = 7
	if err := SaveGame(eng, Slot2, tempDir); err != nil {
		t.Fatalf("failed to save slot 2: %v", err)
	}

	// Re-list
	slots, err = ListSaveSlots(tempDir)
	if err != nil {
		t.Fatalf("failed to list save slots: %v", err)
	}

	foundSlot2 := false
	for _, s := range slots {
		if s.SlotID == Slot2 {
			foundSlot2 = true
			if !s.Exists {
				t.Errorf("expected slot 2 to exist")
			}
			if s.DayCounter != 7 {
				t.Errorf("expected DayCounter 7, got %d", s.DayCounter)
			}
			if s.VillageName != "Benteng Barat" {
				t.Errorf("expected village name Benteng Barat, got %s", s.VillageName)
			}
			if s.PlayerName != "Ksatria" {
				t.Errorf("expected player name Ksatria, got %s", s.PlayerName)
			}
		}
	}
	if !foundSlot2 {
		t.Errorf("expected slot 2 to be in slots list")
	}

	// Test DeleteSave
	if err := DeleteSave(Slot2, tempDir); err != nil {
		t.Fatalf("failed to delete save: %v", err)
	}
	slots, _ = ListSaveSlots(tempDir)
	for _, s := range slots {
		if s.SlotID == Slot2 && s.Exists {
			t.Errorf("expected slot 2 to be deleted")
		}
	}
}
