package ui

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/combat"
	"github.com/Fahmi-mi/terminal-odyssey/internal/dungeon"
	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/save"
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

	selectedAlchemyIdx int
	selectedTavernIdx  int
	tavernTab          int

	selectedTitleIdx    int
	titleMode           views.TitleMenuMode
	saveSlots           []save.SaveSlotInfo
	selectedSaveSlotIdx int
}

// NewAppModel creates a fresh TUI model
func NewAppModel(eng *engine.Engine) *AppModel {
	slots, _ := save.ListSaveSlots()

	return &AppModel{
		Engine:              eng,
		width:               80,
		height:              24,
		selectedWorkerIdx:   0,
		selectedBuildIdx:    0,
		selectedRecipeIdx:   0,
		selectedStatIdx:     0,
		blacksmithTab:       0,
		selectedWeaponIdx:   0,
		selectedMarketIdx:   0,
		marketTab:           0,
		selectedRouteIdx:    0,
		selectedAlchemyIdx:  0,
		selectedTavernIdx:   0,
		tavernTab:           0,
		selectedTitleIdx:    0,
		titleMode:           views.TitleModeMain,
		saveSlots:           slots,
		selectedSaveSlotIdx: 0,
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
		case engine.StateTitleScreen:
			return m.updateTitleScreen(msg)
		case engine.StateSaveMenu:
			return m.updateSaveMenu(msg)
		case engine.StateSiegeReport:
			return m.updateSiegeReport(msg)
		case engine.StateVictoryScreen:
			return m.updateVictoryScreen(msg)
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
		case engine.StateAlchemyLab:
			return m.updateAlchemyLab(msg)
		case engine.StateTavernRecruit:
			return m.updateTavernRecruit(msg)
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
		_ = save.SaveGame(m.Engine, save.SlotAutosave)
	case "s", "S":
		slots, _ := save.ListSaveSlots()
		m.saveSlots = slots
		m.selectedSaveSlotIdx = 0
		m.Engine.SwitchState(engine.StateSaveMenu)
	case "m", "M":
		slots, _ := save.ListSaveSlots()
		m.saveSlots = slots
		m.titleMode = views.TitleModeMain
		m.selectedTitleIdx = 0
		m.Engine.SwitchState(engine.StateTitleScreen)
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
		if m.Engine.Village.Buildings[settlement.BuildingTavern] < 1 {
			m.Engine.SetAlert("Kedai Minum belum didirikan! Bangun di menu [2] Pembangunan")
		} else {
			m.Engine.SwitchState(engine.StateTavernRecruit)
		}
	case "6":
		rationsToTake := 3
		if m.Engine.Village.Rations < rationsToTake {
			rationsToTake = m.Engine.Village.Rations
		}
		if err := m.Engine.StartExpedition(rationsToTake); err != nil {
			m.Engine.SetAlert(err.Error())
		}
	case "7":
		if m.Engine.Village.Buildings[settlement.BuildingApothecary] < 1 {
			m.Engine.SetAlert("Laboratorium Alkimia belum didirikan! Bangun di menu [2] Pembangunan")
		} else {
			m.Engine.SwitchState(engine.StateAlchemyLab)
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
		case "p", "P":
			m.usePotionInDungeon()
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
		case "p", "P":
			m.usePotionInDungeon()
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
		case "p", "P":
			m.usePotionInDungeon()
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
		case "p", "P":
			m.usePotionInDungeon()
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
	case "p", "P":
		m.usePotionInDungeon()
	case "esc":
		m.Engine.FinishExpedition(true)
	}
	return m, nil
}

func (m *AppModel) usePotionInDungeon() {
	exp := m.Engine.ActiveExpedition
	if exp == nil {
		return
	}
	p := m.Engine.Player
	if p.TotalPotions() <= 0 {
		m.Engine.SetAlert("Kantong ramuan Anda kosong")
		return
	}

	var targetPotion string
	if p.HP < p.MaxHP && p.GetPotionCount("salep_pemulih") > 0 {
		targetPotion = "salep_pemulih"
	} else if exp.Torch < 65 && p.GetPotionCount("minyak_obor") > 0 {
		targetPotion = "minyak_obor"
	} else if p.Sanity < p.MaxSanity && p.GetPotionCount("tonik_penenang") > 0 {
		targetPotion = "tonik_penenang"
	} else if p.GetPotionCount("penawar_racun") > 0 {
		targetPotion = "penawar_racun"
	} else {
		order := []string{"salep_pemulih", "minyak_obor", "tonik_penenang", "penawar_racun", "eliksir_kekuatan"}
		for _, pot := range order {
			if p.GetPotionCount(pot) > 0 {
				targetPotion = pot
				break
			}
		}
	}

	if targetPotion == "" {
		m.Engine.SetAlert("Tidak ada ramuan yang dapat digunakan saat ini")
		return
	}

	if _, err := exp.ConsumePotion(targetPotion); err != nil {
		m.Engine.SetAlert(err.Error())
	}
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
	case "p", "P":
		p := m.Engine.Player
		if p.TotalPotions() <= 0 {
			m.Engine.SetAlert("Kantong ramuan Anda kosong")
		} else {
			var targetPotion string
			if p.HP < p.MaxHP && p.GetPotionCount("salep_pemulih") > 0 {
				targetPotion = "salep_pemulih"
			} else if p.GetPotionCount("eliksir_kekuatan") > 0 && session.TemporaryAtkBuff == 0 {
				targetPotion = "eliksir_kekuatan"
			} else if p.Sanity < p.MaxSanity && p.GetPotionCount("tonik_penenang") > 0 {
				targetPotion = "tonik_penenang"
			} else if p.GetPotionCount("penawar_racun") > 0 {
				targetPotion = "penawar_racun"
			} else {
				order := []string{"salep_pemulih", "eliksir_kekuatan", "tonik_penenang", "penawar_racun"}
				for _, pot := range order {
					if p.GetPotionCount(pot) > 0 {
						targetPotion = pot
						break
					}
				}
			}

			if targetPotion == "" {
				m.Engine.SetAlert("Tidak ada ramuan tempur yang dapat diminum saat ini")
			} else {
				if _, err := session.PlayerDrinkPotion(targetPotion); err != nil {
					m.Engine.SetAlert(err.Error())
				}
				if session.Player.HP <= 0 {
					session.IsOver = true
					session.Won = false
				}
			}
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

func (m *AppModel) updateAlchemyLab(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.Engine.Alchemy == nil {
		m.Engine.SwitchState(engine.StateTownMenu)
		return m, nil
	}

	recipes := m.Engine.Alchemy.Recipes
	maxRecipes := len(recipes)

	switch msg.String() {
	case "esc":
		m.Engine.SwitchState(engine.StateTownMenu)
	case "up":
		if m.selectedAlchemyIdx > 0 {
			m.selectedAlchemyIdx--
		} else if maxRecipes > 0 {
			m.selectedAlchemyIdx = maxRecipes - 1
		}
		m.Engine.ClearAlert()
	case "down":
		if m.selectedAlchemyIdx < maxRecipes-1 {
			m.selectedAlchemyIdx++
		} else {
			m.selectedAlchemyIdx = 0
		}
		m.Engine.ClearAlert()
	case "enter":
		if maxRecipes > 0 && m.selectedAlchemyIdx < maxRecipes {
			sel := recipes[m.selectedAlchemyIdx]
			if err := m.Engine.BrewPotion(sel.ID); err != nil {
				m.Engine.SetAlert(err.Error())
			}
		}
	}
	return m, nil
}

func (m *AppModel) updateTavernRecruit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.Engine.Tavern == nil {
		m.Engine.SwitchState(engine.StateTownMenu)
		return m, nil
	}

	mercs := m.Engine.Tavern.Mercenaries
	maxMercs := len(mercs)

	switch msg.String() {
	case "esc":
		m.Engine.SwitchState(engine.StateTownMenu)
	case "tab":
		m.tavernTab = 1 - m.tavernTab
		m.Engine.ClearAlert()
	case "1":
		m.tavernTab = 0
		m.Engine.ClearAlert()
	case "2":
		m.tavernTab = 1
		m.Engine.ClearAlert()
	}

	if m.tavernTab == 0 {
		switch msg.String() {
		case "up":
			if m.selectedTavernIdx > 0 {
				m.selectedTavernIdx--
			} else if maxMercs > 0 {
				m.selectedTavernIdx = maxMercs - 1
			}
			m.Engine.ClearAlert()
		case "down":
			if m.selectedTavernIdx < maxMercs-1 {
				m.selectedTavernIdx++
			} else {
				m.selectedTavernIdx = 0
			}
			m.Engine.ClearAlert()
		case "enter":
			if maxMercs > 0 && m.selectedTavernIdx < maxMercs {
				sel := mercs[m.selectedTavernIdx]
				if sel.IsHired {
					if err := m.Engine.DismissCompanion(sel.Def.ID); err != nil {
						m.Engine.SetAlert(err.Error())
					}
				} else {
					if err := m.Engine.HireCompanion(sel.Def.ID); err != nil {
						m.Engine.SetAlert(err.Error())
					}
				}
			}
		}
	} else {
		switch msg.String() {
		case "enter", "r", "R":
			if err := m.Engine.TavernRest(); err != nil {
				m.Engine.SetAlert(err.Error())
			}
		case "m", "M":
			m.Engine.TavernRumor()
		}
	}

	return m, nil
}

func (m *AppModel) updateTitleScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.titleMode == views.TitleModeMain {
		switch msg.String() {
		case "up", "k":
			if m.selectedTitleIdx > 0 {
				m.selectedTitleIdx--
			} else {
				m.selectedTitleIdx = 2
			}
			m.Engine.ClearAlert()
		case "down", "j":
			if m.selectedTitleIdx < 2 {
				m.selectedTitleIdx++
			} else {
				m.selectedTitleIdx = 0
			}
			m.Engine.ClearAlert()
		case "1":
			m.Engine.SwitchState(engine.StateTownMenu)
		case "2":
			m.titleMode = views.TitleModeSelectSlot
			slots, _ := save.ListSaveSlots()
			m.saveSlots = slots
			m.selectedSaveSlotIdx = 0
			m.Engine.ClearAlert()
		case "3", "q", "Q":
			return m, tea.Quit
		case "enter":
			switch m.selectedTitleIdx {
			case 0:
				m.Engine.SwitchState(engine.StateTownMenu)
			case 1:
				m.titleMode = views.TitleModeSelectSlot
				slots, _ := save.ListSaveSlots()
				m.saveSlots = slots
				m.selectedSaveSlotIdx = 0
				m.Engine.ClearAlert()
			case 2:
				return m, tea.Quit
			}
		}
	} else {
		maxSlots := len(m.saveSlots)
		switch msg.String() {
		case "esc", "b", "B", "q", "Q":
			m.titleMode = views.TitleModeMain
			m.selectedTitleIdx = 1
			m.Engine.ClearAlert()
		case "up", "k":
			if m.selectedSaveSlotIdx > 0 {
				m.selectedSaveSlotIdx--
			} else {
				m.selectedSaveSlotIdx = maxSlots - 1
			}
			m.Engine.ClearAlert()
		case "down", "j":
			if m.selectedSaveSlotIdx < maxSlots-1 {
				m.selectedSaveSlotIdx++
			} else {
				m.selectedSaveSlotIdx = 0
			}
			m.Engine.ClearAlert()
		case "1", "2", "3", "4":
			if len(msg.String()) == 1 {
				idx := int(msg.String()[0] - '1')
				if idx >= 0 && idx < maxSlots {
					m.selectedSaveSlotIdx = idx
					m.Engine.ClearAlert()
				}
			}
		case "enter":
			if m.selectedSaveSlotIdx < maxSlots {
				sel := m.saveSlots[m.selectedSaveSlotIdx]
				if sel.Exists {
					loadedEng, err := save.LoadGame(sel.SlotID)
					if err != nil {
						m.Engine.SetAlert(err.Error())
					} else {
						m.Engine = loadedEng
						m.titleMode = views.TitleModeMain
						m.Engine.SwitchState(engine.StateTownMenu)
					}
				} else {
					m.Engine.SetAlert(fmt.Sprintf("Slot %d masih kosong, belum ada data simpanan", m.selectedSaveSlotIdx+1))
				}
			}
		}
	}
	return m, nil
}

func (m *AppModel) updateSaveMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	manualSlots := []string{save.Slot1, save.Slot2, save.Slot3}
	switch msg.String() {
	case "esc", "b", "B", "q", "Q":
		m.Engine.SwitchState(engine.StateTownMenu)
	case "up", "k":
		if m.selectedSaveSlotIdx > 0 {
			m.selectedSaveSlotIdx--
		} else {
			m.selectedSaveSlotIdx = len(manualSlots) - 1
		}
		m.Engine.ClearAlert()
	case "down", "j":
		if m.selectedSaveSlotIdx < len(manualSlots)-1 {
			m.selectedSaveSlotIdx++
		} else {
			m.selectedSaveSlotIdx = 0
		}
		m.Engine.ClearAlert()
	case "1", "2", "3":
		if len(msg.String()) == 1 {
			idx := int(msg.String()[0] - '1')
			if idx >= 0 && idx < len(manualSlots) {
				m.selectedSaveSlotIdx = idx
				m.Engine.ClearAlert()
			}
		}
	case "enter":
		if m.selectedSaveSlotIdx < len(manualSlots) {
			slotID := manualSlots[m.selectedSaveSlotIdx]
			if err := save.SaveGame(m.Engine, slotID); err != nil {
				m.Engine.SetAlert(err.Error())
			} else {
				m.Engine.SetAlert(fmt.Sprintf("[+] Permainan berhasil disimpan di %s", slotID))
				slots, _ := save.ListSaveSlots()
				m.saveSlots = slots
				m.Engine.SwitchState(engine.StateTownMenu)
			}
		}
	}
	return m, nil
}

func (m *AppModel) updateSiegeReport(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ", "esc":
		if m.Engine.CheckVictoryCondition() {
			m.Engine.SwitchState(engine.StateVictoryScreen)
		} else {
			m.Engine.SwitchState(engine.StateTownMenu)
		}
	}
	return m, nil
}

func (m *AppModel) updateVictoryScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		m.Engine.AcknowledgeVictory()
	case "s", "S":
		slots, _ := save.ListSaveSlots()
		m.saveSlots = slots
		m.selectedSaveSlotIdx = 0
		m.Engine.SwitchState(engine.StateSaveMenu)
	case "q", "Q", "esc":
		m.titleMode = views.TitleModeMain
		m.selectedTitleIdx = 0
		m.Engine.SwitchState(engine.StateTitleScreen)
	}
	return m, nil
}

func (m *AppModel) updateExpeditionSummary(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		_ = save.SaveGame(m.Engine, save.SlotAutosave)
		if m.Engine.CheckVictoryCondition() {
			m.Engine.SwitchState(engine.StateVictoryScreen)
		} else {
			m.Engine.SwitchState(engine.StateTownMenu)
		}
	}
	return m, nil
}

// View delegates rendering to the active screen view
func (m *AppModel) View() string {
	switch m.Engine.CurrentState {
	case engine.StateTitleScreen:
		selIdx := m.selectedTitleIdx
		if m.titleMode == views.TitleModeSelectSlot {
			selIdx = m.selectedSaveSlotIdx
		}
		return views.RenderTitleView(m.titleMode, selIdx, m.saveSlots, m.Engine.StatusAlert, m.width)
	case engine.StateSaveMenu:
		return views.RenderSaveMenuView(m.Engine, m.selectedSaveSlotIdx, m.saveSlots, m.width)
	case engine.StateSiegeReport:
		return views.RenderSiegeView(m.Engine, m.width)
	case engine.StateVictoryScreen:
		return views.RenderVictoryView(m.Engine, m.width)
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
	case engine.StateAlchemyLab:
		return views.RenderAlchemyView(m.Engine, m.selectedAlchemyIdx, m.width)
	case engine.StateTavernRecruit:
		return views.RenderTavernView(m.Engine, m.selectedTavernIdx, m.tavernTab, m.width)
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


