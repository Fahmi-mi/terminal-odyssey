package save

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Fahmi-mi/terminal-odyssey/internal/alchemy"
	"github.com/Fahmi-mi/terminal-odyssey/internal/character"
	"github.com/Fahmi-mi/terminal-odyssey/internal/economy"
	"github.com/Fahmi-mi/terminal-odyssey/internal/engine"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
	"github.com/Fahmi-mi/terminal-odyssey/internal/siege"
	"github.com/Fahmi-mi/terminal-odyssey/internal/tavern"
)

const (
	DefaultSaveDir = "saves"
	Slot1          = "slot_1"
	Slot2          = "slot_2"
	Slot3          = "slot_3"
	SlotAutosave   = "autosave"
)

// SaveData is the serialized schema representing game world state
type SaveData struct {
	Version             string                       `json:"version"`
	SlotID              string                       `json:"slot_id"`
	SaveTime            string                       `json:"save_time"`
	DayCounter          int                          `json:"day_counter"`
	CurrentSeason       int                          `json:"current_season"`
	DaysPerSeason       int                          `json:"days_per_season"`
	Player              *character.Player            `json:"player"`
	Village             *settlement.Settlement       `json:"village"`
	MarketItems         []economy.MarketItem         `json:"market_items"`
	MarketEvent         string                       `json:"market_event"`
	ActiveCaravans      []*economy.CaravanExpedition `json:"active_caravans"`
	Mercenaries         []*tavern.Mercenary          `json:"mercenaries"`
	LastSiegeDay        int                          `json:"last_siege_day"`
	TotalSiegesRepelled int                          `json:"total_sieges_repelled"`
	TotalSiegesFailed   int                          `json:"total_sieges_failed"`
	LastSiegeResult     *siege.SiegeResult           `json:"last_siege_result,omitempty"`
	BossDefeated        bool                         `json:"boss_defeated"`
	HasWonGame          bool                         `json:"has_won_game"`
	VictoryAcknowledged bool                         `json:"victory_acknowledged"`
	DailyLogs           []string                     `json:"daily_logs"`
}

// SaveSlotInfo represents metadata for save slots displayed in menus
type SaveSlotInfo struct {
	SlotID       string `json:"slot_id"`
	Filename     string `json:"filename"`
	Exists       bool   `json:"exists"`
	DayCounter   int    `json:"day_counter"`
	VillageName  string `json:"village_name"`
	PlayerName   string `json:"player_name"`
	SaveTime     string `json:"save_time"`
	BossDefeated bool   `json:"boss_defeated"`
}

// resolveDir returns custom directory or fallback to DefaultSaveDir
func resolveDir(customDir []string) string {
	if len(customDir) > 0 && customDir[0] != "" {
		return customDir[0]
	}
	return DefaultSaveDir
}

// SaveGame saves the entire engine state to a slot file
func SaveGame(eng *engine.Engine, slotID string, customDir ...string) error {
	if eng == nil {
		return fmt.Errorf("engine tidak valid")
	}

	dir := resolveDir(customDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori simpan: %w", err)
	}

	filePath := filepath.Join(dir, slotID+".json")

	var marketItems []economy.MarketItem
	var marketEvent string
	if eng.Market != nil {
		marketItems = eng.Market.Items
		marketEvent = eng.Market.MarketEvent
	}

	var activeCaravans []*economy.CaravanExpedition
	if eng.Caravans != nil {
		activeCaravans = eng.Caravans.ActiveCaravans
	}

	var mercs []*tavern.Mercenary
	if eng.Tavern != nil {
		mercs = eng.Tavern.Mercenaries
	}

	var lastSiegeDay, repelled, failed int
	var lastResult *siege.SiegeResult
	if eng.Siege != nil {
		lastSiegeDay = eng.Siege.LastSiegeDay
		repelled = eng.Siege.TotalSiegesRepelled
		failed = eng.Siege.TotalSiegesFailed
		lastResult = eng.Siege.LastResult
	}

	saveData := &SaveData{
		Version:             "1.0.0",
		SlotID:              slotID,
		SaveTime:            time.Now().Format("2006-01-02 15:04:05"),
		DayCounter:          eng.DayCounter,
		CurrentSeason:       int(eng.CurrentSeason),
		DaysPerSeason:       eng.DaysPerSeason,
		Player:              eng.Player,
		Village:             eng.Village,
		MarketItems:         marketItems,
		MarketEvent:         marketEvent,
		ActiveCaravans:      activeCaravans,
		Mercenaries:         mercs,
		LastSiegeDay:        lastSiegeDay,
		TotalSiegesRepelled: repelled,
		TotalSiegesFailed:   failed,
		LastSiegeResult:     lastResult,
		BossDefeated:        eng.BossDefeated,
		HasWonGame:          eng.HasWonGame,
		VictoryAcknowledged: eng.VictoryAcknowledged,
		DailyLogs:           eng.DailyLogs,
	}

	bytes, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("gagal menyusun data simpanan: %w", err)
	}

	if err := os.WriteFile(filePath, bytes, 0644); err != nil {
		return fmt.Errorf("gagal menulis file simpanan: %w", err)
	}

	return nil
}

// LoadGame reads a saved game slot and reconstructs a functional engine
func LoadGame(slotID string, customDir ...string) (*engine.Engine, error) {
	dir := resolveDir(customDir)
	filePath := filepath.Join(dir, slotID+".json")

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("berkas simpanan tidak ditemukan: %w", err)
	}

	var data SaveData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("gagal mengurai berkas simpanan: %w", err)
	}

	// Reconstruct subsystems
	market, err := economy.NewMarketState()
	if err == nil && len(data.MarketItems) > 0 {
		market.Items = data.MarketItems
		market.MarketEvent = data.MarketEvent
	}

	caravans, err := economy.NewCaravanManager()
	if err == nil && len(data.ActiveCaravans) > 0 {
		caravans.ActiveCaravans = data.ActiveCaravans
	}

	alc, _ := alchemy.NewAlchemyManager()

	tav, err := tavern.NewTavernManager()
	if err == nil && len(data.Mercenaries) > 0 {
		tav.Mercenaries = data.Mercenaries
	}

	sm, err := siege.NewSiegeManager()
	if err == nil {
		sm.LastSiegeDay = data.LastSiegeDay
		sm.TotalSiegesRepelled = data.TotalSiegesRepelled
		sm.TotalSiegesFailed = data.TotalSiegesFailed
		sm.LastResult = data.LastSiegeResult
	}

	if data.Village != nil {
		data.Village.RecalculateDefense()
	}

	eng := &engine.Engine{
		CurrentState:        engine.StateTownMenu,
		PreviousState:       engine.StateTownMenu,
		DayCounter:          data.DayCounter,
		CurrentSeason:       settlement.Season(data.CurrentSeason),
		DaysPerSeason:       data.DaysPerSeason,
		Village:             data.Village,
		Player:              data.Player,
		Market:              market,
		Caravans:            caravans,
		Alchemy:             alc,
		Tavern:              tav,
		Siege:               sm,
		BossDefeated:        data.BossDefeated,
		HasWonGame:          data.HasWonGame,
		VictoryAcknowledged: data.VictoryAcknowledged,
		DailyLogs:           data.DailyLogs,
		StatusAlert:         fmt.Sprintf("Permainan berhasil dimuat dari %s", slotID),
	}

	return eng, nil
}

// ListSaveSlots inspects all standard save slots and retrieves summary info
func ListSaveSlots(customDir ...string) ([]SaveSlotInfo, error) {
	dir := resolveDir(customDir)
	standardSlots := []string{Slot1, Slot2, Slot3, SlotAutosave}
	result := make([]SaveSlotInfo, len(standardSlots))

	for i, slot := range standardSlots {
		filePath := filepath.Join(dir, slot+".json")
		info := SaveSlotInfo{
			SlotID:   slot,
			Filename: slot + ".json",
			Exists:   false,
		}

		if bytes, err := os.ReadFile(filePath); err == nil {
			var d SaveData
			if err := json.Unmarshal(bytes, &d); err == nil {
				info.Exists = true
				info.DayCounter = d.DayCounter
				if d.Village != nil {
					info.VillageName = d.Village.Name
				}
				if d.Player != nil {
					info.PlayerName = d.Player.Name
				}
				info.SaveTime = d.SaveTime
				info.BossDefeated = d.BossDefeated
			}
		}

		result[i] = info
	}

	return result, nil
}

// DeleteSave removes a specific save slot file
func DeleteSave(slotID string, customDir ...string) error {
	dir := resolveDir(customDir)
	filePath := filepath.Join(dir, slotID+".json")
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("gagal menghapus file simpanan: %w", err)
	}
	return nil
}
