package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Fahmi-mi/terminal-odyssey/internal/combat"
	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

func makeKeyMsg(k string) tea.KeyMsg {
	switch k {
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "left":
		return tea.KeyMsg{Type: tea.KeyLeft}
	case "right":
		return tea.KeyMsg{Type: tea.KeyRight}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
	}
}

func TestAppModel_MarketNavigationAndKeybindings(t *testing.T) {
	eng := engine.NewGame("Petualang", "Oakhaven")
	eng.Village.Treasury = 500
	eng.Village.Lumber = 30
	eng.Village.Stone = 30
	eng.Village.Rations = 30
	eng.Village.Buildings[settlement.BuildingCaravanPost] = 1

	model := NewAppModel(eng)

	// 1. From Town Menu, press '1' to enter Market
	model.Update(makeKeyMsg("1"))
	if model.Engine.CurrentState != engine.StateMarketTrade {
		t.Fatalf("expected StateMarketTrade, got %v", model.Engine.CurrentState)
	}

	// 2. Tab switching
	model.Update(makeKeyMsg("tab"))
	if model.marketTab != 1 {
		t.Errorf("expected marketTab 1 after tab press, got %d", model.marketTab)
	}
	model.Update(makeKeyMsg("1"))
	if model.marketTab != 0 {
		t.Errorf("expected marketTab 0 after pressing '1', got %d", model.marketTab)
	}
	model.Update(makeKeyMsg("2"))
	if model.marketTab != 1 {
		t.Errorf("expected marketTab 1 after pressing '2', got %d", model.marketTab)
	}
	model.Update(makeKeyMsg("tab"))
	if model.marketTab != 0 {
		t.Errorf("expected marketTab 0 after pressing tab again, got %d", model.marketTab)
	}

	// 3. Navigation and buying/selling
	model.selectedMarketIdx = 0 // Lumber
	initialLumber := eng.Village.Lumber
	model.Update(makeKeyMsg("b"))
	if eng.Village.Lumber != initialLumber+1 {
		t.Errorf("expected lumber to increase by 1 after pressing 'b', got %d", eng.Village.Lumber)
	}

	model.Update(makeKeyMsg("s"))
	if eng.Village.Lumber != initialLumber {
		t.Errorf("expected lumber to return to initial after pressing 's', got %d", eng.Village.Lumber)
	}

	// Navigate with down
	model.Update(makeKeyMsg("down"))
	if model.selectedMarketIdx != 1 {
		t.Errorf("expected selectedMarketIdx 1 after down, got %d", model.selectedMarketIdx)
	}

	// Verify 'j' does NOT navigate (strictly no undocumented keys!)
	model.Update(makeKeyMsg("j"))
	if model.selectedMarketIdx != 1 {
		t.Errorf("expected 'j' to be ignored, but index changed to %d", model.selectedMarketIdx)
	}

	// Navigate with up
	model.Update(makeKeyMsg("up"))
	if model.selectedMarketIdx != 0 {
		t.Errorf("expected selectedMarketIdx 0 after up, got %d", model.selectedMarketIdx)
	}

	// Verify 'k' does NOT navigate
	model.Update(makeKeyMsg("k"))
	if model.selectedMarketIdx != 0 {
		t.Errorf("expected 'k' to be ignored, but index changed to %d", model.selectedMarketIdx)
	}

	// 4. Switch to Pos Kafilah tab and dispatch
	model.Update(makeKeyMsg("2"))
	model.selectedRouteIdx = 0 // Riverfall
	model.Update(makeKeyMsg("enter"))
	if len(eng.Caravans.ActiveCaravans) != 1 {
		t.Fatalf("expected 1 active caravan after pressing enter on Riverfall route, got %d", len(eng.Caravans.ActiveCaravans))
	}

	// 5. Exit back to Town Menu via ESC
	model.Update(makeKeyMsg("esc"))
	if model.Engine.CurrentState != engine.StateTownMenu {
		t.Fatalf("expected StateTownMenu after pressing ESC, got %v", model.Engine.CurrentState)
	}
}

func TestAppModel_TrainingGroundsStrictKeybindings(t *testing.T) {
	eng := engine.NewGame("Petualang", "Oakhaven")
	eng.Village.Buildings[settlement.BuildingTrainingGround] = 1
	eng.Village.Treasury = 500
	eng.Village.Rations = 100

	model := NewAppModel(eng)

	// Switch to training
	model.Engine.SwitchState(engine.StateTrainingGrounds)
	model.selectedStatIdx = 0

	// '1', '2', '3', '4' must NOT trigger jumping (not listed on page)
	model.Update(makeKeyMsg("3"))
	if model.selectedStatIdx != 0 {
		t.Errorf("expected '3' to be ignored in training grounds, got index %d", model.selectedStatIdx)
	}

	// 'down' navigates down
	model.Update(makeKeyMsg("down"))
	if model.selectedStatIdx != 1 {
		t.Errorf("expected selectedStatIdx 1 after down, got %d", model.selectedStatIdx)
	}

	// 'j' does NOT navigate down
	model.Update(makeKeyMsg("j"))
	if model.selectedStatIdx != 1 {
		t.Errorf("expected 'j' to be ignored, got index %d", model.selectedStatIdx)
	}

	// 'enter' trains stat
	initialAgility := eng.Player.Stats.Agility
	model.Update(makeKeyMsg("enter"))
	if eng.Player.Stats.Agility != initialAgility+1 {
		t.Errorf("expected agility trained (+1), got %d", eng.Player.Stats.Agility)
	}

	// 'esc' returns to town
	model.Update(makeKeyMsg("esc"))
	if model.Engine.CurrentState != engine.StateTownMenu {
		t.Fatalf("expected StateTownMenu after ESC, got %v", model.Engine.CurrentState)
	}
}

func TestAppModel_ViewRendering(t *testing.T) {
	eng := engine.NewGame("Petualang", "Oakhaven")
	model := NewAppModel(eng)

	viewTown := model.View()
	if len(viewTown) == 0 {
		t.Errorf("expected non-empty town view")
	}

	model.Engine.SwitchState(engine.StateMarketTrade)
	viewMarket := model.View()
	if len(viewMarket) == 0 {
		t.Errorf("expected non-empty market view")
	}

	model.Engine.SwitchState(engine.StateAlchemyLab)
	viewAlchemy := model.View()
	if len(viewAlchemy) == 0 {
		t.Errorf("expected non-empty alchemy view")
	}

	model.Engine.SwitchState(engine.StateTavernRecruit)
	viewTavern := model.View()
	if len(viewTavern) == 0 {
		t.Errorf("expected non-empty tavern view")
	}
}

func TestAppModel_AlchemyNavigationAndBrewing(t *testing.T) {
	eng := engine.NewGame("Alchemist", "Oakhaven")
	model := NewAppModel(eng)

	// Press '7' without lab -> alert set, stays in town menu
	model.Update(makeKeyMsg("7"))
	if model.Engine.CurrentState != engine.StateTownMenu {
		t.Errorf("expected StateTownMenu when lab not built")
	}

	// Build lab level 1
	eng.Village.Buildings[settlement.BuildingApothecary] = 1
	eng.Village.Treasury = 200
	eng.Village.AddCommodity("herbal_salve", 10)

	// Press '7' to enter
	model.Update(makeKeyMsg("7"))
	if model.Engine.CurrentState != engine.StateAlchemyLab {
		t.Fatalf("expected StateAlchemyLab, got %v", model.Engine.CurrentState)
	}

	// Navigate with down and up
	model.Update(makeKeyMsg("down"))
	if model.selectedAlchemyIdx != 1 {
		t.Errorf("expected selectedAlchemyIdx 1, got %d", model.selectedAlchemyIdx)
	}
	model.Update(makeKeyMsg("up"))
	if model.selectedAlchemyIdx != 0 {
		t.Errorf("expected selectedAlchemyIdx 0, got %d", model.selectedAlchemyIdx)
	}

	// Press enter to brew salep_pemulih
	model.Update(makeKeyMsg("enter"))
	if eng.Player.GetPotionCount("salep_pemulih") < 1 {
		t.Errorf("expected brewed potion in player pouch")
	}

	// Press esc to return to town
	model.Update(makeKeyMsg("esc"))
	if model.Engine.CurrentState != engine.StateTownMenu {
		t.Errorf("expected StateTownMenu after esc, got %v", model.Engine.CurrentState)
	}
}

func TestAppModel_TavernNavigationAndCompanions(t *testing.T) {
	eng := engine.NewGame("Commander", "Oakhaven")
	model := NewAppModel(eng)

	// Press '5' without tavern -> stays in town menu
	model.Update(makeKeyMsg("5"))
	if model.Engine.CurrentState != engine.StateTownMenu {
		t.Errorf("expected StateTownMenu when tavern not built")
	}

	// Build tavern level 1
	eng.Village.Buildings[settlement.BuildingTavern] = 1
	eng.Village.Treasury = 300
	eng.Village.Rations = 20

	// Press '5' to enter
	model.Update(makeKeyMsg("5"))
	if model.Engine.CurrentState != engine.StateTavernRecruit {
		t.Fatalf("expected StateTavernRecruit, got %v", model.Engine.CurrentState)
	}

	// Tab switching
	model.Update(makeKeyMsg("tab"))
	if model.tavernTab != 1 {
		t.Errorf("expected tavernTab 1, got %d", model.tavernTab)
	}
	// On Tab 1: order dining rest
	eng.Player.Sanity = 50
	eng.Player.MaxSanity = 100
	model.Update(makeKeyMsg("enter"))
	if eng.Player.Sanity <= 50 {
		t.Errorf("expected sanity recovery after tavern meal")
	}

	// Switch back to Tab 0
	model.Update(makeKeyMsg("1"))
	if model.tavernTab != 0 {
		t.Errorf("expected tavernTab 0, got %d", model.tavernTab)
	}

	// Hire first mercenary (Valen Rogue)
	model.selectedTavernIdx = 0
	model.Update(makeKeyMsg("enter"))
	if len(eng.Player.Party) != 1 {
		t.Fatalf("expected 1 hired companion, got %d", len(eng.Player.Party))
	}

	// Press enter again to dismiss
	model.Update(makeKeyMsg("enter"))
	if len(eng.Player.Party) != 0 {
		t.Errorf("expected 0 companions after dismissal, got %d", len(eng.Player.Party))
	}

	// Press esc to return to town
	model.Update(makeKeyMsg("esc"))
	if model.Engine.CurrentState != engine.StateTownMenu {
		t.Errorf("expected StateTownMenu after esc, got %v", model.Engine.CurrentState)
	}
}

func TestAppModel_PotionUseInDungeonAndCombat(t *testing.T) {
	eng := engine.NewGame("Explorer", "Oakhaven")
	eng.Village.Rations = 10
	_ = eng.StartExpedition(3)

	exp := eng.ActiveExpedition
	exp.Player.HP = 40
	exp.Player.MaxHP = 100
	exp.Player.AddPotion("salep_pemulih", 1)

	model := NewAppModel(eng)

	// 1. In dungeon exploration, press 'p' to drink potion
	model.Update(makeKeyMsg("p"))
	if exp.Player.HP != 75 { // 40 + 35 = 75
		t.Errorf("expected player HP 75 after drinking salep_pemulih in dungeon, got %d", exp.Player.HP)
	}

	// 2. In combat, drink strength elixir
	exp.Player.AddPotion("eliksir_kekuatan", 1)
	room := exp.CurrentRoom()
	room.Enemy = &combat.Enemy{
		ID:         "rat",
		Name:       "Tikus Raksasa",
		HP:         50,
		MaxHP:      50,
		MinDamage:  0,
		MaxDamage:  0,
		Initiative: 2,
	}
	model.Engine.SwitchState(engine.StateCombatTurn)
	exp.ActiveCombat = combat.NewCombatSession(exp.Player, room.Enemy)

	model.Update(makeKeyMsg("p"))
	if exp.ActiveCombat.TemporaryAtkBuff != 8 {
		t.Errorf("expected TemporaryAtkBuff 8, got %d", exp.ActiveCombat.TemporaryAtkBuff)
	}
}

