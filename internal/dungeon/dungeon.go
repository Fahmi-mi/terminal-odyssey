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

// Room represents an active room node in the expedition DAG
type Room struct {
	Index           int          `json:"index"`
	GraphIdx        int          `json:"graph_idx"`
	Depth           int          `json:"depth"`
	BranchName      string       `json:"branch_name"`
	NextRoomIndices []int        `json:"next_room_indices"`
	Def             data.RoomDef `json:"def"`
	Enemy           *combat.Enemy `json:"enemy,omitempty"`
	IsResolved      bool         `json:"is_resolved"`
	ResolutionLog   string       `json:"resolution_log"`
}

// Expedition tracks an active run through the catacombs
type Expedition struct {
	Player             *character.Player
	Floor              int
	CurrentRoomIdx     int
	Rooms              []*Room
	TotalDepths        int
	RoomsExploredCount int
	Torch              int
	GoldFound          int
	LumberFound        int
	StoneFound         int
	RationsFound       int
	EnemiesDefeated    int
	Rations            int
	MaxBackpack        int
	IsCompleted        bool
	IsDefeated         bool
	Logs               []string
	ActiveCombat       *combat.CombatSession
}

type weightedChoice struct {
	roomType string
	weight   int
}

func pickWeightedType(choices []weightedChoice) string {
	totalWeight := 0
	for _, c := range choices {
		if c.weight > 0 {
			totalWeight += c.weight
		}
	}
	if totalWeight <= 0 {
		return RoomTypeCombat
	}
	roll := rand.Intn(totalWeight)
	accum := 0
	for _, c := range choices {
		if c.weight <= 0 {
			continue
		}
		accum += c.weight
		if roll < accum {
			return c.roomType
		}
	}
	return choices[0].roomType
}

func pickDefByType(roomType string, depth int, defMap map[string]data.RoomDef, usedIDs map[string]bool) data.RoomDef {
	switch roomType {
	case RoomTypeCombat:
		var candidates []string
		switch depth {
		case 1:
			candidates = []string{"combat_corridor_1"}
		case 2:
			candidates = []string{"combat_spider_nest", "combat_corridor_1", "combat_goblin_ambush"}
		case 3:
			candidates = []string{"combat_goblin_ambush", "combat_spider_nest", "combat_crypt_hall"}
		case 4:
			candidates = []string{"combat_crypt_hall", "combat_goblin_ambush"}
		case 5:
			candidates = []string{"combat_boss_chamber"}
		default:
			candidates = []string{"combat_corridor_1"}
		}
		var unused []string
		for _, c := range candidates {
			if !usedIDs[c] {
				unused = append(unused, c)
			}
		}
		var chosenID string
		if len(unused) > 0 {
			chosenID = unused[rand.Intn(len(unused))]
		} else {
			chosenID = candidates[rand.Intn(len(candidates))]
		}
		usedIDs[chosenID] = true
		if def, ok := defMap[chosenID]; ok {
			return def
		}
	case RoomTypeTreasure:
		if def, ok := defMap["treasure_cache"]; ok {
			return def
		}
	case RoomTypeRest:
		if def, ok := defMap["rest_shrine"]; ok {
			return def
		}
	case RoomTypeMystery:
		if def, ok := defMap["mystery_altar"]; ok {
			return def
		}
	case RoomTypeExit:
		if def, ok := defMap["exit_stairs"]; ok {
			return def
		}
	}
	for _, def := range defMap {
		return def
	}
	return data.RoomDef{Type: roomType, Title: "Lorong Bawah Tanah"}
}

func createRoomNode(graphIdx int, depth int, branchName string, nextIndices []int, def data.RoomDef) *Room {
	r := &Room{
		Index:           depth,
		GraphIdx:        graphIdx,
		Depth:           depth,
		BranchName:      branchName,
		NextRoomIndices: nextIndices,
		Def:             def,
	}

	if def.Type == RoomTypeCombat && len(def.EnemyIDs) > 0 {
		chosenEnemyID := def.EnemyIDs[rand.Intn(len(def.EnemyIDs))]
		enemy, err := combat.NewEnemyByID(chosenEnemyID)
		if err == nil {
			r.Enemy = enemy
		}
	}

	return r
}

// NewExpedition generates a procedural 9-node DAG dungeon run across 6 depths
func NewExpedition(p *character.Player, startingRations int) (*Expedition, error) {
	roomDefs, err := data.LoadRoomDefs()
	if err != nil {
		return nil, fmt.Errorf("gagal memuat data ruangan: %w", err)
	}

	defMap := make(map[string]data.RoomDef)
	for _, rd := range roomDefs {
		defMap[rd.ID] = rd
	}

	// 1. Procedural type selection with strict pacing and lifeline guarantees
	// Depth 1: Node 0 (Entrance) - 90% Combat, 5% Mystery, 5% Treasure (rare inactive catacomb entrance)
	d1Choices := []weightedChoice{
		{roomType: RoomTypeCombat, weight: 90},
		{roomType: RoomTypeMystery, weight: 5},
		{roomType: RoomTypeTreasure, weight: 5},
	}
	type0 := pickWeightedType(d1Choices)
	depth1HasTreasure := (type0 == RoomTypeTreasure)

	// Depth 2: Nodes 1 & 2
	d2Choices1 := []weightedChoice{
		{roomType: RoomTypeCombat, weight: 45},
		{roomType: RoomTypeTreasure, weight: 25},
		{roomType: RoomTypeRest, weight: 15},
		{roomType: RoomTypeMystery, weight: 15},
	}
	if depth1HasTreasure {
		d2Choices1 = []weightedChoice{
			{roomType: RoomTypeCombat, weight: 60},
			{roomType: RoomTypeTreasure, weight: 15},
			{roomType: RoomTypeRest, weight: 10},
			{roomType: RoomTypeMystery, weight: 15},
		}
	}
	type1 := pickWeightedType(d2Choices1)

	var d2Choices2 []weightedChoice
	switch type1 {
	case RoomTypeRest:
		d2Choices2 = []weightedChoice{
			{roomType: RoomTypeCombat, weight: 55},
			{roomType: RoomTypeTreasure, weight: 30},
			{roomType: RoomTypeMystery, weight: 15},
		}
	case RoomTypeTreasure:
		d2Choices2 = []weightedChoice{
			{roomType: RoomTypeCombat, weight: 60},
			{roomType: RoomTypeRest, weight: 25},
			{roomType: RoomTypeMystery, weight: 15},
		}
	case RoomTypeCombat:
		d2Choices2 = []weightedChoice{
			{roomType: RoomTypeTreasure, weight: 45},
			{roomType: RoomTypeRest, weight: 30},
			{roomType: RoomTypeMystery, weight: 25},
		}
	default:
		d2Choices2 = []weightedChoice{
			{roomType: RoomTypeCombat, weight: 50},
			{roomType: RoomTypeTreasure, weight: 30},
			{roomType: RoomTypeRest, weight: 20},
		}
	}
	type2 := pickWeightedType(d2Choices2)

	depth2HasRest := (type1 == RoomTypeRest || type2 == RoomTypeRest)
	depth2HasTreasure := (type1 == RoomTypeTreasure || type2 == RoomTypeTreasure)

	// Depth 3: Nodes 3 & 4
	var type3, type4 string
	if !depth2HasRest {
		// Lifeline guarantee: at least one rest sanctuary appears in Depth 3 if none in Depth 2
		restSlot := rand.Intn(2)
		otherChoices := []weightedChoice{
			{roomType: RoomTypeCombat, weight: 60},
			{roomType: RoomTypeTreasure, weight: 25},
			{roomType: RoomTypeMystery, weight: 15},
		}
		otherType := pickWeightedType(otherChoices)
		if restSlot == 0 {
			type3 = RoomTypeRest
			type4 = otherType
		} else {
			type3 = otherType
			type4 = RoomTypeRest
		}
	} else {
		// Pacing rule: no consecutive rests, rest weight set to 0
		combatWeight := 55
		treasureWeight := 30
		mysteryWeight := 15
		if depth2HasTreasure {
			combatWeight = 65
			treasureWeight = 20
			mysteryWeight = 15
		}

		d3Choices1 := []weightedChoice{
			{roomType: RoomTypeCombat, weight: combatWeight},
			{roomType: RoomTypeTreasure, weight: treasureWeight},
			{roomType: RoomTypeMystery, weight: mysteryWeight},
		}
		type3 = pickWeightedType(d3Choices1)

		var d3Choices2 []weightedChoice
		switch type3 {
		case RoomTypeCombat:
			d3Choices2 = []weightedChoice{
				{roomType: RoomTypeTreasure, weight: 60},
				{roomType: RoomTypeMystery, weight: 40},
			}
		case RoomTypeTreasure:
			d3Choices2 = []weightedChoice{
				{roomType: RoomTypeCombat, weight: 70},
				{roomType: RoomTypeMystery, weight: 30},
			}
		default:
			d3Choices2 = []weightedChoice{
				{roomType: RoomTypeCombat, weight: 60},
				{roomType: RoomTypeTreasure, weight: 40},
			}
		}
		type4 = pickWeightedType(d3Choices2)
	}

	depth3HasRest := (type3 == RoomTypeRest || type4 == RoomTypeRest)
	depth3HasTreasure := (type3 == RoomTypeTreasure || type4 == RoomTypeTreasure)

	// Depth 4: Nodes 5 & 6 (Pre-boss chambers)
	var type5, type6 string
	if depth3HasRest {
		// No consecutive rest rooms
		combatWeight := 60
		treasureWeight := 25
		mysteryWeight := 15
		if depth3HasTreasure {
			combatWeight = 70
			treasureWeight = 15
			mysteryWeight = 15
		}
		d4Choices1 := []weightedChoice{
			{roomType: RoomTypeCombat, weight: combatWeight},
			{roomType: RoomTypeTreasure, weight: treasureWeight},
			{roomType: RoomTypeMystery, weight: mysteryWeight},
		}
		type5 = pickWeightedType(d4Choices1)

		var d4Choices2 []weightedChoice
		switch type5 {
		case RoomTypeCombat:
			d4Choices2 = []weightedChoice{
				{roomType: RoomTypeTreasure, weight: 60},
				{roomType: RoomTypeMystery, weight: 40},
			}
		case RoomTypeTreasure:
			d4Choices2 = []weightedChoice{
				{roomType: RoomTypeCombat, weight: 70},
				{roomType: RoomTypeMystery, weight: 30},
			}
		default:
			d4Choices2 = []weightedChoice{
				{roomType: RoomTypeCombat, weight: 60},
				{roomType: RoomTypeTreasure, weight: 40},
			}
		}
		type6 = pickWeightedType(d4Choices2)
	} else {
		// Depth 3 had no rest, rest sanctuary available before boss
		restWeight := 25
		combatWeight := 45
		treasureWeight := 20
		mysteryWeight := 10
		if depth3HasTreasure {
			combatWeight = 65
			restWeight = 10
			treasureWeight = 10
			mysteryWeight = 15
		}

		d4Choices1 := []weightedChoice{
			{roomType: RoomTypeCombat, weight: combatWeight},
			{roomType: RoomTypeTreasure, weight: treasureWeight},
			{roomType: RoomTypeRest, weight: restWeight},
			{roomType: RoomTypeMystery, weight: mysteryWeight},
		}
		type5 = pickWeightedType(d4Choices1)

		var d4Choices2 []weightedChoice
		switch type5 {
		case RoomTypeRest:
			d4Choices2 = []weightedChoice{
				{roomType: RoomTypeCombat, weight: 65},
				{roomType: RoomTypeTreasure, weight: 20},
				{roomType: RoomTypeMystery, weight: 15},
			}
		case RoomTypeTreasure:
			d4Choices2 = []weightedChoice{
				{roomType: RoomTypeCombat, weight: 60},
				{roomType: RoomTypeRest, weight: restWeight},
				{roomType: RoomTypeMystery, weight: 15},
			}
		case RoomTypeCombat:
			d4Choices2 = []weightedChoice{
				{roomType: RoomTypeTreasure, weight: 35},
				{roomType: RoomTypeRest, weight: restWeight},
				{roomType: RoomTypeMystery, weight: 25},
			}
		default:
			d4Choices2 = []weightedChoice{
				{roomType: RoomTypeCombat, weight: 55},
				{roomType: RoomTypeTreasure, weight: 25},
				{roomType: RoomTypeRest, weight: restWeight},
			}
		}
		type6 = pickWeightedType(d4Choices2)
	}

	usedIDs := make(map[string]bool)

	// Node 0: Depth 1 (Entrance)
	def0 := pickDefByType(type0, 1, defMap, usedIDs)
	node0 := createRoomNode(0, 1, "Lorong Masuk", []int{1, 2}, def0)

	// Node 1: Depth 2 (Lorong Kiri)
	def1 := pickDefByType(type1, 2, defMap, usedIDs)
	node1 := createRoomNode(1, 2, "Lorong Kiri", []int{3, 4}, def1)

	// Node 2: Depth 2 (Lorong Kanan)
	def2 := pickDefByType(type2, 2, defMap, usedIDs)
	node2 := createRoomNode(2, 2, "Lorong Kanan", []int{3, 4}, def2)

	// Node 3: Depth 3 (Celah Kiri)
	def3 := pickDefByType(type3, 3, defMap, usedIDs)
	node3 := createRoomNode(3, 3, "Celah Kiri", []int{5, 6}, def3)

	// Node 4: Depth 3 (Celah Kanan)
	def4 := pickDefByType(type4, 3, defMap, usedIDs)
	node4 := createRoomNode(4, 3, "Celah Kanan", []int{5, 6}, def4)

	// Node 5: Depth 4 (Kubah Kiri)
	def5 := pickDefByType(type5, 4, defMap, usedIDs)
	node5 := createRoomNode(5, 4, "Kubah Kiri", []int{7}, def5)

	// Node 6: Depth 4 (Kubah Kanan)
	def6 := pickDefByType(type6, 4, defMap, usedIDs)
	node6 := createRoomNode(6, 4, "Kubah Kanan", []int{7}, def6)

	// Node 7: Depth 5 (Boss Chamber)
	def7 := pickDefByType(RoomTypeCombat, 5, defMap, usedIDs)
	node7 := createRoomNode(7, 5, "Kamar Persemayaman", []int{8}, def7)

	// Node 8: Depth 6 (Exit Stairs)
	def8 := pickDefByType(RoomTypeExit, 6, defMap, usedIDs)
	node8 := createRoomNode(8, 6, "Tangga Keluar", []int{}, def8)

	rooms := []*Room{node0, node1, node2, node3, node4, node5, node6, node7, node8}

	exp := &Expedition{
		Player:             p,
		Floor:              1,
		CurrentRoomIdx:     0,
		Rooms:              rooms,
		TotalDepths:        6,
		RoomsExploredCount: 1,
		Torch:              100,
		Rations:            startingRations,
		MaxBackpack:        p.MaxBackpack,
		Logs:               make([]string, 0),
	}

	exp.AddLog("[+] Pintu gerbang Katakombe terbuka, hawa dingin menyambut Anda")
	exp.AddLog(fmt.Sprintf("[i] Ransel memuat %d ransum perbekalan dan obor menyala penuh 100%%", startingRations))

	// Setup initial room
	firstRoom := exp.CurrentRoom()
	if firstRoom.Def.Type == RoomTypeCombat && firstRoom.Enemy != nil {
		exp.ActiveCombat = combat.NewCombatSession(p, firstRoom.Enemy)
		exp.AddLog(fmt.Sprintf("[!] Ruang 1: Musuh mendekat! Bersiap bertarung melawan %s", firstRoom.Enemy.Name))
	} else if firstRoom.Def.Type == RoomTypeTreasure {
		exp.AddLog("[*] Ruang 1: Lorong masuk lengang dan terbengkalai, tampak peti tua tak bertuan")
	} else if firstRoom.Def.Type == RoomTypeMystery {
		exp.AddLog("[*] Ruang 1: Lorong masuk hening tanpa penjaga, sebuah altar kuno menyambut Anda")
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

// NextRoomChoices returns the available connected rooms from current room
func (exp *Expedition) NextRoomChoices() []*Room {
	curr := exp.CurrentRoom()
	if curr == nil || len(curr.NextRoomIndices) == 0 {
		return nil
	}
	choices := make([]*Room, 0, len(curr.NextRoomIndices))
	for _, idx := range curr.NextRoomIndices {
		if idx >= 0 && idx < len(exp.Rooms) {
			choices = append(choices, exp.Rooms[idx])
		}
	}
	return choices
}

// AdvanceToRoom moves the expedition to a specific connected room
func (exp *Expedition) AdvanceToRoom(targetIdx int) error {
	if exp.IsCompleted || exp.IsDefeated {
		return fmt.Errorf("ekspedisi telah selesai")
	}

	if exp.ActiveCombat != nil && !exp.ActiveCombat.IsOver {
		return fmt.Errorf("selesaikan pertempuran terlebih dahulu")
	}

	curr := exp.CurrentRoom()
	if !curr.IsResolved {
		return fmt.Errorf("selesaikan peristiwa di ruangan ini terlebih dahulu")
	}

	valid := false
	for _, idx := range curr.NextRoomIndices {
		if idx == targetIdx {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("ruangan tujuan tidak dapat diakses dari lorong ini")
	}

	if targetIdx < 0 || targetIdx >= len(exp.Rooms) {
		return fmt.Errorf("indeks ruangan tidak valid")
	}

	// Torch decay
	exp.Torch -= 10
	if exp.Torch < 0 {
		exp.Torch = 0
	}
	if exp.Torch <= 0 {
		exp.AddLog("[!] Obor padam total! Kegelapan malam menyulitkan langkah Anda")
	}

	exp.CurrentRoomIdx = targetIdx
	exp.RoomsExploredCount++
	room := exp.CurrentRoom()

	exp.AddLog(fmt.Sprintf("[*] Melangkah ke Ruang %d (%s): %s", room.Depth, room.BranchName, room.Def.Title))

	// If combat room, initialize battle
	if room.Def.Type == RoomTypeCombat && room.Enemy != nil && !room.IsResolved {
		exp.ActiveCombat = combat.NewCombatSession(exp.Player, room.Enemy)
		exp.AddLog(fmt.Sprintf("[!] Musuh terdeteksi! Bersiap menghadapi %s", room.Enemy.Name))
	} else {
		exp.ActiveCombat = nil
	}

	return nil
}

// AdvanceRoom moves the expedition to the first connected room
func (exp *Expedition) AdvanceRoom() error {
	curr := exp.CurrentRoom()
	if len(curr.NextRoomIndices) == 0 {
		return fmt.Errorf("sudah berada di ruangan paling ujung")
	}
	return exp.AdvanceToRoom(curr.NextRoomIndices[0])
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

	if exp.Player.HasCompanionRole("Rogue") {
		rogueBonus := goldGained / 4
		if rogueBonus < 5 {
			rogueBonus = 5
		}
		goldGained += rogueBonus
		exp.AddLog(fmt.Sprintf("[*] Pencuri Valen membongkar kompartemen rahasia peti (+%d Emas bonus)", rogueBonus))
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

	if exp.Player.HasCompanionRole("Acolyte") {
		acolyteHeal := 25
		exp.Player.HP += acolyteHeal
		if exp.Player.HP > exp.Player.MaxHP {
			exp.Player.HP = exp.Player.MaxHP
		}
		exp.Player.Sanity += 15
		if exp.Player.Sanity > exp.Player.MaxSanity {
			exp.Player.Sanity = exp.Player.MaxSanity
		}
		exp.AddLog("[*] Suster Selene memanjatkan doa ketenangan (+25 HP & +15 Sanity)")
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

	hasScholar := exp.Player.HasCompanionRole("Scholar")
	success := hasScholar || exp.Player.Stats.Ingenuity >= room.Def.ReqValue
	room.IsResolved = true

	if success {
		bonusReward := 0
		if hasScholar {
			bonusReward = 30
			exp.AddLog("[*] Sarjana Alden menerjemahkan inskripsi rune purba dengan sempurna")
		}
		exp.GoldFound += room.Def.RewardGold + bonusReward
		resText := fmt.Sprintf("Teka-teki aksara terpecahkan! Altar membuka relik kuno (+%d Emas)", room.Def.RewardGold+bonusReward)
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

// ConsumePotion consumes an alchemy potion from player pouch during expedition
func (exp *Expedition) ConsumePotion(potionID string) (int, error) {
	if exp.Player.GetPotionCount(potionID) <= 0 {
		return 0, fmt.Errorf("stok ramuan tidak tersedia di ransel")
	}

	switch potionID {
	case "salep_pemulih":
		if exp.Player.HP >= exp.Player.MaxHP {
			return 0, fmt.Errorf("HP karakter sudah maksimal")
		}
		exp.Player.UsePotion(potionID)
		healAmount := 35
		oldHP := exp.Player.HP
		exp.Player.HP += healAmount
		if exp.Player.HP > exp.Player.MaxHP {
			exp.Player.HP = exp.Player.MaxHP
		}
		actualHeal := exp.Player.HP - oldHP
		exp.AddLog(fmt.Sprintf("[+] Mengoleskan Salep Pemulih Herbal (+%d HP, HP: %d/%d)", actualHeal, exp.Player.HP, exp.Player.MaxHP))
		return actualHeal, nil

	case "minyak_obor":
		if exp.Torch >= 100 {
			return 0, fmt.Errorf("nyala obor masih maksimal (100%%)")
		}
		exp.Player.UsePotion(potionID)
		exp.Torch += 35
		if exp.Torch > 100 {
			exp.Torch = 100
		}
		exp.AddLog(fmt.Sprintf("[+] Menuangkan Minyak Obor Murni (Penerangan pulih ke %d%%)", exp.Torch))
		return 35, nil

	case "tonik_penenang":
		if exp.Player.Sanity >= exp.Player.MaxSanity {
			return 0, fmt.Errorf("kewarasan karakter sudah maksimal")
		}
		exp.Player.UsePotion(potionID)
		gainSanity := 25
		exp.Player.Sanity += gainSanity
		if exp.Player.Sanity > exp.Player.MaxSanity {
			exp.Player.Sanity = exp.Player.MaxSanity
		}
		exp.AddLog(fmt.Sprintf("[+] Meminum Tonik Penenang Jiwa (+%d Sanity, Kewarasan: %d/%d)", gainSanity, exp.Player.Sanity, exp.Player.MaxSanity))
		return gainSanity, nil

	case "penawar_racun":
		exp.Player.UsePotion(potionID)
		healAmount := 15
		exp.Player.HP += healAmount
		if exp.Player.HP > exp.Player.MaxHP {
			exp.Player.HP = exp.Player.MaxHP
		}
		exp.AddLog(fmt.Sprintf("[+] Menenggak Penawar Racun Alami (+%d HP, racun dan pendarahan ternetralisir)", healAmount))
		return healAmount, nil

	case "eliksir_kekuatan":
		exp.Player.UsePotion(potionID)
		exp.AddLog("[+] Menenggak Eliksir Kekuatan Tempur (Kekuatan tebasan bertambah +8 ATK)")
		return 8, nil

	default:
		return 0, fmt.Errorf("ramuan tidak dikenal")
	}
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
