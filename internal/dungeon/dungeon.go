package dungeon

import (
	"fmt"
	"math/rand"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/character"
	"github.com/Fahmi-mi/terminal-odyssey/internal/combat"
)

const (
	RoomTypeCombat   = "combat"
	RoomTypeTreasure = "treasure"
	RoomTypeRest     = "rest"
	RoomTypeMystery  = "mystery"
	RoomTypeExit     = "exit"
)

// Room represents an active room node in the current expedition
type Room struct {
	Index         int
	Def           data.RoomDef
	Enemy         *combat.Enemy
	IsResolved    bool
	ResolutionLog string
}

// Expedition tracks an active run through the catacombs
type Expedition struct {
	Player          *character.Player
	Floor           int
	CurrentRoomIdx  int
	Rooms           []*Room
	Torch           int
	GoldFound       int
	LumberFound     int
	StoneFound      int
	RationsFound    int
	EnemiesDefeated int
	Rations         int
	MaxBackpack     int
	IsCompleted     bool
	IsDefeated      bool
	Logs            []string
	ActiveCombat    *combat.CombatSession
}

// NewExpedition generates a structured 6-room dungeon run
func NewExpedition(p *character.Player, startingRations int) (*Expedition, error) {
	roomDefs, err := data.LoadRoomDefs()
	if err != nil {
		return nil, fmt.Errorf("gagal memuat data ruangan: %w", err)
	}

	defMap := make(map[string]data.RoomDef)
	for _, rd := range roomDefs {
		defMap[rd.ID] = rd
	}

	// Preset room sequence for catacombs level 1
	sequenceIDs := []string{
		"combat_corridor_1",
		"treasure_cache",
		"combat_spider_nest",
		"rest_shrine",
		"combat_crypt_hall",
		"exit_stairs",
	}

	rooms := make([]*Room, 0, len(sequenceIDs))
	for i, id := range sequenceIDs {
		def, ok := defMap[id]
		if !ok {
			// Fallback if specific room ID missing
			def = roomDefs[i%len(roomDefs)]
		}

		r := &Room{
			Index: i + 1,
			Def:   def,
		}

		if def.Type == RoomTypeCombat && len(def.EnemyIDs) > 0 {
			chosenEnemyID := def.EnemyIDs[rand.Intn(len(def.EnemyIDs))]
			enemy, err := combat.NewEnemyByID(chosenEnemyID)
			if err == nil {
				r.Enemy = enemy
			}
		}

		rooms = append(rooms, r)
	}

	exp := &Expedition{
		Player:         p,
		Floor:          1,
		CurrentRoomIdx: 0,
		Rooms:          rooms,
		Torch:          100,
		Rations:        startingRations,
		MaxBackpack:    p.MaxBackpack,
		Logs:           make([]string, 0),
	}

	exp.AddLog("[+] Pintu gerbang Katakombe terbuka, hawa dingin menyambut Anda")
	exp.AddLog(fmt.Sprintf("[i] Ransel memuat %d ransum perbekalan dan obor menyala penuh 100%%", startingRations))

	// Setup initial room
	firstRoom := exp.CurrentRoom()
	if firstRoom.Def.Type == RoomTypeCombat && firstRoom.Enemy != nil {
		exp.ActiveCombat = combat.NewCombatSession(p, firstRoom.Enemy)
		exp.AddLog(fmt.Sprintf("[!] Ruang 1: Musuh mendekat! Bersiap bertarung melawan %s", firstRoom.Enemy.Name))
	} else {
		exp.AddLog(fmt.Sprintf("[*] Ruang 1: %s", firstRoom.Def.Title))
	}

	return exp, nil
}

// CurrentRoom returns the room the player is currently occupying
func (exp *Expedition) CurrentRoom() *Room {
	if exp.CurrentRoomIdx < 0 {
		return exp.Rooms[0]
	}
	if exp.CurrentRoomIdx >= len(exp.Rooms) {
		return exp.Rooms[len(exp.Rooms)-1]
	}
	return exp.Rooms[exp.CurrentRoomIdx]
}

// AdvanceRoom moves the expedition to the next room
func (exp *Expedition) AdvanceRoom() error {
	if exp.IsCompleted || exp.IsDefeated {
		return fmt.Errorf("ekspedisi telah selesai")
	}

	if exp.ActiveCombat != nil && !exp.ActiveCombat.IsOver {
		return fmt.Errorf("selesaikan pertempuran terlebih dahulu")
	}

	if exp.CurrentRoomIdx >= len(exp.Rooms)-1 {
		return fmt.Errorf("sudah berada di ruangan paling ujung")
	}

	// Torch decay
	exp.Torch -= 10
	if exp.Torch < 0 {
		exp.Torch = 0
	}
	if exp.Torch <= 0 {
		exp.AddLog("[!] Obor padam total! Kegelapan malam menyulitkan langkah Anda")
	}

	exp.CurrentRoomIdx++
	room := exp.CurrentRoom()

	exp.AddLog(fmt.Sprintf("[*] Melangkah ke Ruang %d: %s", room.Index, room.Def.Title))

	// If combat room, initialize battle
	if room.Def.Type == RoomTypeCombat && room.Enemy != nil && !room.IsResolved {
		exp.ActiveCombat = combat.NewCombatSession(exp.Player, room.Enemy)
		exp.AddLog(fmt.Sprintf("[!] Musuh terdeteksi! Bersiap menghadapi %s", room.Enemy.Name))
	} else {
		exp.ActiveCombat = nil
	}

	return nil
}

// ConsumeRation consumes 1 ration to restore 25 HP
func (exp *Expedition) ConsumeRation() (int, error) {
	if exp.Rations <= 0 {
		return 0, fmt.Errorf("ransum perbekalan sudah habis")
	}
	if exp.Player.HP >= exp.Player.MaxHP {
		return 0, fmt.Errorf("HP karakter sudah maksimal")
	}

	exp.Rations--
	healAmount := 25
	oldHP := exp.Player.HP
	exp.Player.HP += healAmount
	if exp.Player.HP > exp.Player.MaxHP {
		exp.Player.HP = exp.Player.MaxHP
	}
	actualHeal := exp.Player.HP - oldHP

	exp.AddLog(fmt.Sprintf("[+] Mengonsumsi ransum kering (+%d HP, Sisa Ransum: %d)", actualHeal, exp.Rations))
	return actualHeal, nil
}

// ConsumeTorch uses an emergency torch to recover +30% torch meter
func (exp *Expedition) ConsumeTorch() error {
	if exp.Torch >= 100 {
		return fmt.Errorf("nyala obor masih maksimal (100%%)")
	}

	exp.Torch += 30
	if exp.Torch > 100 {
		exp.Torch = 100
	}
	exp.AddLog(fmt.Sprintf("[+] Menyalakan cadangan obor (Penerangan pulih ke %d%%)", exp.Torch))
	return nil
}

// ResolveTreasure loots treasure chests and materials
func (exp *Expedition) ResolveTreasure() string {
	room := exp.CurrentRoom()
	if room.Def.Type != RoomTypeTreasure || room.IsResolved {
		return "Ruangan ini sudah diperiksa"
	}

	goldGained := room.Def.MinGold
	if room.Def.MaxGold > room.Def.MinGold {
		goldGained += rand.Intn(room.Def.MaxGold - room.Def.MinGold + 1)
	}

	rationsGained := room.Def.MinRations
	if room.Def.MaxRations > room.Def.MinRations {
		rationsGained += rand.Intn(room.Def.MaxRations - room.Def.MinRations + 1)
	}

	exp.GoldFound += goldGained
	exp.Rations += rationsGained
	exp.RationsFound += rationsGained

	for _, mat := range room.Def.Materials {
		if mat.Type == "lumber" {
			exp.LumberFound += mat.Amount
		} else if mat.Type == "stone" {
			exp.StoneFound += mat.Amount
		}
	}

	room.IsResolved = true
	resText := fmt.Sprintf("Membuka peti: +%d Emas, +%d Ransum", goldGained, rationsGained)
	room.ResolutionLog = resText
	exp.AddLog(fmt.Sprintf("[+] %s", resText))
	return resText
}

// ResolveRest rests at sanctuary, healing HP and restoring torch
func (exp *Expedition) ResolveRest() string {
	room := exp.CurrentRoom()
	if room.Def.Type != RoomTypeRest || room.IsResolved {
		return "Suaka sudah pernah digunakan"
	}

	healPct := room.Def.HealPercent
	if healPct <= 0 {
		healPct = 35
	}
	healVal := (exp.Player.MaxHP * healPct) / 100
	oldHP := exp.Player.HP
	exp.Player.HP += healVal
	if exp.Player.HP > exp.Player.MaxHP {
		exp.Player.HP = exp.Player.MaxHP
	}
	actualHealed := exp.Player.HP - oldHP

	torchBonus := room.Def.TorchBonus
	if torchBonus <= 0 {
		torchBonus = 30
	}
	exp.Torch += torchBonus
	if exp.Torch > 100 {
		exp.Torch = 100
	}

	room.IsResolved = true
	resText := fmt.Sprintf("Beristirahat di suaka (+%d HP, +%d%% Obor)", actualHealed, torchBonus)
	room.ResolutionLog = resText
	exp.AddLog(fmt.Sprintf("[+] %s", resText))
	return resText
}

// ResolveMystery tests character stats against ancient altar
func (exp *Expedition) ResolveMystery() (bool, string) {
	room := exp.CurrentRoom()
	if room.Def.Type != RoomTypeMystery || room.IsResolved {
		return false, "Altar sudah diperiksa"
	}

	success := exp.Player.Stats.Ingenuity >= room.Def.ReqValue
	room.IsResolved = true

	if success {
		exp.GoldFound += room.Def.RewardGold
		resText := fmt.Sprintf("Teka-teki aksara terpecahkan! Altar membuka relik kuno (+%d Emas)", room.Def.RewardGold)
		room.ResolutionLog = resText
		exp.AddLog(fmt.Sprintf("[+] %s", resText))
		return true, resText
	}

	exp.Player.HP -= room.Def.PenaltyDamage
	if exp.Player.HP <= 0 {
		exp.Player.HP = 0
		exp.IsDefeated = true
		exp.IsCompleted = true
	}
	resText := fmt.Sprintf("Gagal memahami rune altar! Terkena jebakan gas beracun (-%d HP)", room.Def.PenaltyDamage)
	room.ResolutionLog = resText
	exp.AddLog(fmt.Sprintf("[-] %s", resText))
	return false, resText
}

// OnCombatWon handles enemy victory in current room
func (exp *Expedition) OnCombatWon() {
	room := exp.CurrentRoom()
	room.IsResolved = true
	exp.EnemiesDefeated++
	if room.Enemy != nil {
		exp.GoldFound += room.Enemy.GoldReward
	}
	exp.ActiveCombat = nil
}

// Evacuate successfully ends the expedition and brings back all loot
func (exp *Expedition) Evacuate() {
	exp.IsCompleted = true
	exp.AddLog(fmt.Sprintf("[+] Ekspedisi selesai! Berhasil membawa pulang +%d Emas ke desa", exp.GoldFound))
}

// HandleDefeat handles player defeat and strips unbanked loot
func (exp *Expedition) HandleDefeat() {
	exp.IsDefeated = true
	exp.IsCompleted = true
	exp.GoldFound = 0
	exp.LumberFound = 0
	exp.StoneFound = 0
	exp.RationsFound = 0
	exp.AddLog("[!] Ekspedisi gagal! Karakter dievakuasi darurat, seluruh jarahan hilang")
}

// AddLog appends an entry to expedition logs
func (exp *Expedition) AddLog(entry string) {
	exp.Logs = append(exp.Logs, entry)
	if len(exp.Logs) > 6 {
		exp.Logs = exp.Logs[len(exp.Logs)-6:]
	}
}
