package character

import (
	"fmt"

	"github.com/Fahmi-mi/terminal-odyssey/data"
)

// CharacterStats holds the four core attributes
type CharacterStats struct {
	Might     int // Physical damage & inventory slots
	Agility   int // Dodge, crit chance, initiative
	Resolve   int // Sanity resistance, max HP
	Ingenuity int // Dialogue options, barter bonus, alchemy
}

// Weapon represents an equipped or crafted weapon
type Weapon struct {
	ID           string
	Name         string
	WeaponType   string // Sword, Blunt, Daggers, Polearm
	BaseDamage   [2]int // min, max
	CritRate     float64
	Initiative   int
	Durability   int
	MaxDura      int
	SpecialAffix string
}

// Companion represents a recruited party member
type Companion struct {
	Name       string
	Role       string // Vanguard, Rogue, Scholar, Acolyte
	CutPercent int
	HP         int
	MaxHP      int
	IsAlive    bool
}

// Player represents the protagonist stats and status
type Player struct {
	Name           string
	Stats          CharacterStats
	HP             int
	MaxHP          int
	Sanity         int
	MaxSanity      int
	TorchMeter     int // 0 to 100%
	EquippedWeapon Weapon
	Party          []Companion
	Backpack       []string
	MaxBackpack    int
}

// NewPlayerFromScenario creates a player from scenario data
func NewPlayerFromScenario(p data.ScenarioPlayer) *Player {
	return &Player{
		Name: p.Name,
		Stats: CharacterStats{
			Might:     p.Stats.Might,
			Agility:   p.Stats.Agility,
			Resolve:   p.Stats.Resolve,
			Ingenuity: p.Stats.Ingenuity,
		},
		HP:          p.HP,
		MaxHP:       p.MaxHP,
		Sanity:      p.Sanity,
		MaxSanity:   p.MaxSanity,
		TorchMeter:  p.TorchMeter,
		MaxBackpack: p.MaxBackpack,
		Backpack:    make([]string, 0),
		Party:       make([]Companion, 0),
		EquippedWeapon: Weapon{
			ID:           p.EquippedWeapon.ID,
			Name:         p.EquippedWeapon.Name,
			WeaponType:   p.EquippedWeapon.WeaponType,
			BaseDamage:   p.EquippedWeapon.BaseDamage,
			CritRate:     p.EquippedWeapon.CritRate,
			Initiative:   p.EquippedWeapon.Initiative,
			Durability:   p.EquippedWeapon.Durability,
			MaxDura:      p.EquippedWeapon.MaxDura,
			SpecialAffix: p.EquippedWeapon.SpecialAffix,
		},
	}
}

// NewDefaultPlayer creates a starter protagonist using default scenario
func NewDefaultPlayer(name string) *Player {
	sc, err := data.LoadDefaultScenario()
	if err == nil {
		p := NewPlayerFromScenario(sc.Player)
		if name != "" {
			p.Name = name
		}
		return p
	}

	// Fallback in case of error
	return &Player{
		Name: name,
		Stats: CharacterStats{
			Might:     10,
			Agility:   10,
			Resolve:   10,
			Ingenuity: 10,
		},
		HP:          100,
		MaxHP:       100,
		Sanity:      100,
		MaxSanity:   100,
		TorchMeter:  100,
		MaxBackpack: 12,
		Backpack:    make([]string, 0),
		Party:       make([]Companion, 0),
		EquippedWeapon: Weapon{
			ID:           "rusty_sword",
			Name:         "Pedang Besi Tua",
			WeaponType:   "Sword",
			BaseDamage:   [2]int{5, 9},
			CritRate:     0.05,
			Initiative:   10,
			Durability:   50,
			MaxDura:      50,
			SpecialAffix: "Standar",
		},
	}
}

// EquipWeapon equips a new weapon on the player
func (p *Player) EquipWeapon(w Weapon) {
	p.EquippedWeapon = w
}

// RepairWeapon restores equipped weapon durability to its maximum
func (p *Player) RepairWeapon() int {
	missing := p.EquippedWeapon.MaxDura - p.EquippedWeapon.Durability
	if missing < 0 {
		missing = 0
	}
	p.EquippedWeapon.Durability = p.EquippedWeapon.MaxDura
	return missing
}

// RecalculateDerivedStats updates MaxHP and MaxBackpack based on Resolve and Might
func (p *Player) RecalculateDerivedStats() {
	resolveDiff := p.Stats.Resolve - 10
	if resolveDiff < 0 {
		resolveDiff = 0
	}
	oldMaxHP := p.MaxHP
	p.MaxHP = 100 + (resolveDiff * 5)
	if p.MaxHP > oldMaxHP {
		p.HP += (p.MaxHP - oldMaxHP)
	}
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}

	mightDiff := p.Stats.Might - 10
	if mightDiff < 0 {
		mightDiff = 0
	}
	p.MaxBackpack = 12 + (mightDiff / 2)
}

// UpgradeStat increments the chosen stat by 1 if below stat cap
func (p *Player) UpgradeStat(statName string, cap int) (int, error) {
	switch statName {
	case "might", "Might":
		if p.Stats.Might >= cap {
			return p.Stats.Might, fmt.Errorf("stat Might sudah mencapai batas maksimal fasilitas (Cap: %d)", cap)
		}
		p.Stats.Might++
		p.RecalculateDerivedStats()
		return p.Stats.Might, nil
	case "agility", "Agility":
		if p.Stats.Agility >= cap {
			return p.Stats.Agility, fmt.Errorf("stat Agility sudah mencapai batas maksimal fasilitas (Cap: %d)", cap)
		}
		p.Stats.Agility++
		p.RecalculateDerivedStats()
		return p.Stats.Agility, nil
	case "resolve", "Resolve":
		if p.Stats.Resolve >= cap {
			return p.Stats.Resolve, fmt.Errorf("stat Resolve sudah mencapai batas maksimal fasilitas (Cap: %d)", cap)
		}
		p.Stats.Resolve++
		p.RecalculateDerivedStats()
		return p.Stats.Resolve, nil
	case "ingenuity", "Ingenuity":
		if p.Stats.Ingenuity >= cap {
			return p.Stats.Ingenuity, fmt.Errorf("stat Ingenuity sudah mencapai batas maksimal fasilitas (Cap: %d)", cap)
		}
		p.Stats.Ingenuity++
		p.RecalculateDerivedStats()
		return p.Stats.Ingenuity, nil
	default:
		return 0, fmt.Errorf("nama atribut %s tidak valid", statName)
	}
}

