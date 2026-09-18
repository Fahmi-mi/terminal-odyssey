package economy

import (
	"strings"
	"testing"

	"github.com/Fahmi-mi/terminal-odyssey/internal/settlement"
)

func TestMarketState_PricesAndSeasons(t *testing.T) {
	market, err := NewMarketState()
	if err != nil {
		t.Fatalf("failed to create market state: %v", err)
	}

	if len(market.Items) < 8 {
		t.Fatalf("expected at least 8 commodities, got %d", len(market.Items))
	}

	// Spring
	market.UpdateDailyPrices(1, settlement.SeasonSpring)
	lumberSpring := market.GetItem("lumber").BuyPrice

	// Winter: lumber should surge
	market.UpdateDailyPrices(1, settlement.SeasonWinter)
	lumberWinter := market.GetItem("lumber").BuyPrice
	if lumberWinter <= lumberSpring {
		t.Errorf("expected winter lumber price (%d) > spring lumber price (%d)", lumberWinter, lumberSpring)
	}

	// Summer: rations should be cheaper
	market.UpdateDailyPrices(1, settlement.SeasonSummer)
	rationsSummer := market.GetItem("rations").BuyPrice
	market.UpdateDailyPrices(1, settlement.SeasonWinter)
	rationsWinter := market.GetItem("rations").BuyPrice
	if rationsWinter <= rationsSummer {
		t.Errorf("expected winter rations price (%d) > summer rations price (%d)", rationsWinter, rationsSummer)
	}

	// Verify buy price > sell price across all items
	for _, item := range market.Items {
		if item.BuyPrice <= item.SellPrice {
			t.Errorf("expected %s BuyPrice (%d) > SellPrice (%d)", item.Def.ID, item.BuyPrice, item.SellPrice)
		}
	}
}

func TestMarketState_IngenuityBarter(t *testing.T) {
	market, err := NewMarketState()
	if err != nil {
		t.Fatalf("failed to create market state: %v", err)
	}

	market.UpdateDailyPrices(1, settlement.SeasonSpring)
	item := market.GetItem("spices")

	baseBuy := market.GetEffectiveBuyPrice(item.Def.ID, 10)
	baseSell := market.GetEffectiveSellPrice(item.Def.ID, 10)

	// Ingenuity 20: 10% discount on buy, 15% bonus on sell
	smartBuy := market.GetEffectiveBuyPrice(item.Def.ID, 20)
	smartSell := market.GetEffectiveSellPrice(item.Def.ID, 20)

	if smartBuy >= baseBuy {
		t.Errorf("expected smart buy price (%d) < base buy price (%d)", smartBuy, baseBuy)
	}
	if smartSell <= baseSell {
		t.Errorf("expected smart sell price (%d) > base sell price (%d)", smartSell, baseSell)
	}

	// Sell price should always be strictly less than buy price to prevent local arbitrage
	if smartSell >= smartBuy {
		t.Errorf("expected smart sell price (%d) < smart buy price (%d)", smartSell, smartBuy)
	}
}

func TestCaravanManager_DispatchAndAdvance(t *testing.T) {
	cm, err := NewCaravanManager()
	if err != nil {
		t.Fatalf("failed to create caravan manager: %v", err)
	}

	route := cm.GetRoute("riverfall")
	if route == nil {
		t.Fatalf("expected route riverfall to exist")
	}

	// Pos Kafilah Level 0 cannot dispatch
	err = cm.CanDispatch("riverfall", 0, 100, nil, 20, 20)
	if err == nil {
		t.Errorf("expected error for level 0 building")
	}

	// Insufficient treasury
	err = cm.CanDispatch("riverfall", 1, 10, nil, 20, 20)
	if err == nil {
		t.Errorf("expected error for insufficient treasury")
	}

	// Insufficient cargo
	err = cm.CanDispatch("riverfall", 1, 100, nil, 5, 20)
	if err == nil {
		t.Errorf("expected error for insufficient cargo")
	}

	// Valid dispatch in spring
	err = cm.CanDispatch("riverfall", 1, 100, nil, 20, 20)
	if err != nil {
		t.Fatalf("unexpected error validating dispatch: %v", err)
	}

	exp, err := cm.Dispatch("riverfall", 1, settlement.SeasonSpring, 2)
	if err != nil {
		t.Fatalf("failed to dispatch caravan: %v", err)
	}
	if exp.TotalDays != 2 {
		t.Errorf("expected 2 total days, got %d", exp.TotalDays)
	}

	// Winter doubles travel duration
	expWinter, err := cm.Dispatch("ironpeak", 2, settlement.SeasonWinter, 2)
	if err != nil {
		t.Fatalf("failed to dispatch winter caravan: %v", err)
	}
	if expWinter.TotalDays != 6 {
		t.Errorf("expected 6 total days in winter for ironpeak, got %d", expWinter.TotalDays)
	}

	// Advance day 1
	results1 := cm.AdvanceDay(2, 5)
	if len(results1) != 0 {
		t.Errorf("expected no caravans finished on day 1")
	}

	// Advance day 2: riverfall finishes
	results2 := cm.AdvanceDay(3, 5)
	if len(results2) != 1 {
		t.Fatalf("expected 1 caravan finished on day 2, got %d", len(results2))
	}

	res := results2[0]
	if !strings.Contains(res.Log, "Lembah Riverfall") {
		t.Errorf("expected riverfall in log, got %s", res.Log)
	}
	if res.Expedition.ReturnedGold <= 0 {
		t.Errorf("expected positive returned gold, got %d", res.Expedition.ReturnedGold)
	}
}
