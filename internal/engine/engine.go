package engine

import (
	"fmt"
	"strings"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/alchemy"
	"github.com/Fahmi-mi/terminal-odyssey/internal/character"
	"github.com/Fahmi-mi/terminal-odyssey/internal/dungeon"
	"github.com/Fahmi-mi/terminal-odyssey/internal/economy"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
	"github.com/Fahmi-mi/terminal-odyssey/internal/tavern"
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
	StateAlchemyLab
)

// ExpeditionSummary stores the results of a finished dungeon run
type ExpeditionSummary struct {
	WasEvacuated    bool
	GoldEarned      int
	CompanionCut    int
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

	Market   *economy.MarketState
	Caravans *economy.CaravanManager
	Alchemy  *alchemy.AlchemyManager
	Tavern   *tavern.TavernManager

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

	market, _ := economy.NewMarketState()
	caravans, _ := economy.NewCaravanManager()
	alc, _ := alchemy.NewAlchemyManager()
	tav, _ := tavern.NewTavernManager()

	e := &Engine{
		CurrentState:  StateTownMenu,
		PreviousState: StateTownMenu,
		DayCounter:    1,
		CurrentSeason: settlement.SeasonSpring,
		DaysPerSeason: 15, // Every 15 days = 1 season
		Village:       settlement.NewSettlement(villageName),
		Player:        character.NewDefaultPlayer(playerName),
		Market:        market,
		Caravans:      caravans,
		Alchemy:       alc,
		Tavern:        tav,
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

	// Update market commodity prices
	if e.Market != nil {
		e.Market.UpdateDailyPrices(e.DayCounter, e.CurrentSeason)
	}

	// Advance active trade caravans
	if e.Caravans != nil {
		caravanResults := e.Caravans.AdvanceDay(e.DayCounter, e.Village.Workers.Militia)
		for _, res := range caravanResults {
			e.Village.Treasury += res.Expedition.ReturnedGold
			if res.Expedition.BonusItem != "" && res.Expedition.BonusAmount > 0 {
				e.Village.AddCommodity(res.Expedition.BonusItem, res.Expedition.BonusAmount)
			}
			simResult.Logs = append([]string{res.Log}, simResult.Logs...)
		}
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

	totalCompanionCut := 0
	if wasEvac {
		for _, comp := range e.Player.Party {
			if comp.IsAlive && comp.HP > 0 {
				cut := (exp.GoldFound * comp.CutPercent) / 100
				totalCompanionCut += cut
			}
		}
	}
	netGold := exp.GoldFound - totalCompanionCut
	if netGold < 0 {
		netGold = 0
	}

	summary := &ExpeditionSummary{
		WasEvacuated:    wasEvac,
		GoldEarned:      exp.GoldFound,
		CompanionCut:    totalCompanionCut,
		LumberEarned:    exp.LumberFound,
		StoneEarned:     exp.StoneFound,
		EnemiesDefeated: exp.EnemiesDefeated,
		RoomsExplored:   roomsExplored,
		TotalRooms:      totalRooms,
	}

	if wasEvac {
		exp.Evacuate()
		e.Village.Treasury += netGold
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

		if totalCompanionCut > 0 {
			e.DailyLogs = append([]string{
				fmt.Sprintf("[+] Ekspedisi sukses: membawa pulang +%d Emas (+%d kas desa, -%d upah rekan), +%d Kayu, +%d Batu", exp.GoldFound, netGold, totalCompanionCut, exp.LumberFound, exp.StoneFound),
			}, e.DailyLogs...)
		} else {
			e.DailyLogs = append([]string{
				fmt.Sprintf("[+] Ekspedisi sukses: membawa pulang +%d Emas, +%d Kayu, +%d Batu", exp.GoldFound, exp.LumberFound, exp.StoneFound),
			}, e.DailyLogs...)
		}
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

	// Handle fallen companions
	var survivingParty []character.Companion
	for _, c := range e.Player.Party {
		if c.IsAlive && c.HP > 0 {
			survivingParty = append(survivingParty, c)
		} else {
			if e.Tavern != nil {
				_ = e.Tavern.Dismiss(c.ID)
			}
			e.DailyLogs = append([]string{
				fmt.Sprintf("[!] Rekan %s gugur di dalam katakombe dan tidak kembali", c.Name),
			}, e.DailyLogs...)
		}
	}
	e.Player.Party = survivingParty

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

// SwitchWeapon switches the player's active weapon to an already owned weapon
func (e *Engine) SwitchWeapon(weaponID string) error {
	if err := e.Player.SwitchWeapon(weaponID); err != nil {
		return err
	}
	e.SetAlert(fmt.Sprintf("Berhasil memasang %s sebagai senjata aktif", e.Player.EquippedWeapon.Name))
	return nil
}

// RepairWeaponByID repairs a specific weapon by ID (either equipped or in armory)
func (e *Engine) RepairWeaponByID(weaponID string) error {
	bsLvl := e.Village.Buildings[settlement.BuildingBlacksmith]
	if bsLvl < 1 {
		return fmt.Errorf("bengkel Pandai Besi belum dibangun")
	}

	if weaponID == "" || weaponID == e.Player.EquippedWeapon.ID {
		return e.RepairEquippedWeapon()
	}

	var target *character.Weapon
	for i := range e.Player.OwnedWeapons {
		if e.Player.OwnedWeapons[i].ID == weaponID {
			target = &e.Player.OwnedWeapons[i]
			break
		}
	}

	if target == nil {
		return fmt.Errorf("senjata tidak ditemukan di inventaris")
	}

	if target.Durability >= target.MaxDura {
		return fmt.Errorf("ketahanan %s masih maksimal (%d/%d)", target.Name, target.Durability, target.MaxDura)
	}

	missing := target.MaxDura - target.Durability
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

	target.Durability = target.MaxDura
	e.SetAlert(fmt.Sprintf("Berhasil memperbaiki %s (+%d Durabilitas, -%d Gold, -%d Batu)", target.Name, missing, goldCost, stoneCost))
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

// BuyCommodity purchases a quantity of a commodity from the local market
func (e *Engine) BuyCommodity(commodityID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("jumlah beli harus minimal 1 unit")
	}

	if e.Market == nil {
		return fmt.Errorf("pasar belum tersedia")
	}

	item := e.Market.GetItem(commodityID)
	if item == nil {
		return fmt.Errorf("komoditas %s tidak ditemukan di katalog pasar", commodityID)
	}

	unitPrice := e.Market.GetEffectiveBuyPrice(commodityID, e.Player.Stats.Ingenuity)
	totalCost := unitPrice * amount

	if e.Village.Treasury < totalCost {
		return fmt.Errorf("kas emas tidak mencukupi (butuh %d Gold, ada %d)", totalCost, e.Village.Treasury)
	}

	maxL, maxS, maxR := e.Village.StorageCap()
	switch commodityID {
	case "lumber":
		if e.Village.Lumber+amount > maxL {
			return fmt.Errorf("kapasitas gudang kayu tidak mencukupi (maksimal %d)", maxL)
		}
	case "stone":
		if e.Village.Stone+amount > maxS {
			return fmt.Errorf("kapasitas gudang batu tidak mencukupi (maksimal %d)", maxS)
		}
	case "rations":
		if e.Village.Rations+amount > maxR {
			return fmt.Errorf("kapasitas gudang ransum tidak mencukupi (maksimal %d)", maxR)
		}
	}

	e.Village.Treasury -= totalCost
	e.Village.AddCommodity(commodityID, amount)

	e.SetAlert(fmt.Sprintf("Berhasil membeli %d %s (-%d Gold)", amount, item.Def.Name, totalCost))
	return nil
}

// SellCommodity sells a quantity of a commodity from village stock to the local market
func (e *Engine) SellCommodity(commodityID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("jumlah jual harus minimal 1 unit")
	}

	if e.Market == nil {
		return fmt.Errorf("pasar belum tersedia")
	}

	item := e.Market.GetItem(commodityID)
	if item == nil {
		return fmt.Errorf("komoditas %s tidak ditemukan di katalog pasar", commodityID)
	}

	currentStock := e.Village.GetCommodityStock(commodityID)
	if currentStock < amount {
		return fmt.Errorf("stok %s di desa tidak mencukupi (ada %d, butuh %d)", item.Def.Name, currentStock, amount)
	}

	unitPrice := e.Market.GetEffectiveSellPrice(commodityID, e.Player.Stats.Ingenuity)
	totalEarned := unitPrice * amount

	if err := e.Village.DeductCommodity(commodityID, amount); err != nil {
		return err
	}

	e.Village.Treasury += totalEarned
	e.SetAlert(fmt.Sprintf("Berhasil menjual %d %s (+%d Gold)", amount, item.Def.Name, totalEarned))
	return nil
}

// DispatchCaravan launches a regional trade caravan along the chosen route
func (e *Engine) DispatchCaravan(routeID string) error {
	if e.Caravans == nil {
		return fmt.Errorf("sistem kafilah belum siap")
	}

	postLvl := e.Village.Buildings[settlement.BuildingCaravanPost]
	if err := e.Caravans.CanDispatch(routeID, postLvl, e.Village.Treasury, e.Village.Commodities, e.Village.Lumber, e.Village.Rations); err != nil {
		return err
	}

	route := e.Caravans.GetRoute(routeID)
	if route == nil {
		return fmt.Errorf("rute tidak ditemukan")
	}

	// Deduct investments and cargo
	e.Village.Treasury -= route.GoldInvestment
	if err := e.Village.DeductCommodity(route.CargoCommodity, route.CargoAmount); err != nil {
		e.Village.Treasury += route.GoldInvestment
		return err
	}

	exp, err := e.Caravans.Dispatch(routeID, e.DayCounter, e.CurrentSeason, e.Village.Workers.Militia)
	if err != nil {
		e.Village.Treasury += route.GoldInvestment
		e.Village.AddCommodity(route.CargoCommodity, route.CargoAmount)
		return err
	}

	e.SetAlert(fmt.Sprintf("Kafilah dagang menuju %s berhasil diberangkatkan (%d hari perjalanan)", route.Name, exp.TotalDays))
	return nil
}

// BrewPotion crafts a potion using the alchemy lab and adds it to player pouch
func (e *Engine) BrewPotion(recipeID string) error {
	if e.Alchemy == nil {
		return fmt.Errorf("sistem alkimia belum siap")
	}

	apothecaryLvl := e.Village.Buildings[settlement.BuildingApothecary]
	if apothecaryLvl < 1 {
		return fmt.Errorf("laboratorium Alkimia belum didirikan di desa")
	}

	result, err := e.Alchemy.Brew(
		recipeID,
		apothecaryLvl,
		e.Village.Treasury,
		e.Village.Commodities,
		e.Village.Lumber,
		e.Village.Stone,
		e.Village.Rations,
		e.Player.Stats.Ingenuity,
	)
	if err != nil {
		return err
	}

	// Deduct costs
	e.Village.Treasury -= result.Recipe.GoldCost
	for ingID, amt := range result.Recipe.Ingredients {
		switch ingID {
		case "lumber":
			e.Village.Lumber -= amt
		case "stone":
			e.Village.Stone -= amt
		case "rations":
			e.Village.Rations -= amt
		default:
			_ = e.Village.DeductCommodity(ingID, amt)
		}
	}

	// Add to player pouch
	e.Player.AddPotion(result.Recipe.ID, result.YieldAmount)
	e.SetAlert(result.Log)
	return nil
}

// DrinkPotionInTown consumes a potion from player pouch while in town
func (e *Engine) DrinkPotionInTown(potionID string) error {
	if e.Player.GetPotionCount(potionID) <= 0 {
		return fmt.Errorf("anda tidak memiliki ramuan tersebut")
	}

	switch potionID {
	case "salep_pemulih":
		if e.Player.HP >= e.Player.MaxHP {
			return fmt.Errorf("darah (HP) karakter sudah maksimal")
		}
		e.Player.UsePotion(potionID)
		e.Player.HP += 45
		if e.Player.HP > e.Player.MaxHP {
			e.Player.HP = e.Player.MaxHP
		}
		e.SetAlert(fmt.Sprintf("Meminum Salep Pemulih (HP pulih ke %d/%d)", e.Player.HP, e.Player.MaxHP))

	case "tonik_penenang":
		if e.Player.Sanity >= e.Player.MaxSanity {
			return fmt.Errorf("kewarasan karakter sudah maksimal")
		}
		e.Player.UsePotion(potionID)
		e.Player.Sanity += 40
		if e.Player.Sanity > e.Player.MaxSanity {
			e.Player.Sanity = e.Player.MaxSanity
		}
		e.SetAlert(fmt.Sprintf("Meminum Tonik Penenang Jiwa (Kewarasan pulih ke %d/%d)", e.Player.Sanity, e.Player.MaxSanity))

	case "penawar_racun":
		e.Player.UsePotion(potionID)
		e.Player.HP += 15
		if e.Player.HP > e.Player.MaxHP {
			e.Player.HP = e.Player.MaxHP
		}
		e.Player.Sanity += 15
		if e.Player.Sanity > e.Player.MaxSanity {
			e.Player.Sanity = e.Player.MaxSanity
		}
		e.SetAlert("Meminum Penawar Racun (+15 HP, +15 Sanity)")

	case "minyak_obor":
		return fmt.Errorf("minyak obor hanya dapat digunakan saat penjelajahan lorong katakombe")

	case "eliksir_kekuatan":
		return fmt.Errorf("eliksir kekuatan hanya dapat digunakan saat pertempuran aktif")

	default:
		return fmt.Errorf("ramuan tidak dikenal")
	}
	return nil
}

// HireCompanion recruits a companion from the tavern into the party
func (e *Engine) HireCompanion(companionID string) error {
	if e.Tavern == nil {
		return fmt.Errorf("kedai Minum belum siap")
	}

	tavLvl := e.Village.Buildings[settlement.BuildingTavern]
	if tavLvl < 1 {
		return fmt.Errorf("kedai Minum belum didirikan di desa")
	}

	merc, err := e.Tavern.Hire(companionID, tavLvl, e.Village.Treasury, len(e.Player.Party))
	if err != nil {
		return err
	}

	e.Village.Treasury -= merc.Def.HireCost
	e.Player.AddCompanion(character.Companion{
		ID:         merc.Def.ID,
		Name:       merc.Def.Name,
		Role:       merc.Def.Role,
		PerkDesc:   merc.Def.PerkDesc,
		CutPercent: merc.Def.CutPercent,
		HP:         merc.Def.HP,
		MaxHP:      merc.Def.MaxHP,
	})

	e.SetAlert(fmt.Sprintf("%s (%s) bergabung dengan rombongan petualang", merc.Def.Name, merc.Def.RoleDisplay))
	return nil
}

// DismissCompanion releases a companion from the active party
func (e *Engine) DismissCompanion(companionID string) error {
	if e.Tavern == nil {
		return fmt.Errorf("kedai Minum belum siap")
	}

	if err := e.Tavern.Dismiss(companionID); err != nil {
		return err
	}
	e.Player.RemoveCompanion(companionID)
	e.SetAlert("Pendamping berhasil dikeluarkan dari rombongan")
	return nil
}

// TavernRest provides sanity and HP recovery through dining at the tavern
func (e *Engine) TavernRest() error {
	if e.Tavern == nil {
		return fmt.Errorf("kedai Minum belum siap")
	}

	tavLvl := e.Village.Buildings[settlement.BuildingTavern]
	if tavLvl < 1 {
		return fmt.Errorf("kedai Minum belum didirikan di desa")
	}

	if err := e.Tavern.CanRest(e.Village.Treasury, e.Village.Rations); err != nil {
		return err
	}

	e.Village.Treasury -= 10
	e.Village.Rations -= 1

	oldSanity := e.Player.Sanity
	e.Player.Sanity += 20
	if e.Player.Sanity > e.Player.MaxSanity {
		e.Player.Sanity = e.Player.MaxSanity
	}

	oldHP := e.Player.HP
	e.Player.HP += 15
	if e.Player.HP > e.Player.MaxHP {
		e.Player.HP = e.Player.MaxHP
	}

	e.SetAlert(fmt.Sprintf("Menikmati santapan hangat di kedai (+%d Sanity, +%d HP)", e.Player.Sanity-oldSanity, e.Player.HP-oldHP))
	return nil
}

// TavernRumor gets an atmospheric rumor from the tavern
func (e *Engine) TavernRumor() string {
	if e.Tavern == nil {
		return "Suasana kedai minum tampak sepi"
	}
	rumor := e.Tavern.GetRandomRumor()
	e.SetAlert(rumor)
	return rumor
}



