package views

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
)

func TestRenderViewsBorderAlignment(t *testing.T) {
	eng := engine.NewGame("Sang Petualang", "Oakhaven")

	// Prepare views for expedition
	_ = eng.StartExpedition(2)
	dungeonView := RenderDungeonView(eng, 80)
	combatView := RenderCombatView(eng, 80)

	// Prepare summary view
	eng.FinishExpedition(true)
	summaryView := RenderExpeditionSummaryView(eng, 80)

	// Prepare town view with many daily logs
	engWithLogs := engine.NewGame("Sang Petualang", "Oakhaven")
	engWithLogs.DailyLogs = []string{
		"Log 1: Panen pertama",
		"Log 2: Kayu terkumpul",
		"Log 3: Batu tertambang",
		"Log 4: Warga baru datang",
		"Log 5: Satu peristiwa panjang yang terjadi di perkemahan dan melampaui batas lebar teks standar",
	}

	testCases := []struct {
		name     string
		rendered string
	}{
		{"TownView", RenderTownView(eng, 80)},
		{"TownViewManyLogs", RenderTownView(engWithLogs, 80)},
		{"WorkerView", RenderWorkerView(eng, 0, 80)},
		{"BuildView", RenderBuildView(eng, 0, 80)},
		{"DungeonView", dungeonView},
		{"CombatView", combatView},
		{"SummaryView", summaryView},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lines := strings.Split(tc.rendered, "\n")
			for i, line := range lines {
				if len(strings.TrimSpace(line)) == 0 {
					continue
				}
				w := lipgloss.Width(line)
				if w != 78 {
					t.Errorf("line %d has width %d (expected 78): %q", i, w, line)
				}
			}
		})
	}
}

func TestSelectPriorityDailyLogs(t *testing.T) {
	// Case 1: <= 3 logs returned as is
	initialLogs := []string{"Log A", "Log B"}
	res1 := selectPriorityDailyLogs(initialLogs, 3)
	if len(res1) != 2 || res1[0] != "Log A" || res1[1] != "Log B" {
		t.Errorf("expected 2 logs preserved, got %v", res1)
	}

	// Case 2: Priority ordering among > 3 logs
	logs := []string{
		"[i] Suasana perkemahan tenang",
		"[+] Produksi Harian : +4 Ransum | +3 Kayu | +2 Batu",
		"[+] Seorang pengembara terkesan dengan stabilitas desa dan memutuskan menetap",
		"[!] KELAPARAN: Defisit 3 ransum! Warga menderita kekurangan pangan",
		"[*] PERUBAHAN MUSIM: Memasuki Musim Dingin",
	}

	res2 := selectPriorityDailyLogs(logs, 3)
	if len(res2) != 3 {
		t.Fatalf("expected 3 logs, got %d", len(res2))
	}

	// Priority 1: [!] KELAPARAN
	if !strings.Contains(res2[0], "KELAPARAN") {
		t.Errorf("expected highest priority to be KELAPARAN, got %s", res2[0])
	}
	// Priority 2: [*] PERUBAHAN MUSIM
	if !strings.Contains(res2[1], "PERUBAHAN MUSIM") {
		t.Errorf("expected second priority to be PERUBAHAN MUSIM, got %s", res2[1])
	}
	// Priority 3: [+] Seorang pengembara
	if !strings.Contains(res2[2], "pengembara") {
		t.Errorf("expected third priority to be pengembara, got %s", res2[2])
	}
}

