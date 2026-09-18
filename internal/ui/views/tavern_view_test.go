package views

import (
	"strings"
	"testing"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

func TestRenderTavernView(t *testing.T) {
	eng := engine.NewGame("Commander", "Oakhaven")

	// 1. Render without tavern building (Level 0)
	v0 := RenderTavernView(eng, 0, 0, 80)
	if !strings.Contains(v0, "KEDAI MINUM") {
		t.Errorf("expected title in view")
	}
	if !strings.Contains(v0, "Terkunci") {
		t.Errorf("expected Terkunci status when tavern is level 0")
	}
	if !strings.Contains(v0, "▶ ") {
		t.Errorf("expected solid triangle cursor ▶")
	}

	// 2. Render Tab 0 with Tavern Level 1
	eng.Village.Buildings[settlement.BuildingTavern] = 1
	eng.Village.Treasury = 200

	v1 := RenderTavernView(eng, 0, 0, 80)
	if !strings.Contains(v1, "Tersedia") {
		t.Errorf("expected Tersedia status when companion is recruitable")
	}

	// 3. Render after hiring
	_ = eng.HireCompanion("valen_rogue")
	vHired := RenderTavernView(eng, 0, 0, 80)
	if !strings.Contains(vHired, "Bergabung") {
		t.Errorf("expected Bergabung status for hired companion")
	}

	// 4. Render Tab 1 (Rumors and Dining)
	vTab1 := RenderTavernView(eng, 0, 1, 80)
	if !strings.Contains(vTab1, "SANTAPAN HANGAT") {
		t.Errorf("expected dining section in Tab 1")
	}
}
