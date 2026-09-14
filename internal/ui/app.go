package ui

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/combat"
	"github.com/Fahmi-mi/terminal-odyssey/internal/dungeon"
	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
	"github.com/Fahmi-mi/terminal-odyssey/internal/ui/views"
)

// AppModel is the root Bubble Tea model
type AppModel struct {
	Engine *engine.Engine

	width  int
	height int

	selectedWorkerIdx int
	selectedBuildIdx  int
	selectedRecipeIdx int
	selectedStatIdx   int
}

// NewAppModel creates a fresh TUI model
func NewAppModel(eng *engine.Engine) *AppModel {
	return &AppModel{
		Engine:            eng,
		width:             80,
		height:            24,
		selectedWorkerIdx: 0,
		selectedBuildIdx:  0,
		selectedRecipeIdx: 0,
		selectedStatIdx:   0,
	}
}

// Init initializes the Bubble Tea loop
func (m *AppModel) Init() tea.Cmd {
	return nil
}

// Update processes incoming messages and keyboard events
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.Engine.CurrentState {
		case engine.StateTownMenu:
			return m.updateTownMenu(msg)
		case engine.StateWorkerAssign:
			return m.updateWorkerAssign(msg)
		case engine.StateTownBuild:
			return m.updateTownBuild(msg)
		case engine.StateBlacksmithCraft:
			return m.updateBlacksmithCraft(msg)
		case engine.StateTrainingGrounds:
			return m.updateTrainingGrounds(msg)
		case engine.StateDungeonExplore:
			return m.updateDungeonExplore(msg)
		case engine.StateCombatTurn:
			return m.updateCombatTurn(msg)
		case engine.StateExpeditionSummary:
			return m.updateExpeditionSummary(msg)
		default:
			// Fallback back to town menu
			if msg.String() == "esc" || msg.String() == "q" {
				m.Engine.SwitchState(engine.StateTownMenu)
			}
		}
	}

	return m, nil
}

func (m *AppModel) updateTownMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "Q":
		return m, tea.Quit
	case "d", "D":
		m.Engine.PassDay()
	case "w", "W":
		m.Engine.SwitchState(engine.StateWorkerAssign)
	case "2":
		m.Engine.SwitchState(engine.StateTownBuild)
	case "1":
		m.Engine.SetAlert("Pasar & Perdagangan antar-kota sedang dipersiapkan (Milestone 4)")
	case "3":
		if m.Engine.Village.Buildings[settlement.BuildingBlacksmith] < 1 {
			m.Engine.SetAlert("Bengkel Pandai Besi belum didirikan! Bangun di menu [2] Pembangunan")
		} else {
			m.Engine.SwitchState(engine.StateBlacksmithCraft)
		}
	case "4":
		if m.Engine.Village.Buildings[settlement.BuildingTrainingGround] < 1 {
			m.Engine.SetAlert("Pusat Latihan belum didirikan! Bangun di menu [2] Pembangunan")
		} else {
			m.Engine.SwitchState(engine.StateTrainingGrounds)
		}
	case "5":
		m.Engine.SetAlert("Kedai Minum belum memiliki rumor baru hari ini")
	case "6":
		rationsToTake := 3
		if m.Engine.Village.Rations < rationsToTake {
			rationsToTake = m.Engine.Village.Rations
		}
		if err := m.Engine.StartExpedition(rationsToTake); err != nil {
			m.Engine.SetAlert(err.Error())
		}
	}
	return m, nil
}

func (m *AppModel) updateWorkerAssign(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	maxRoles := len(views.RoleItems)

	switch msg.String() {
	case "esc":
		m.Engine.SwitchState(engine.StateTownMenu)
	case "up", "k":
		if m.selectedWorkerIdx > 0 {
			m.selectedWorkerIdx--
		} else {
			m.selectedWorkerIdx = maxRoles - 1
		}
		m.Engine.ClearAlert()
	case "down", "j":
		if m.selectedWorkerIdx < maxRoles-1 {
			m.selectedWorkerIdx++
		} else {
			m.selectedWorkerIdx = 0
		}
		m.Engine.ClearAlert()
	case "right", "+":
		role := views.RoleItems[m.selectedWorkerIdx].Role
		if err := m.Engine.Village.AssignWorker(role); err != nil {
			m.Engine.SetAlert(err.Error())
		} else {
			m.Engine.SetAlert(fmt.Sprintf("Berhasil menugaskan +1 %s", role))
		}
	case "left", "-":
		role := views.RoleItems[m.selectedWorkerIdx].Role
		if err := m.Engine.Village.UnassignWorker(role); err != nil {
			m.Engine.SetAlert(err.Error())
		} else {
			m.Engine.SetAlert(fmt.Sprintf("Berhasil menarik 1 pekerja dari %s", role))
		}
	}
	return m, nil
}

func (m *AppModel) updateTownBuild(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	buildings := settlement.AllBuildingsList()
	maxBuildings := len(buildings)

	switch msg.String() {
	case "esc":
		m.Engine.SwitchState(engine.StateTownMenu)
	case "up", "k":
		if m.selectedBuildIdx > 0 {
			m.selectedBuildIdx--
		} else {
			m.selectedBuildIdx = maxBuildings - 1
		}
		m.Engine.ClearAlert()
	case "down", "j":
		if m.selectedBuildIdx < maxBuildings-1 {
			m.selectedBuildIdx++
		} else {
			m.selectedBuildIdx = 0
		}
		m.Engine.ClearAlert()
	case "1", "2", "3", "4", "5", "6", "7":
		val, _ := strconv.Atoi(msg.String())
		idx := val - 1
		if idx >= 0 && idx < maxBuildings {
			m.selectedBuildIdx = idx
			m.Engine.ClearAlert()
		}
	case "enter":
		bName := buildings[m.selectedBuildIdx]
		if err := m.Engine.Village.UpgradeBuilding(bName); err != nil {
			m.Engine.SetAlert(err.Error())
		} else {
			m.Engine.SetAlert(fmt.Sprintf("[+] %s berhasil ditingkatkan ke Level %d", bName, m.Engine.Village.Buildings[bName]))
		}
	}
	return m, nil
}

func (m *AppModel) updateDungeonExplore(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	exp := m.Engine.ActiveExpedition
	if exp == nil {
		m.Engine.SwitchState(engine.StateTownMenu)
		return m, nil
	}

	room := exp.CurrentRoom()

	switch msg.String() {
	case "esc":
		// Mundur dari ekspedisi dan bawa pulang jarahan
		m.Engine.FinishExpedition(true)
	case "enter":
		if room.Def.Type == dungeon.RoomTypeExit {
			m.Engine.FinishExpedition(true)
		} else if room.Def.Type == dungeon.RoomTypeCombat && !room.IsResolved {
			if exp.ActiveCombat == nil && room.Enemy != nil {
				exp.ActiveCombat = combat.NewCombatSession(exp.Player, room.Enemy)
			}
			m.Engine.SwitchState(engine.StateCombatTurn)
		} else if !room.IsResolved {
			switch room.Def.Type {
			case dungeon.RoomTypeTreasure:
				exp.ResolveTreasure()
			case dungeon.RoomTypeRest:
				exp.ResolveRest()
			case dungeon.RoomTypeMystery:
				exp.ResolveMystery()
				if exp.IsDefeated {
					m.Engine.FinishExpedition(false)
				}
			}
		} else {
			choices := exp.NextRoomChoices()
			if len(choices) == 1 {
				if err := exp.AdvanceRoom(); err != nil {
					m.Engine.SetAlert(err.Error())
				}
			} else if len(choices) > 1 {
				m.Engine.SetAlert("Pilih jalur lorong dengan menekan 1 atau 2")
			}
		}
	case "1":
		if room.IsResolved {
			choices := exp.NextRoomChoices()
			if len(choices) >= 1 {
				if err := exp.AdvanceToRoom(choices[0].GraphIdx); err != nil {
					m.Engine.SetAlert(err.Error())
				}
			}
		} else {
			m.Engine.SetAlert("Selesaikan peristiwa di ruangan ini terlebih dahulu")
		}
	case "2":
		if room.IsResolved {
			choices := exp.NextRoomChoices()
			if len(choices) >= 2 {
				if err := exp.AdvanceToRoom(choices[1].GraphIdx); err != nil {
					m.Engine.SetAlert(err.Error())
				}
			}
		} else {
			m.Engine.SetAlert("Selesaikan peristiwa di ruangan ini terlebih dahulu")
		}
	case "m", "M":
		if _, err := exp.ConsumeRation(); err != nil {
			m.Engine.SetAlert(err.Error())
		}
	case "o", "O":
		if err := exp.ConsumeTorch(); err != nil {
			m.Engine.SetAlert(err.Error())
		}
	}
	return m, nil
}

func (m *AppModel) updateCombatTurn(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	exp := m.Engine.ActiveExpedition
	if exp == nil || exp.ActiveCombat == nil {
		m.Engine.SwitchState(engine.StateDungeonExplore)
		return m, nil
	}

	session := exp.ActiveCombat

	// If combat is already resolved, wait for confirmation
	if session.IsOver {
		switch msg.String() {
		case "enter":
			if session.Won {
				exp.OnCombatWon()
				m.Engine.SwitchState(engine.StateDungeonExplore)
			} else if session.Fled {
				exp.ActiveCombat = nil
				m.Engine.SwitchState(engine.StateDungeonExplore)
			} else {
				// Player defeated
				m.Engine.FinishExpedition(false)
			}
		}
		return m, nil
	}

	// Active turn actions
	switch msg.String() {
	case "1":
		session.PlayerAttack()
		if session.Player.HP <= 0 {
			session.IsOver = true
			session.Won = false
		}
	case "2":
		session.PlayerDefend()
		if session.Player.HP <= 0 {
			session.IsOver = true
			session.Won = false
		}
	case "3":
		if exp.Rations <= 0 {
			m.Engine.SetAlert("Tidak ada ransum tersisa di dalam ransel")
		} else {
			exp.Rations--
			session.PlayerHeal(25)
			if session.Player.HP <= 0 {
				session.IsOver = true
				session.Won = false
			}
		}
	case "4":
		session.PlayerFlee()
		if session.Player.HP <= 0 {
			session.IsOver = true
			session.Won = false
		}
	}

	return m, nil
}

func (m *AppModel) updateBlacksmithCraft(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	recipes, err := data.LoadRecipeDefs()
	maxRecipes := len(recipes)
	if err != nil || maxRecipes == 0 {
		m.Engine.SwitchState(engine.StateTownMenu)
		return m, nil
	}

	switch msg.String() {
	case "esc":
		m.Engine.SwitchState(engine.StateTownMenu)
	case "up", "k":
		if m.selectedRecipeIdx > 0 {
			m.selectedRecipeIdx--
		} else {
			m.selectedRecipeIdx = maxRecipes - 1
		}
		m.Engine.ClearAlert()
	case "down", "j":
		if m.selectedRecipeIdx < maxRecipes-1 {
			m.selectedRecipeIdx++
		} else {
			m.selectedRecipeIdx = 0
		}
		m.Engine.ClearAlert()
	case "enter":
		recipe := recipes[m.selectedRecipeIdx]
		if err := m.Engine.CraftWeapon(recipe.ID); err != nil {
			m.Engine.SetAlert(err.Error())
		}
	case "r", "R":
		if err := m.Engine.RepairEquippedWeapon(); err != nil {
			m.Engine.SetAlert(err.Error())
		}
	}
	return m, nil
}

func (m *AppModel) updateTrainingGrounds(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	statNames := []string{"Might", "Agility", "Resolve", "Ingenuity"}
	maxStats := len(statNames)

	switch msg.String() {
	case "esc":
		m.Engine.SwitchState(engine.StateTownMenu)
	case "up", "k":
		if m.selectedStatIdx > 0 {
			m.selectedStatIdx--
		} else {
			m.selectedStatIdx = maxStats - 1
		}
		m.Engine.ClearAlert()
	case "down", "j":
		if m.selectedStatIdx < maxStats-1 {
			m.selectedStatIdx++
		} else {
			m.selectedStatIdx = 0
		}
		m.Engine.ClearAlert()
	case "1":
		m.selectedStatIdx = 0
		m.Engine.ClearAlert()
	case "2":
		m.selectedStatIdx = 1
		m.Engine.ClearAlert()
	case "3":
		m.selectedStatIdx = 2
		m.Engine.ClearAlert()
	case "4":
		m.selectedStatIdx = 3
		m.Engine.ClearAlert()
	case "enter":
		stat := statNames[m.selectedStatIdx]
		if err := m.Engine.TrainStat(stat); err != nil {
			m.Engine.SetAlert(err.Error())
		}
	}
	return m, nil
}

func (m *AppModel) updateExpeditionSummary(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.Engine.SwitchState(engine.StateTownMenu)
	}
	return m, nil
}

// View delegates rendering to the active screen view
func (m *AppModel) View() string {
	switch m.Engine.CurrentState {
	case engine.StateTownMenu:
		return views.RenderTownView(m.Engine, m.width)
	case engine.StateWorkerAssign:
		return views.RenderWorkerView(m.Engine, m.selectedWorkerIdx, m.width)
	case engine.StateTownBuild:
		return views.RenderBuildView(m.Engine, m.selectedBuildIdx, m.width)
	case engine.StateBlacksmithCraft:
		return views.RenderBlacksmithView(m.Engine, m.selectedRecipeIdx, m.width)
	case engine.StateTrainingGrounds:
		return views.RenderTrainingView(m.Engine, m.selectedStatIdx, m.width)
	case engine.StateDungeonExplore:
		return views.RenderDungeonView(m.Engine, m.width)
	case engine.StateCombatTurn:
		return views.RenderCombatView(m.Engine, m.width)
	case engine.StateExpeditionSummary:
		return views.RenderExpeditionSummaryView(m.Engine, m.width)
	default:
		return views.RenderTownView(m.Engine, m.width)
	}
}

