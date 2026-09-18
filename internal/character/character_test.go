package character

import (
	"testing"
)

func TestEquipAndRepairWeapon(t *testing.T) {
	p := NewDefaultPlayer("Hero")

	newWeapon := Weapon{
		ID:           "iron_sword",
		Name:         "Pedang Besi Tempa",
		WeaponType:   "Sword",
		BaseDamage:   [2]int{12, 18},
		CritRate:     0.10,
		Initiative:   12,
		Durability:   30,
		MaxDura:      50,
		SpecialAffix: "Tajam",
	}

	p.EquipWeapon(newWeapon)
	if p.EquippedWeapon.ID != "iron_sword" {
		t.Errorf("expected equipped weapon ID iron_sword, got %s", p.EquippedWeapon.ID)
	}
	if p.EquippedWeapon.Durability != 30 {
		t.Errorf("expected durability 30, got %d", p.EquippedWeapon.Durability)
	}

	repaired := p.RepairWeapon()
	if repaired != 20 {
		t.Errorf("expected 20 repaired durability, got %d", repaired)
	}
	if p.EquippedWeapon.Durability != 50 {
		t.Errorf("expected durability restored to 50, got %d", p.EquippedWeapon.Durability)
	}

	// Verify both starter weapon and new weapon are stored in OwnedWeapons
	if len(p.OwnedWeapons) != 2 {
		t.Fatalf("expected 2 owned weapons, got %d", len(p.OwnedWeapons))
	}
	if !p.OwnsWeapon("rusty_sword") || !p.OwnsWeapon("iron_sword") {
		t.Errorf("expected both rusty_sword and iron_sword to be owned")
	}

	// Switch back to rusty_sword
	if err := p.SwitchWeapon("rusty_sword"); err != nil {
		t.Fatalf("unexpected error switching to rusty_sword: %v", err)
	}
	if p.EquippedWeapon.ID != "rusty_sword" {
		t.Errorf("expected equipped weapon rusty_sword, got %s", p.EquippedWeapon.ID)
	}

	// Switch back to iron_sword
	if err := p.SwitchWeapon("iron_sword"); err != nil {
		t.Fatalf("unexpected error switching to iron_sword: %v", err)
	}
	if p.EquippedWeapon.ID != "iron_sword" {
		t.Errorf("expected equipped weapon iron_sword, got %s", p.EquippedWeapon.ID)
	}

	// Switch to non-existent weapon
	if err := p.SwitchWeapon("phantom_blade"); err == nil {
		t.Errorf("expected error switching to unowned weapon")
	}
}

func TestUpgradeStatAndDerivedStats(t *testing.T) {
	p := NewDefaultPlayer("Warrior")

	// Initial stats
	if p.Stats.Resolve != 10 || p.MaxHP != 100 {
		t.Fatalf("expected initial resolve 10 and max HP 100, got %d and %d", p.Stats.Resolve, p.MaxHP)
	}
	if p.Stats.Might != 10 || p.MaxBackpack != 12 {
		t.Fatalf("expected initial might 10 and backpack 12, got %d and %d", p.Stats.Might, p.MaxBackpack)
	}

	// Upgrade Resolve up to cap 12
	newRes, err := p.UpgradeStat("Resolve", 12)
	if err != nil {
		t.Fatalf("unexpected error upgrading resolve: %v", err)
	}
	if newRes != 11 {
		t.Errorf("expected resolve 11, got %d", newRes)
	}
	if p.MaxHP != 105 {
		t.Errorf("expected max HP 105 (100 + 5), got %d", p.MaxHP)
	}

	// Upgrade Might
	_, _ = p.UpgradeStat("Might", 15)
	_, _ = p.UpgradeStat("Might", 15)
	if p.Stats.Might != 12 {
		t.Errorf("expected might 12, got %d", p.Stats.Might)
	}
	if p.MaxBackpack != 13 {
		t.Errorf("expected backpack 13 (12 + 2/2), got %d", p.MaxBackpack)
	}

	// Test cap limit
	_, _ = p.UpgradeStat("Resolve", 12) // Resolve reaches 12
	_, errCap := p.UpgradeStat("Resolve", 12)
	if errCap == nil {
		t.Errorf("expected error when upgrading beyond cap 12")
	}

	// Test invalid stat name
	_, errInvalid := p.UpgradeStat("UnknownStat", 15)
	if errInvalid == nil {
		t.Errorf("expected error on invalid stat name")
	}
}

func TestPlayerPotionsAndCompanions(t *testing.T) {
	p := NewDefaultPlayer("Alchemist")

	// 1. Potions
	if p.TotalPotions() != 0 {
		t.Errorf("expected 0 initial potions, got %d", p.TotalPotions())
	}
	p.AddPotion("salep_pemulih", 2)
	p.AddPotion("minyak_obor", 1)

	if p.GetPotionCount("salep_pemulih") != 2 {
		t.Errorf("expected 2 salep_pemulih, got %d", p.GetPotionCount("salep_pemulih"))
	}
	if p.TotalPotions() != 3 {
		t.Errorf("expected 3 total potions, got %d", p.TotalPotions())
	}

	used := p.UsePotion("salep_pemulih")
	if !used || p.GetPotionCount("salep_pemulih") != 1 {
		t.Errorf("expected 1 salep_pemulih after use, got %d", p.GetPotionCount("salep_pemulih"))
	}

	usedNonExistent := p.UsePotion("unknown_potion")
	if usedNonExistent {
		t.Errorf("expected false when using non-existent potion")
	}

	// 2. Companions
	c := Companion{
		ID:         "valen_rogue",
		Name:       "Valen",
		Role:       "Rogue",
		CutPercent: 10,
		HP:         60,
		MaxHP:      60,
		IsAlive:    true,
	}
	p.AddCompanion(c)

	if !p.HasCompanionRole("Rogue") {
		t.Errorf("expected Rogue companion in party")
	}
	if p.HasCompanionRole("Scholar") {
		t.Errorf("expected no Scholar in party")
	}

	comp := p.GetCompanion("Rogue")
	if comp == nil || comp.Name != "Valen" {
		t.Errorf("expected Valen companion, got %v", comp)
	}

	removed := p.RemoveCompanion("valen_rogue")
	if !removed || p.HasCompanionRole("Rogue") {
		t.Errorf("expected Valen to be removed from party")
	}
}
