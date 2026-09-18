package views

import (
	"strings"
	"testing"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

func TestRenderAlchemyView(t *testing.T) {
	eng := engine.NewGame("Alchemist", "Oakhaven")

	// 1. Render without building (Level 0)
	v0 := RenderAlchemyView(eng, 0, 80)
	if !strings.Contains(v0, "LABORATORIUM ALKIMIA") {
		t.Errorf("expected title in view")
	}
	if !strings.Contains(v0, "Terkunci") {
		t.Errorf("expected Terkunci status for recipes when lab is level 0")
	}
	if !strings.Contains(v0, "▶ ") {
		t.Errorf("expected solid triangle cursor ▶")
	}

	// 2. Render with building Level 1
	eng.Village.Buildings[settlement.BuildingApothecary] = 1
	eng.Village.Treasury = 100
	eng.Village.AddCommodity("herbal_salve", 5)

	v1 := RenderAlchemyView(eng, 0, 80)
	if !strings.Contains(v1, "Siap Racik") {
		t.Errorf("expected Siap Racik for salep_pemulih with lab lvl 1 and ingredients")
	}

	// 3. Render with alert
	eng.SetAlert("Racikan Alkimia Sempurna!")
	vAlert := RenderAlchemyView(eng, 0, 80)
	if !strings.Contains(vAlert, "Racikan Alkimia Sempurna!") {
		t.Errorf("expected alert message rendered in view")
	}
}
