package combat

import (
	"testing"

	"github.com/Fahmi-mi/terminal-odyssey/internal/character"
)

func TestNewEnemyByID(t *testing.T) {
	e, err := NewEnemyByID("skeleton_scout")
	if err != nil {
		t.Fatalf("failed to load skeleton_scout: %v", err)
	}
	if e.Name != "Prajurit Kerangka" {
		t.Errorf("expected Prajurit Kerangka, got %s", e.Name)
	}
	if e.MaxHP <= 0 || e.HP <= 0 {
		t.Errorf("expected positive HP, got %d", e.HP)
	}

	_, errNotFound := NewEnemyByID("non_existent_enemy")
	if errNotFound == nil {
		t.Errorf("expected error for non existent enemy")
	}
}

func TestCombatAttackAndVictory(t *testing.T) {
	p := character.NewDefaultPlayer("Hero")
	e := &Enemy{
		ID:         "dummy",
		Name:       "Dummy Target",
		HP:         10,
		MaxHP:      10,
		MinDamage:  2,
		MaxDamage:  2,
		Initiative: 5,
		Defense:    0,
		GoldReward: 20,
	}

	session := NewCombatSession(p, e)
	if session.IsOver {
		t.Fatalf("combat session should not start as over")
	}

	// Attack until enemy dies
	for !session.IsOver && e.HP > 0 {
		_, _, err := session.PlayerAttack()
		if err != nil {
			t.Fatalf("unexpected attack error: %v", err)
		}
	}

	if !session.Won {
		t.Errorf("expected player to win the battle")
	}
	if !session.IsOver {
		t.Errorf("expected combat to be marked as over")
	}

	// Attacking after combat over should return error
	_, _, errPost := session.PlayerAttack()
	if errPost == nil {
		t.Errorf("expected error when attacking after combat is over")
	}
}

func TestCombatDefend(t *testing.T) {
	p := character.NewDefaultPlayer("Defender")
	p.HP = 100
	e := &Enemy{
		ID:         "strong_dummy",
		Name:       "Strong Dummy",
		HP:         100,
		MaxHP:      100,
		MinDamage:  10,
		MaxDamage:  10,
		Initiative: 5,
		Defense:    0,
	}

	session := NewCombatSession(p, e)

	// Player defends
	err := session.PlayerDefend()
	if err != nil {
		t.Fatalf("unexpected defend error: %v", err)
	}

	// Damage should be halved: 10 / 2 = 5 damage
	expectedHP := 100 - 5
	if p.HP != expectedHP {
		t.Errorf("expected HP %d after defended hit, got %d", expectedHP, p.HP)
	}
}

func TestCombatHeal(t *testing.T) {
	p := character.NewDefaultPlayer("Healer")
	p.HP = 50
	p.MaxHP = 100
	e := &Enemy{
		ID:         "weak_dummy",
		Name:       "Weak Dummy",
		HP:         50,
		MaxHP:      50,
		MinDamage:  2,
		MaxDamage:  2,
		Initiative: 5,
		Defense:    0,
	}

	session := NewCombatSession(p, e)

	err := session.PlayerHeal(20)
	if err != nil {
		t.Fatalf("unexpected heal error: %v", err)
	}

	// 50 + 20 - 2 = 68 HP
	expectedHP := 50 + 20 - 2
	if p.HP != expectedHP {
		t.Errorf("expected HP %d, got %d", expectedHP, p.HP)
	}
}

func TestWeaponDurabilityAndAffixes(t *testing.T) {
	p := character.NewDefaultPlayer("AffixTester")
	p.EquippedWeapon = character.Weapon{
		ID:           "test_daggers",
		Name:         "Belati Beracun",
		WeaponType:   "Daggers",
		BaseDamage:   [2]int{10, 10},
		CritRate:     0.0,
		Initiative:   15,
		Durability:   2,
		MaxDura:      2,
		SpecialAffix: "Bisa Beracun (Bleed)",
	}

	e := &Enemy{
		ID:         "target_dummy",
		Name:       "Target Dummy",
		HP:         100,
		MaxHP:      100,
		MinDamage:  1,
		MaxDamage:  1,
		Initiative: 5,
		Defense:    0,
	}

	session := NewCombatSession(p, e)

	// First hit: durability decreases from 2 to 1, bleed dealt
	_, _, err := session.PlayerAttack()
	if err != nil {
		t.Fatalf("unexpected attack error: %v", err)
	}
	if p.EquippedWeapon.Durability != 1 {
		t.Errorf("expected durability 1, got %d", p.EquippedWeapon.Durability)
	}

	// Second hit: durability decreases from 1 to 0
	_, _, _ = session.PlayerAttack()
	if p.EquippedWeapon.Durability != 0 {
		t.Errorf("expected durability 0, got %d", p.EquippedWeapon.Durability)
	}

	// Third hit: broken weapon, durability stays 0, damage halved
	preHP := e.HP
	dmg, _, _ := session.PlayerAttack()
	if p.EquippedWeapon.Durability != 0 {
		t.Errorf("expected durability to stay 0, got %d", p.EquippedWeapon.Durability)
	}
	if dmg > 6 {
		t.Errorf("expected halved damage from broken weapon, got %d", dmg)
	}
	_ = preHP
}

