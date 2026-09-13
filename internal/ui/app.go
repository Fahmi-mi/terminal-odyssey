package ui

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"

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
}

// NewAppModel creates a fresh TUI model
func NewAppModel(eng *engine.Engine) *AppModel {
	return &AppModel{
		Engine:            eng,
		width:             80,
		height:            24,
		selectedWorkerIdx: 0,
		selectedBuildIdx:  0,
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
	case "d", "D", " ":
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
			m.Engine.SetAlert("Bengkel Pandai Besi siap dibuka (Milestone 4)")
		}
	case "4":
		m.Engine.SetAlert("Pusat Latihan Stat Karakter sedang dipersiapkan")
	case "5":
		m.Engine.SetAlert("Kedai Minum belum memiliki rumor baru hari ini")
	case "6":
		m.Engine.SetAlert("Pintu Katakombe Bawah Tanah sedang dipersiapkan untuk Ekspedisi (Milestone 2)")
	}
	return m, nil
}

func (m *AppModel) updateWorkerAssign(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	maxRoles := len(views.RoleItems)

	switch msg.String() {
	case "esc", "q", "b":
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
	case "right", "+", "l", "enter":
		role := views.RoleItems[m.selectedWorkerIdx].Role
		if err := m.Engine.Village.AssignWorker(role); err != nil {
			m.Engine.SetAlert(err.Error())
		} else {
			m.Engine.SetAlert(fmt.Sprintf("Berhasil menugaskan +1 %s", role))
		}
	case "left", "-", "h":
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
	case "esc", "q", "b":
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
			bName := buildings[m.selectedBuildIdx]
			if err := m.Engine.Village.UpgradeBuilding(bName); err != nil {
				m.Engine.SetAlert(err.Error())
			} else {
				m.Engine.SetAlert(fmt.Sprintf("[+] %s berhasil ditingkatkan ke Level %d", bName, m.Engine.Village.Buildings[bName]))
			}
		}
	case "enter", "u", "U":
		bName := buildings[m.selectedBuildIdx]
		if err := m.Engine.Village.UpgradeBuilding(bName); err != nil {
			m.Engine.SetAlert(err.Error())
		} else {
			m.Engine.SetAlert(fmt.Sprintf("[+] %s berhasil ditingkatkan ke Level %d", bName, m.Engine.Village.Buildings[bName]))
		}
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
	default:
		return views.RenderTownView(m.Engine, m.width)
	}
}
