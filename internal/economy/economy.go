package economy

import (
	"fmt"
	"math"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

// MarketItem tracks current market price and trends for a commodity
type MarketItem struct {
	Def              data.CommodityDef
	BuyPrice         int
	SellPrice        int
	PreviousBuyPrice int
	Trend            string // ▲ Naik, ▼ Turun, ─ Stabil
}

// MarketState coordinates commodity prices, seasonal shifts, and daily fluctuations
type MarketState struct {
	Items       []MarketItem
	MarketEvent string
}

// NewMarketState creates a market initialized from commodities.json
func NewMarketState() (*MarketState, error) {
	defs, err := data.LoadCommodityDefs()
	if err != nil {
		return nil, err
	}

	items := make([]MarketItem, len(defs))
	for i, d := range defs {
		items[i] = MarketItem{
			Def:              d,
			BuyPrice:         d.BaseBuyPrice,
			SellPrice:        d.BaseSellPrice,
			PreviousBuyPrice: d.BaseBuyPrice,
			Trend:            "─ Stabil",
		}
	}

	m := &MarketState{
		Items:       items,
		MarketEvent: "Pasar lokal beroperasi normal",
	}
	m.UpdateDailyPrices(1, settlement.SeasonSpring)
	return m, nil
}

// UpdateDailyPrices recalculates commodity prices based on season and day variance
func (m *MarketState) UpdateDailyPrices(day int, season settlement.Season) {
	var eventText string

	switch season {
	case settlement.SeasonSpring:
		eventText = "Musim Semi: Pasokan herba segar melimpah, jalur perdagangan terbuka"
	case settlement.SeasonSummer:
		eventText = "Musim Panas: Panen raya melimpah, pasokan ransum murah dan melimpah"
	case settlement.SeasonAutumn:
		eventText = "Musim Gugur: Penebangan kayu meningkat, permintaan ransum awetan naik"
	case settlement.SeasonWinter:
		eventText = "Musim Dingin: Suhu beku ekstrem, harga kayu bakar dan ransum melonjak tinggi"
	}
	m.MarketEvent = eventText

	for i := range m.Items {
		item := &m.Items[i]
		item.PreviousBuyPrice = item.BuyPrice

		// 1. Season multiplier
		seasonMult := 1.0
		switch item.Def.ID {
		case "lumber":
			switch season {
			case settlement.SeasonAutumn:
				seasonMult = 0.80 // Surplus lumber
			case settlement.SeasonWinter:
				seasonMult = 1.45 // High firewood demand
			}
		case "stone":
			switch season {
			case settlement.SeasonWinter:
				seasonMult = 1.15 // Mining frozen ground
			}
		case "rations":
			switch season {
			case settlement.SeasonSummer:
				seasonMult = 0.75 // Abundant food
			case settlement.SeasonAutumn:
				seasonMult = 1.20 // Preserving demand
			case settlement.SeasonWinter:
				seasonMult = 1.40 // Scarce food
			}
		case "herbal_salve":
			switch season {
			case settlement.SeasonSpring:
				seasonMult = 0.85 // Fresh spring herbs
			case settlement.SeasonWinter:
				seasonMult = 1.25 // Illness remedies
			}
		case "spices", "caravan_silk":
			switch season {
			case settlement.SeasonWinter:
				seasonMult = 1.25 // Snow-blocked trade routes
			}
		}

		// 2. Daily pseudo-random fluctuation (-15% to +15%)
		// Deterministic based on day and commodity index
		seed := (day*17 + i*31 + int(season)*13) % 100
		dailyVariance := float64(seed-50) / 333.0 // roughly -0.15 to +0.15

		newBuy := int(math.Round(float64(item.Def.BaseBuyPrice) * seasonMult * (1.0 + dailyVariance)))
		newSell := int(math.Round(float64(item.Def.BaseSellPrice) * seasonMult * (1.0 + dailyVariance)))

		if newBuy < 2 {
			newBuy = 2
		}
		if newSell < 1 {
			newSell = 1
		}
		if newSell >= newBuy {
			newSell = newBuy - 1
		}

		item.BuyPrice = newBuy
		item.SellPrice = newSell

		if item.BuyPrice > item.PreviousBuyPrice {
			item.Trend = "▲ Naik"
		} else if item.BuyPrice < item.PreviousBuyPrice {
			item.Trend = "▼ Turun"
		} else {
			item.Trend = "─ Stabil"
		}
	}
}

// GetItem retrieves a market item by ID
func (m *MarketState) GetItem(id string) *MarketItem {
	for i := range m.Items {
		if m.Items[i].Def.ID == id {
			return &m.Items[i]
		}
	}
	return nil
}

// GetEffectiveBuyPrice calculates buy price with Ingenuity discount (1% per point above 10, max 15%)
func (m *MarketState) GetEffectiveBuyPrice(id string, ingenuity int) int {
	item := m.GetItem(id)
	if item == nil {
		return 0
	}

	diff := ingenuity - 10
	if diff < 0 {
		diff = 0
	}
	discountPct := float64(diff) * 0.01
	if discountPct > 0.15 {
		discountPct = 0.15
	}

	price := int(math.Round(float64(item.BuyPrice) * (1.0 - discountPct)))
	if price < 1 {
		price = 1
	}
	return price
}

// GetEffectiveSellPrice calculates sell price with Ingenuity barter bonus (1.5% per point above 10, max 22%)
func (m *MarketState) GetEffectiveSellPrice(id string, ingenuity int) int {
	item := m.GetItem(id)
	if item == nil {
		return 0
	}

	diff := ingenuity - 10
	if diff < 0 {
		diff = 0
	}
	bonusPct := float64(diff) * 0.015
	if bonusPct > 0.22 {
		bonusPct = 0.22
	}

	price := int(math.Round(float64(item.SellPrice) * (1.0 + bonusPct)))
	buyPrice := m.GetEffectiveBuyPrice(id, ingenuity)
	if price >= buyPrice {
		price = buyPrice - 1
	}
	if price < 1 {
		price = 1
	}
	return price
}

// CaravanExpedition represents an active caravan on a regional trade route
type CaravanExpedition struct {
	Route          data.CaravanRouteDef
	DispatchDay    int
	TotalDays      int
	DaysRemaining  int
	GoldInvestment int
	CargoCommodity string
	CargoAmount    int
	IsFinished     bool
	Ambushed       bool
	ReturnedGold   int
	BonusItem      string
	BonusAmount    int
	ResultLog      string
}

// CaravanResult holds the outcome of a finished caravan journey
type CaravanResult struct {
	Expedition *CaravanExpedition
	Log        string
}

// CaravanManager tracks regional routes and in-flight trade caravans
type CaravanManager struct {
	Routes         []data.CaravanRouteDef
	ActiveCaravans []*CaravanExpedition
}

// NewCaravanManager initializes caravan routes from data
func NewCaravanManager() (*CaravanManager, error) {
	routes, err := data.LoadCaravanRouteDefs()
	if err != nil {
		return nil, err
	}
	return &CaravanManager{
		Routes:         routes,
		ActiveCaravans: make([]*CaravanExpedition, 0),
	}, nil
}

// GetRoute returns route definition by ID
func (cm *CaravanManager) GetRoute(routeID string) *data.CaravanRouteDef {
	for i := range cm.Routes {
		if cm.Routes[i].ID == routeID {
			return &cm.Routes[i]
		}
	}
	return nil
}

// MaxConcurrentCaravans returns maximum caravans allowed by Pos Kafilah level
func (cm *CaravanManager) MaxConcurrentCaravans(buildingLevel int) int {
	if buildingLevel < 1 {
		return 0
	}
	return buildingLevel
}

// CanDispatch validates if a caravan can be launched
func (cm *CaravanManager) CanDispatch(routeID string, buildingLevel, treasury int, commodities map[string]int, lumber, rations int) error {
	route := cm.GetRoute(routeID)
	if route == nil {
		return fmt.Errorf("rute dagang tidak ditemukan")
	}

	if buildingLevel < route.RequiredLevel {
		return fmt.Errorf("butuh Pos Kafilah Level %d untuk rute %s", route.RequiredLevel, route.Name)
	}

	maxAllowed := cm.MaxConcurrentCaravans(buildingLevel)
	if len(cm.ActiveCaravans) >= maxAllowed {
		return fmt.Errorf("seluruh armada kafilah (%d/%d) sedang bertugas di perjalanan", len(cm.ActiveCaravans), maxAllowed)
	}

	for _, active := range cm.ActiveCaravans {
		if active.Route.ID == routeID {
			return fmt.Errorf("kafilah menuju %s sedang aktif dalam perjalanan", route.Name)
		}
	}

	if treasury < route.GoldInvestment {
		return fmt.Errorf("kas emas tidak mencukupi modal kafilah (butuh %d Gold, ada %d)", route.GoldInvestment, treasury)
	}

	switch route.CargoCommodity {
	case "lumber":
		if lumber < route.CargoAmount {
			return fmt.Errorf("pasokan kayu tidak mencukupi kargo (butuh %d Kayu, ada %d)", route.CargoAmount, lumber)
		}
	case "rations":
		if rations < route.CargoAmount {
			return fmt.Errorf("pasokan ransum tidak mencukupi kargo (butuh %d Ransum, ada %d)", route.CargoAmount, rations)
		}
	default:
		currentStock := 0
		if commodities != nil {
			currentStock = commodities[route.CargoCommodity]
		}
		if currentStock < route.CargoAmount {
			return fmt.Errorf("komoditas kargo %s tidak mencukupi (butuh %d, ada %d)", route.CargoCommodity, route.CargoAmount, currentStock)
		}
	}

	return nil
}

// Dispatch launches a new caravan along the route
func (cm *CaravanManager) Dispatch(routeID string, currentDay int, season settlement.Season, militiaCount int) (*CaravanExpedition, error) {
	route := cm.GetRoute(routeID)
	if route == nil {
		return nil, fmt.Errorf("rute dagang tidak ditemukan")
	}

	totalDays := route.BaseDays
	if season == settlement.SeasonWinter {
		totalDays = route.BaseDays * 2 // Winter snow doubles travel duration
	}

	exp := &CaravanExpedition{
		Route:          *route,
		DispatchDay:    currentDay,
		TotalDays:      totalDays,
		DaysRemaining:  totalDays,
		GoldInvestment: route.GoldInvestment,
		CargoCommodity: route.CargoCommodity,
		CargoAmount:    route.CargoAmount,
		IsFinished:     false,
	}

	cm.ActiveCaravans = append(cm.ActiveCaravans, exp)
	return exp, nil
}

// AdvanceDay decrements travel days and resolves returning caravans
func (cm *CaravanManager) AdvanceDay(currentDay int, militiaCount int) []CaravanResult {
	var results []CaravanResult
	var remaining []*CaravanExpedition

	for _, exp := range cm.ActiveCaravans {
		exp.DaysRemaining--
		if exp.DaysRemaining <= 0 {
			exp.IsFinished = true

			// Ambush calculation: baseline risk reduced by militia defense
			effectiveRisk := exp.Route.AmbushRisk - (militiaCount * 3)
			if effectiveRisk < 0 {
				effectiveRisk = 0
			}

			// Deterministic roll based on dispatch day and route
			roll := (currentDay*23 + exp.Route.RequiredLevel*19) % 100
			if roll < effectiveRisk {
				exp.Ambushed = true
				exp.ReturnedGold = exp.Route.RewardGoldMin - (exp.GoldInvestment / 3)
				if exp.ReturnedGold < exp.GoldInvestment {
					exp.ReturnedGold = exp.GoldInvestment
				}
				exp.BonusItem = ""
				exp.BonusAmount = 0
				exp.ResultLog = fmt.Sprintf("[!] KAFILAH DISERGAP: Kafilah dari %s dihadang bandit! Sebagian muatan lenyap namun berhasil membawa pulang %d Gold",
					exp.Route.Name, exp.ReturnedGold)
			} else {
				exp.Ambushed = false
				exp.ReturnedGold = (exp.Route.RewardGoldMin + exp.Route.RewardGoldMax) / 2
				exp.BonusItem = exp.Route.BonusItemID
				exp.BonusAmount = exp.Route.BonusItemAmount
				exp.ResultLog = fmt.Sprintf("[+] KAFILAH KEMBALI: Armada dari %s tiba dengan sukses! Membawa pulang +%d Gold dan +%d %s",
					exp.Route.Name, exp.ReturnedGold, exp.BonusAmount, exp.BonusItem)
			}

			results = append(results, CaravanResult{
				Expedition: exp,
				Log:        exp.ResultLog,
			})
		} else {
			remaining = append(remaining, exp)
		}
	}

	cm.ActiveCaravans = remaining
	return results
}
