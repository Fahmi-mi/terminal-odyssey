package engine

import (
	"fmt"
	"strings"

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

// CraftWeapon crafts a new weapon from recipe and equips it on the player
func (e *Engine) CraftWeapon(recipeID string) error {
	bsLvl := e.Village.Buildings[settlement.BuildingBlacksmith]
	if bsLvl < 1 {
		return fmt.Errorf("bengkel Pandai Besi belum dibangun")
	}

	recipes, err := data.LoadRecipeDefs()
	if err != nil {
		return fmt.Errorf("gagal memuat data resep senjata: %w", err)
	}

	var targetRecipe *data.RecipeDef
	for i := range recipes {
		if recipes[i].ID == recipeID {
			targetRecipe = &recipes[i]
			break
		}
	}
	if targetRecipe == nil {
		return fmt.Errorf("resep senjata dengan ID %s tidak ditemukan", recipeID)
	}

	if bsLvl < targetRecipe.BlacksmithLevel {
		return fmt.Errorf("butuh Bengkel Pandai Besi Level %d untuk menempa %s", targetRecipe.BlacksmithLevel, targetRecipe.Name)
	}

	if e.Village.Lumber < targetRecipe.WoodCost {
		return fmt.Errorf("kayu tidak mencukupi (butuh %d, ada %d)", targetRecipe.WoodCost, e.Village.Lumber)
	}
	if e.Village.Stone < targetRecipe.StoneCost {
		return fmt.Errorf("batu tidak mencukupi (butuh %d, ada %d)", targetRecipe.StoneCost, e.Village.Stone)
	}
	if e.Village.Treasury < targetRecipe.GoldCost {
		return fmt.Errorf("emas tidak mencukupi (butuh %d, ada %d)", targetRecipe.GoldCost, e.Village.Treasury)
	}

	e.Village.Lumber -= targetRecipe.WoodCost
	e.Village.Stone -= targetRecipe.StoneCost
	e.Village.Treasury -= targetRecipe.GoldCost

	newWeapon := character.Weapon{
		ID:           targetRecipe.ID,
		Name:         targetRecipe.Name,
		WeaponType:   targetRecipe.Type,
		BaseDamage:   [2]int{targetRecipe.MinDamage, targetRecipe.MaxDamage},
		CritRate:     targetRecipe.CritRate,
		Initiative:   targetRecipe.Initiative,
		Durability:   targetRecipe.Durability,
		MaxDura:      targetRecipe.Durability,
		SpecialAffix: targetRecipe.SpecialAffix,
	}

	e.Player.EquipWeapon(newWeapon)
	e.SetAlert(fmt.Sprintf("Berhasil menempa %s (%s, %d-%d ATK)", targetRecipe.Name, targetRecipe.Type, targetRecipe.MinDamage, targetRecipe.MaxDamage))
	return nil
}

// RepairEquippedWeapon repairs the currently equipped weapon in Blacksmith
func (e *Engine) RepairEquippedWeapon() error {
	bsLvl := e.Village.Buildings[settlement.BuildingBlacksmith]
	if bsLvl < 1 {
		return fmt.Errorf("bengkel Pandai Besi belum dibangun")
	}

	w := e.Player.EquippedWeapon
	if w.Durability >= w.MaxDura {
		return fmt.Errorf("ketahanan %s masih maksimal (%d/%d)", w.Name, w.Durability, w.MaxDura)
	}

	missing := w.MaxDura - w.Durability
	goldCost := (missing * 2) / 5
	if goldCost < 5 {
		goldCost = 5
	}
	stoneCost := (missing * 1) / 5
	if stoneCost < 2 {
		stoneCost = 2
	}

	if e.Village.Treasury < goldCost {
		return fmt.Errorf("kas emas tidak cukup untuk reparasi (butuh %d Gold, ada %d)", goldCost, e.Village.Treasury)
	}
	if e.Village.Stone < stoneCost {
		return fmt.Errorf("batu tidak cukup untuk reparasi (butuh %d Batu, ada %d)", stoneCost, e.Village.Stone)
	}

	e.Village.Treasury -= goldCost
	e.Village.Stone -= stoneCost

	repaired := e.Player.RepairWeapon()
	e.SetAlert(fmt.Sprintf("Berhasil memperbaiki %s (+%d Durabilitas, -%d Gold, -%d Batu)", w.Name, repaired, goldCost, stoneCost))
	return nil
}

// TrainStat upgrades a character attribute at the Training Grounds
func (e *Engine) TrainStat(statName string) error {
	tgLvl := e.Village.Buildings[settlement.BuildingTrainingGround]
	if tgLvl < 1 {
		return fmt.Errorf("pusat Latihan belum dibangun")
	}

	cap := 10 + (tgLvl * 5)

	var currentVal int
	normalized := strings.ToLower(statName)
	switch normalized {
	case "might":
		currentVal = e.Player.Stats.Might
	case "agility":
		currentVal = e.Player.Stats.Agility
	case "resolve":
		currentVal = e.Player.Stats.Resolve
	case "ingenuity":
		currentVal = e.Player.Stats.Ingenuity
	default:
		return fmt.Errorf("nama atribut %s tidak valid", statName)
	}

	if currentVal >= cap {
		return fmt.Errorf("stat %s sudah mencapai batas maksimal Pusat Latihan Level %d (Cap: %d)", statName, tgLvl, cap)
	}

	goldCost := currentVal * 5
	rationCost := 2 + (currentVal - 10)
	if rationCost < 2 {
		rationCost = 2
	}

	if e.Village.Treasury < goldCost {
		return fmt.Errorf("kas emas tidak mencukupi (butuh %d Gold, ada %d)", goldCost, e.Village.Treasury)
	}
	if e.Village.Rations < rationCost {
		return fmt.Errorf("lumbung ransum tidak mencukupi (butuh %d Ransum, ada %d)", rationCost, e.Village.Rations)
	}

	newVal, err := e.Player.UpgradeStat(statName, cap)
	if err != nil {
		return err
	}

	e.Village.Treasury -= goldCost
	e.Village.Rations -= rationCost

	e.SetAlert(fmt.Sprintf("Latihan berhasil! %s meningkat menjadi %d (-%d Gold, -%d Ransum)", statName, newVal, goldCost, rationCost))
	return nil
}


