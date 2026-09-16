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
	blacksmithTab     int
	selectedWeaponIdx int

	selectedMarketIdx int
	marketTab         int
	selectedRouteIdx  int
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
		blacksmithTab:     0,
		selectedWeaponIdx: 0,
		selectedMarketIdx: 0,
		marketTab:         0,
		selectedRouteIdx:  0,
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
		case engine.StateMarketTrade:
			return m.updateMarketTrade(msg)
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
	case "1":
		m.Engine.SwitchState(engine.StateMarketTrade)
	case "2":
		m.Engine.SwitchState(engine.StateTownBuild)
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
	case "up":
		if m.selectedWorkerIdx > 0 {
			m.selectedWorkerIdx--
		} else {
			m.selectedWorkerIdx = maxRoles - 1
		}
		m.Engine.ClearAlert()
	case "down":
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
	case "up":
		if m.selectedBuildIdx > 0 {
			m.selectedBuildIdx--
		} else {
			m.selectedBuildIdx = maxBuildings - 1
		}
		m.Engine.ClearAlert()
	case "down":
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
	choices := exp.NextRoomChoices()

	if room.Def.Type == dungeon.RoomTypeExit {
		switch msg.String() {
		case "enter", "esc":
			m.Engine.FinishExpedition(true)
		case "m", "M":
			if exp.Rations > 0 && exp.Player.HP < exp.Player.MaxHP {
				if _, err := exp.ConsumeRation(); err != nil {
					m.Engine.SetAlert(err.Error())
				}
			}
		}
		return m, nil
	}

	if room.Def.Type == dungeon.RoomTypeCombat && !room.IsResolved {
		switch msg.String() {
		case "enter":
			if exp.ActiveCombat == nil && room.Enemy != nil {
				exp.ActiveCombat = combat.NewCombatSession(exp.Player, room.Enemy)
			}
			m.Engine.SwitchState(engine.StateCombatTurn)
		case "m", "M":
			if _, err := exp.ConsumeRation(); err != nil {
				m.Engine.SetAlert(err.Error())
			}
		case "esc":
			m.Engine.FinishExpedition(true)
		}
		return m, nil
	}

	if !room.IsResolved {
		switch msg.String() {
		case "enter":
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
		case "m", "M":
			if _, err := exp.ConsumeRation(); err != nil {
				m.Engine.SetAlert(err.Error())
			}
		case "o", "O":
			if err := exp.ConsumeTorch(); err != nil {
				m.Engine.SetAlert(err.Error())
			}
		case "esc":
			m.Engine.FinishExpedition(true)
		}
		return m, nil
	}

	// Room is resolved
	if len(choices) >= 2 {
		switch msg.String() {
		case "1":
			if err := exp.AdvanceToRoom(choices[0].GraphIdx); err != nil {
				m.Engine.SetAlert(err.Error())
			}
		case "2":
			if err := exp.AdvanceToRoom(choices[1].GraphIdx); err != nil {
				m.Engine.SetAlert(err.Error())
			}
		case "m", "M":
			if _, err := exp.ConsumeRation(); err != nil {
				m.Engine.SetAlert(err.Error())
			}
		case "o", "O":
			if err := exp.ConsumeTorch(); err != nil {
				m.Engine.SetAlert(err.Error())
			}
		case "esc":
			m.Engine.FinishExpedition(true)
		}
		return m, nil
	}

	// Resolved with single choice
	switch msg.String() {
	case "enter":
		if err := exp.AdvanceRoom(); err != nil {
			m.Engine.SetAlert(err.Error())
		}
	case "m", "M":
		if _, err := exp.ConsumeRation(); err != nil {
			m.Engine.SetAlert(err.Error())
		}
	case "o", "O":
		if err := exp.ConsumeTorch(); err != nil {
			m.Engine.SetAlert(err.Error())
		}
	case "esc":
		m.Engine.FinishExpedition(true)
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
	recipes, _ := data.LoadRecipeDefs()
	maxRecipes := len(recipes)
	owned := m.Engine.Player.OwnedWeapons
	maxOwned := len(owned)

	switch msg.String() {
	case "esc":
		m.Engine.SwitchState(engine.StateTownMenu)
	case "tab":
		if m.blacksmithTab == 0 {
			m.blacksmithTab = 1
		} else {
			m.blacksmithTab = 0
		}
		m.Engine.ClearAlert()
	case "1":
		m.blacksmithTab = 0
		m.Engine.ClearAlert()
	case "2":
		m.blacksmithTab = 1
		m.Engine.ClearAlert()
	case "up":
		if m.blacksmithTab == 0 {
			if m.selectedRecipeIdx > 0 {
				m.selectedRecipeIdx--
			} else if maxRecipes > 0 {
				m.selectedRecipeIdx = maxRecipes - 1
			}
		} else {
			if m.selectedWeaponIdx > 0 {
				m.selectedWeaponIdx--
			} else if maxOwned > 0 {
				m.selectedWeaponIdx = maxOwned - 1
			}
		}
		m.Engine.ClearAlert()
	case "down":
		if m.blacksmithTab == 0 {
			if m.selectedRecipeIdx < maxRecipes-1 {
				m.selectedRecipeIdx++
			} else {
				m.selectedRecipeIdx = 0
			}
		} else {
			if m.selectedWeaponIdx < maxOwned-1 {
				m.selectedWeaponIdx++
			} else {
				m.selectedWeaponIdx = 0
			}
		}
		m.Engine.ClearAlert()
	case "enter":
		if m.blacksmithTab == 0 {
			if maxRecipes > 0 && m.selectedRecipeIdx < maxRecipes {
				recipe := recipes[m.selectedRecipeIdx]
				if m.Engine.Player.OwnsWeapon(recipe.ID) {
					if err := m.Engine.SwitchWeapon(recipe.ID); err != nil {
						m.Engine.SetAlert(err.Error())
					}
				} else {
					if err := m.Engine.CraftWeapon(recipe.ID); err != nil {
						m.Engine.SetAlert(err.Error())
					}
				}
			}
		} else {
			if maxOwned > 0 && m.selectedWeaponIdx < maxOwned {
				target := m.Engine.Player.OwnedWeapons[m.selectedWeaponIdx]
				if err := m.Engine.SwitchWeapon(target.ID); err != nil {
					m.Engine.SetAlert(err.Error())
				}
			}
		}
	case "r", "R":
		if m.blacksmithTab == 0 {
			if err := m.Engine.RepairEquippedWeapon(); err != nil {
				m.Engine.SetAlert(err.Error())
			}
		} else {
			if maxOwned > 0 && m.selectedWeaponIdx < maxOwned {
				target := m.Engine.Player.OwnedWeapons[m.selectedWeaponIdx]
				if err := m.Engine.RepairWeaponByID(target.ID); err != nil {
					m.Engine.SetAlert(err.Error())
				}
			}
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
	case "up":
		if m.selectedStatIdx > 0 {
			m.selectedStatIdx--
		} else {
			m.selectedStatIdx = maxStats - 1
		}
		m.Engine.ClearAlert()
	case "down":
		if m.selectedStatIdx < maxStats-1 {
			m.selectedStatIdx++
		} else {
			m.selectedStatIdx = 0
		}
		m.Engine.ClearAlert()
	case "enter":
		stat := statNames[m.selectedStatIdx]
		if err := m.Engine.TrainStat(stat); err != nil {
			m.Engine.SetAlert(err.Error())
		}
	}
	return m, nil
}

func (m *AppModel) updateMarketTrade(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.Engine.Market == nil {
		m.Engine.SwitchState(engine.StateTownMenu)
		return m, nil
	}

	postLvl := m.Engine.Village.Buildings[settlement.BuildingCaravanPost]

	if m.marketTab == 0 {
		maxItems := len(m.Engine.Market.Items)
		switch msg.String() {
		case "esc":
			m.Engine.SwitchState(engine.StateTownMenu)
		case "tab", "2":
			m.marketTab = 1
			m.Engine.ClearAlert()
		case "1":
			m.marketTab = 0
			m.Engine.ClearAlert()
		case "up":
			if m.selectedMarketIdx > 0 {
				m.selectedMarketIdx--
			} else if maxItems > 0 {
				m.selectedMarketIdx = maxItems - 1
			}
			m.Engine.ClearAlert()
		case "down":
			if m.selectedMarketIdx < maxItems-1 {
				m.selectedMarketIdx++
			} else {
				m.selectedMarketIdx = 0
			}
			m.Engine.ClearAlert()
		case "b", "B":
			if maxItems > 0 && m.selectedMarketIdx < maxItems {
				item := m.Engine.Market.Items[m.selectedMarketIdx]
				if err := m.Engine.BuyCommodity(item.Def.ID, 1); err != nil {
					m.Engine.SetAlert(err.Error())
				}
			}
		case "s", "S":
			if maxItems > 0 && m.selectedMarketIdx < maxItems {
				item := m.Engine.Market.Items[m.selectedMarketIdx]
				if err := m.Engine.SellCommodity(item.Def.ID, 1); err != nil {
					m.Engine.SetAlert(err.Error())
				}
			}
		}
	} else {
		if postLvl < 1 {
			switch msg.String() {
			case "esc":
				m.Engine.SwitchState(engine.StateTownMenu)
			case "tab", "1":
				m.marketTab = 0
				m.Engine.ClearAlert()
			case "2":
				m.marketTab = 1
				m.Engine.ClearAlert()
			}
			return m, nil
		}

		maxRoutes := len(m.Engine.Caravans.Routes)
		switch msg.String() {
		case "esc":
			m.Engine.SwitchState(engine.StateTownMenu)
		case "tab", "1":
			m.marketTab = 0
			m.Engine.ClearAlert()
		case "2":
			m.marketTab = 1
			m.Engine.ClearAlert()
		case "up":
			if m.selectedRouteIdx > 0 {
				m.selectedRouteIdx--
			} else if maxRoutes > 0 {
				m.selectedRouteIdx = maxRoutes - 1
			}
			m.Engine.ClearAlert()
		case "down":
			if m.selectedRouteIdx < maxRoutes-1 {
				m.selectedRouteIdx++
			} else {
				m.selectedRouteIdx = 0
			}
			m.Engine.ClearAlert()
		case "enter":
			if maxRoutes > 0 && m.selectedRouteIdx < maxRoutes {
				route := m.Engine.Caravans.Routes[m.selectedRouteIdx]
				if err := m.Engine.DispatchCaravan(route.ID); err != nil {
					m.Engine.SetAlert(err.Error())
				}
			}
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
		return views.RenderBlacksmithView(m.Engine, m.selectedRecipeIdx, m.blacksmithTab, m.selectedWeaponIdx, m.width)
	case engine.StateTrainingGrounds:
		return views.RenderTrainingView(m.Engine, m.selectedStatIdx, m.width)
	case engine.StateMarketTrade:
		return views.RenderMarketView(m.Engine, m.selectedMarketIdx, m.marketTab, m.selectedRouteIdx, m.width)
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


