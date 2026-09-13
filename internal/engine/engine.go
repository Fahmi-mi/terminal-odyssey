package engine

import (
	"fmt"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/character"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

// GameState represents the current active screen or phase of the game
type GameState int

const (
	StateTownMenu GameState = iota
	StateTownBuild
	StateWorkerAssign
	StateBlacksmithCraft
	StateTrainingGrounds
	StateMarketTrade
	StateTavernRecruit
	StateDungeonExplore
	StateCombatTurn
	StateExpeditionSummary
)

// Engine is the central game manager coordinating state, settlement, player, and world simulation
type Engine struct {
	CurrentState  GameState
	PreviousState GameState

	DayCounter    int
	CurrentSeason settlement.Season
	DaysPerSeason int

	Village *settlement.Settlement
	Player  *character.Player

	DailyLogs   []string
	StatusAlert string // temporary flash message (e.g. error or success info)
}

// NewGame initializes a fresh game session
func NewGame(playerName, villageName string) *Engine {
	initialLogs := []string{
		"[+] Pemukiman darurat didirikan di lembah berkabut Oakhaven",
		"[i] Para pekerja siap menerima instruksi Anda",
		"[i] Tekan [SPACE] untuk memulai hari dan menjalankan siklus produksi",
	}

	if sc, err := data.LoadDefaultScenario(); err == nil && len(sc.InitialLogs) > 0 {
		initialLogs = sc.InitialLogs
	}

	e := &Engine{
		CurrentState:  StateTownMenu,
		PreviousState: StateTownMenu,
		DayCounter:    1,
		CurrentSeason: settlement.SeasonSpring,
		DaysPerSeason: 15, // Every 15 days = 1 season
		Village:       settlement.NewSettlement(villageName),
		Player:        character.NewDefaultPlayer(playerName),
		DailyLogs:     initialLogs,
	}
	return e
}

// PassDay advances time by 1 day, simulates village production, and advances season if threshold reached
func (e *Engine) PassDay() settlement.DailyResult {
	e.DayCounter++

	// Check season transition
	seasonIdx := ((e.DayCounter - 1) / e.DaysPerSeason) % 4
	newSeason := settlement.Season(seasonIdx)

	simResult := e.Village.SimulateDay(e.DayCounter, e.CurrentSeason)

	if newSeason != e.CurrentSeason {
		simResult.SeasonChanged = true
		simResult.OldSeason = e.CurrentSeason
		simResult.NewSeason = newSeason
		e.CurrentSeason = newSeason
		simResult.Logs = append([]string{
			fmt.Sprintf("[*] PERUBAHAN MUSIM: Memasuki %s", newSeason.String()),
		}, simResult.Logs...)
	}

	// Update daily logs with the latest results
	e.DailyLogs = simResult.Logs
	e.StatusAlert = fmt.Sprintf("Hari ke-%d telah berlalu", e.DayCounter)
	return simResult
}

// SetAlert sets a temporary feedback banner
func (e *Engine) SetAlert(msg string) {
	e.StatusAlert = msg
}

// ClearAlert clears the feedback banner
func (e *Engine) ClearAlert() {
	e.StatusAlert = ""
}

// SwitchState changes the active state and records previous state for navigation back
func (e *Engine) SwitchState(next GameState) {
	e.PreviousState = e.CurrentState
	e.CurrentState = next
	e.ClearAlert()
}
