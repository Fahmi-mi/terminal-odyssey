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
	ID         string
	Name       string
	Role       string // Vanguard, Rogue, Scholar, Acolyte
	PerkDesc   string
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
	OwnedWeapons   []Weapon
	Party          []Companion
	Potions        map[string]int
	Backpack       []string
	MaxBackpack    int
}

// NewPlayerFromScenario creates a player from scenario data
func NewPlayerFromScenario(p data.ScenarioPlayer) *Player {
	starter := Weapon{
		ID:           p.EquippedWeapon.ID,
		Name:         p.EquippedWeapon.Name,
		WeaponType:   p.EquippedWeapon.WeaponType,
		BaseDamage:   p.EquippedWeapon.BaseDamage,
		CritRate:     p.EquippedWeapon.CritRate,
		Initiative:   p.EquippedWeapon.Initiative,
		Durability:   p.EquippedWeapon.Durability,
		MaxDura:      p.EquippedWeapon.MaxDura,
		SpecialAffix: p.EquippedWeapon.SpecialAffix,
	}

	return &Player{
		Name: p.Name,
		Stats: CharacterStats{
			Might:     p.Stats.Might,
			Agility:   p.Stats.Agility,
			Resolve:   p.Stats.Resolve,
			Ingenuity: p.Stats.Ingenuity,
		},
		HP:             p.HP,
		MaxHP:          p.MaxHP,
		Sanity:         p.Sanity,
		MaxSanity:      p.MaxSanity,
		TorchMeter:     p.TorchMeter,
		MaxBackpack:    p.MaxBackpack,
		Backpack:       make([]string, 0),
		Party:          make([]Companion, 0),
		Potions:        make(map[string]int),
		EquippedWeapon: starter,
		OwnedWeapons:   []Weapon{starter},
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
	starter := Weapon{
		ID:           "rusty_sword",
		Name:         "Pedang Besi Tua",
		WeaponType:   "Sword",
		BaseDamage:   [2]int{5, 9},
		CritRate:     0.05,
		Initiative:   10,
		Durability:   50,
		MaxDura:      50,
		SpecialAffix: "Standar",
	}

	return &Player{
		Name: name,
		Stats: CharacterStats{
			Might:     10,
			Agility:   10,
			Resolve:   10,
			Ingenuity: 10,
		},
		HP:             100,
		MaxHP:          100,
		Sanity:         100,
		MaxSanity:      100,
		TorchMeter:     100,
		MaxBackpack:    12,
		Backpack:       make([]string, 0),
		Party:          make([]Companion, 0),
		EquippedWeapon: starter,
		OwnedWeapons:   []Weapon{starter},
	}
}

// EquipWeapon equips a new weapon on the player and saves it into owned weapons
func (p *Player) EquipWeapon(w Weapon) {
	p.SyncEquippedToOwned()
	p.EquippedWeapon = w
	for i, ow := range p.OwnedWeapons {
		if ow.ID == w.ID {
			p.OwnedWeapons[i] = w
			return
		}
	}
	p.OwnedWeapons = append(p.OwnedWeapons, w)
}

// SyncEquippedToOwned updates the equipped weapon state inside the owned list
func (p *Player) SyncEquippedToOwned() {
	if p.EquippedWeapon.ID == "" {
		return
	}
	for i, ow := range p.OwnedWeapons {
		if ow.ID == p.EquippedWeapon.ID {
			p.OwnedWeapons[i] = p.EquippedWeapon
			return
		}
	}
	p.OwnedWeapons = append(p.OwnedWeapons, p.EquippedWeapon)
}

// SwitchWeapon equips an already owned weapon by ID
func (p *Player) SwitchWeapon(weaponID string) error {
	p.SyncEquippedToOwned()
	for _, ow := range p.OwnedWeapons {
		if ow.ID == weaponID {
			p.EquippedWeapon = ow
			return nil
		}
	}
	return fmt.Errorf("senjata dengan ID %s belum dimiliki", weaponID)
}

// OwnsWeapon checks if player has weapon in owned list
func (p *Player) OwnsWeapon(weaponID string) bool {
	for _, ow := range p.OwnedWeapons {
		if ow.ID == weaponID {
			return true
		}
	}
	return false
}

// RepairWeapon restores equipped weapon durability to its maximum
func (p *Player) RepairWeapon() int {
	missing := p.EquippedWeapon.MaxDura - p.EquippedWeapon.Durability
	if missing < 0 {
		missing = 0
	}
	p.EquippedWeapon.Durability = p.EquippedWeapon.MaxDura
	p.SyncEquippedToOwned()
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

// AddPotion adds potions to player pouch
func (p *Player) AddPotion(id string, amount int) {
	if p.Potions == nil {
		p.Potions = make(map[string]int)
	}
	p.Potions[id] += amount
}

// UsePotion decrements potion count if available
func (p *Player) UsePotion(id string) bool {
	if p.Potions == nil || p.Potions[id] <= 0 {
		return false
	}
	p.Potions[id]--
	if p.Potions[id] <= 0 {
		delete(p.Potions, id)
	}
	return true
}

// GetPotionCount returns count of specific potion
func (p *Player) GetPotionCount(id string) int {
	if p.Potions == nil {
		return 0
	}
	return p.Potions[id]
}

// TotalPotions returns the sum of all potions in pouch
func (p *Player) TotalPotions() int {
	total := 0
	for _, cnt := range p.Potions {
		total += cnt
	}
	return total
}

// HasCompanionRole checks if party has a living companion with the given role
func (p *Player) HasCompanionRole(role string) bool {
	for _, c := range p.Party {
		if c.Role == role && c.IsAlive {
			return true
		}
	}
	return false
}

// GetCompanion returns the living companion with the given role
func (p *Player) GetCompanion(role string) *Companion {
	for i := range p.Party {
		if p.Party[i].Role == role && p.Party[i].IsAlive {
			return &p.Party[i]
		}
	}
	return nil
}

// AddCompanion adds a companion to the party
func (p *Player) AddCompanion(c Companion) {
	c.IsAlive = true
	if c.MaxHP > 0 && c.HP <= 0 {
		c.HP = c.MaxHP
	}
	p.Party = append(p.Party, c)
}

// RemoveCompanion removes a companion by ID from the party
func (p *Player) RemoveCompanion(id string) bool {
	for i, c := range p.Party {
		if c.ID == id {
			p.Party = append(p.Party[:i], p.Party[i+1:]...)
			return true
		}
	}
	return false
}


