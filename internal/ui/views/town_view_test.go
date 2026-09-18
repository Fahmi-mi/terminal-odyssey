package views

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/Fahmi-mi/terminal-odyssey/internal/combat"
	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/save"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

func TestRenderViewsBorderAlignment(t *testing.T) {
	eng := engine.NewGame("Sang Petualang", "Oakhaven")

	// Prepare views for expedition
	_ = eng.StartExpedition(2)
	dungeonView := RenderDungeonView(eng, 80)
	if eng.ActiveExpedition.ActiveCombat == nil {
		enemy, _ := combat.NewEnemyByID("skeleton_scout")
		eng.ActiveExpedition.ActiveCombat = combat.NewCombatSession(eng.Player, enemy)
	}
	combatView := RenderCombatView(eng, 80)

	// Prepare dungeon view variants
	engDungeon := engine.NewGame("Sang Petualang", "Oakhaven")
	_ = engDungeon.StartExpedition(2)
	exp := engDungeon.ActiveExpedition
	exp.Rooms[0].IsResolved = true
	dungeonViewBranching := RenderDungeonView(engDungeon, 80)

	exp.CurrentRoomIdx = 5
	exp.Rooms[5].IsResolved = true
	dungeonViewSingle := RenderDungeonView(engDungeon, 80)

	exp.CurrentRoomIdx = 8
	exp.Rooms[8].IsResolved = true
	dungeonViewExit := RenderDungeonView(engDungeon, 80)

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
		{"DungeonViewUnresolved", dungeonView},
		{"DungeonViewBranching", dungeonViewBranching},
		{"DungeonViewSingleChoice", dungeonViewSingle},
		{"DungeonViewExit", dungeonViewExit},
		{"CombatView", combatView},
		{"SummaryView", summaryView},
		{"BlacksmithViewCraft", RenderBlacksmithView(eng, 0, 0, 0, 80)},
		{"BlacksmithViewArmory", RenderBlacksmithView(eng, 0, 1, 0, 80)},
		{"TrainingView", RenderTrainingView(eng, 0, 80)},
		{"MarketViewLocal", RenderMarketView(eng, 0, 0, 0, 80)},
		{"MarketViewCaravanLocked", RenderMarketView(eng, 0, 1, 0, 80)},
		{"MarketViewCaravanUnlocked", RenderMarketView(func() *engine.Engine {
			e2 := engine.NewGame("Sang Petualang", "Oakhaven")
			e2.Village.Buildings[settlement.BuildingCaravanPost] = 2
			return e2
		}(), 0, 1, 0, 80)},
		{"AlchemyViewLvl0", RenderAlchemyView(eng, 0, 80)},
		{"AlchemyViewLvl1", RenderAlchemyView(func() *engine.Engine {
			e2 := engine.NewGame("Sang Petualang", "Oakhaven")
			e2.Village.Buildings[settlement.BuildingApothecary] = 1
			return e2
		}(), 0, 80)},
		{"TavernViewRecruit", RenderTavernView(eng, 0, 0, 80)},
		{"TavernViewRumorAndDining", RenderTavernView(eng, 0, 1, 80)},
		{"TavernViewWithParty", RenderTavernView(func() *engine.Engine {
			e2 := engine.NewGame("Sang Petualang", "Oakhaven")
			e2.Village.Buildings[settlement.BuildingTavern] = 2
			_ = e2.HireCompanion("valen_rogue")
			return e2
		}(), 0, 0, 80)},
		{"TitleViewMain", RenderTitleView(TitleModeMain, 0, nil, "", 80)},
		{"TitleViewSlots", RenderTitleView(TitleModeSelectSlot, 0, []save.SaveSlotInfo{
			{SlotID: save.Slot1, Exists: true, DayCounter: 5, VillageName: "Oakhaven", PlayerName: "Pahlawan", SaveTime: "2026-09-18 12:00:00"},
			{SlotID: save.Slot2, Exists: false},
		}, "", 80)},
		{"TitleViewWithAlert", RenderTitleView(TitleModeSelectSlot, 0, []save.SaveSlotInfo{
			{SlotID: save.Slot1, Exists: true, DayCounter: 5, VillageName: "Oakhaven", PlayerName: "Pahlawan", SaveTime: "2026-09-18 12:00:00"},
			{SlotID: save.Slot2, Exists: false},
		}, "Slot simpanan masih kosong", 80)},
		{"SaveMenuView", RenderSaveMenuView(eng, 0, []save.SaveSlotInfo{
			{SlotID: save.Slot1, Exists: true, DayCounter: 5, VillageName: "Oakhaven", SaveTime: "2026-09-18 12:00:00"},
			{SlotID: save.Slot2, Exists: false},
			{SlotID: save.Slot3, Exists: false},
		}, 80)},
		{"SiegeViewVictory", RenderSiegeView(func() *engine.Engine {
			e2 := engine.NewGame("Sang Petualang", "Oakhaven")
			e2.Village.DefenseVal = 100
			e2.TriggerDirectSiege("bandit_raiders")
			return e2
		}(), 80)},
		{"SiegeViewDefeat", RenderSiegeView(func() *engine.Engine {
			e2 := engine.NewGame("Sang Petualang", "Oakhaven")
			e2.Village.DefenseVal = 5
			e2.TriggerDirectSiege("corrupted_legion")
			return e2
		}(), 80)},
		{"VictoryView", RenderVictoryView(func() *engine.Engine {
			e2 := engine.NewGame("Sang Petualang", "Oakhaven")
			e2.BossDefeated = true
			e2.Village.Buildings[settlement.BuildingTownHall] = 3
			e2.Village.Buildings[settlement.BuildingFortification] = 2
			e2.Village.RecalculateDefense()
			return e2
		}(), 80)},
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

func TestTownViewLineCount(t *testing.T) {
	eng := engine.NewGame("Sang Petualang", "Oakhaven")
	rendered := RenderTownView(eng, 80)
	lines := strings.Split(rendered, "\n")
	for i, l := range lines {
		t.Logf("%2d: %s", i+1, l)
	}
	t.Logf("TownView total line count: %d", len(lines))
	if len(lines) > 34 {
		t.Errorf("TownView exceeds expected height: %d lines (max 34)", len(lines))
	}
}

