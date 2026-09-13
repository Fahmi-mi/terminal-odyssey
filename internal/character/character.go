package character

import (
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
