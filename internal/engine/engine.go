package engine

import (
	"fmt"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/character"
	"github.com/Fahmi-mi/terminal-odyssey/internal/dungeon"
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

// ExpeditionSummary stores the results of a finished dungeon run
type ExpeditionSummary struct {
	WasEvacuated    bool
	GoldEarned      int
	LumberEarned    int
	StoneEarned     int
	EnemiesDefeated int
	RoomsExplored   int
	TotalRooms      int
}

// Engine is the central game manager coordinating state, settlement, player, and world simulation
type Engine struct {
	CurrentState  GameState
	PreviousState GameState

	DayCounter    int
	CurrentSeason settlement.Season
	DaysPerSeason int

	Village *settlement.Settlement
	Player  *character.Player

	ActiveExpedition      *dungeon.Expedition
	LastExpeditionSummary *ExpeditionSummary

	DailyLogs   []string
	StatusAlert string // temporary flash message (e.g. error or success info)
}

// NewGame initializes a fresh game session
func NewGame(playerName, villageName string) *Engine {
	initialLogs := []string{
		"[+] Pemukiman darurat didirikan di lembah berkabut Oakhaven",
		"[i] Para pekerja siap menerima instruksi Anda",
		"[i] Tekan [D] untuk melewati hari dan menjalankan siklus produksi",
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

	// Resting in village restores player HP if village is not starving
	if !simResult.StarvationEvent {
		e.Player.HP += 20
		if e.Player.HP > e.Player.MaxHP {
			e.Player.HP = e.Player.MaxHP
		}
	}

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

// StartExpedition initializes a new dungeon run, deducting rations from village storage
func (e *Engine) StartExpedition(rationsToTake int) error {
	if rationsToTake < 0 {
		rationsToTake = 0
	}
	if e.Village.Rations < rationsToTake {
		return fmt.Errorf("lumbung desa hanya memiliki %d ransum (butuh %d)", e.Village.Rations, rationsToTake)
	}

	e.Village.Rations -= rationsToTake

	exp, err := dungeon.NewExpedition(e.Player, rationsToTake)
	if err != nil {
		e.Village.Rations += rationsToTake
		return err
	}

	e.ActiveExpedition = exp
	e.SwitchState(StateDungeonExplore)
	return nil
}

// FinishExpedition concludes the active dungeon run and transfers loot to village
func (e *Engine) FinishExpedition(evacuated bool) {
	if e.ActiveExpedition == nil {
		return
	}

	exp := e.ActiveExpedition
	wasEvac := evacuated && !exp.IsDefeated

	totalRooms := exp.TotalDepths
	if totalRooms == 0 {
		totalRooms = len(exp.Rooms)
	}
	roomsExplored := exp.RoomsExploredCount
	if roomsExplored == 0 {
		roomsExplored = exp.CurrentRoomIdx + 1
	}

	summary := &ExpeditionSummary{
		WasEvacuated:    wasEvac,
		GoldEarned:      exp.GoldFound,
		LumberEarned:    exp.LumberFound,
		StoneEarned:     exp.StoneFound,
		EnemiesDefeated: exp.EnemiesDefeated,
		RoomsExplored:   roomsExplored,
		TotalRooms:      totalRooms,
	}

	if wasEvac {
		exp.Evacuate()
		e.Village.Treasury += exp.GoldFound
		e.Village.Lumber += exp.LumberFound
		e.Village.Stone += exp.StoneFound
		// Return leftover rations to village store
		e.Village.Rations += exp.Rations

		maxL, maxS, maxR := e.Village.StorageCap()
		if e.Village.Lumber > maxL {
			e.Village.Lumber = maxL
		}
		if e.Village.Stone > maxS {
			e.Village.Stone = maxS
		}
		if e.Village.Rations > maxR {
			e.Village.Rations = maxR
		}

		e.DailyLogs = append([]string{
			fmt.Sprintf("[+] Ekspedisi sukses: membawa pulang +%d Emas, +%d Kayu, +%d Batu", exp.GoldFound, exp.LumberFound, exp.StoneFound),
		}, e.DailyLogs...)
	} else {
		exp.HandleDefeat()
		// 1 day passes for medical recovery
		e.PassDay()
		// Player wakes up in convalescence with 25 HP
		e.Player.HP = 25
		e.DailyLogs = append([]string{
			"[!] Ekspedisi gagal: Karakter dievakuasi darurat ke desa dan seluruh jarahan hilang",
		}, e.DailyLogs...)
	}

	e.LastExpeditionSummary = summary
	e.SwitchState(StateExpeditionSummary)
}

