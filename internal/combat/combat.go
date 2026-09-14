package combat

import (
	"fmt"
	"math/rand"

	"github.com/Fahmi-mi/terminal-odyssey/data"
	"github.com/Fahmi-mi/terminal-odyssey/internal/character"
)

// Enemy represents an instantiated dungeon monster
type Enemy struct {
	ID          string
	Name        string
	Description string
	HP          int
	MaxHP       int
	MinDamage   int
	MaxDamage   int
	Initiative  int
	Defense     int
	GoldReward  int
	ExpReward   int
}

// NewEnemyFromDef creates an Enemy from data definition
func NewEnemyFromDef(def data.EnemyDef) *Enemy {
	return &Enemy{
		ID:          def.ID,
		Name:        def.Name,
		Description: def.Description,
		HP:          def.MaxHP,
		MaxHP:       def.MaxHP,
		MinDamage:   def.MinDamage,
		MaxDamage:   def.MaxDamage,
		Initiative:  def.Initiative,
		Defense:     def.Defense,
		GoldReward:  def.GoldReward,
		ExpReward:   def.ExpReward,
	}
}

// NewEnemyByID loads enemy definition and instantiates Enemy
func NewEnemyByID(id string) (*Enemy, error) {
	defs, err := data.LoadEnemyDefs()
	if err != nil {
		return nil, err
	}
	for _, d := range defs {
		if d.ID == id {
			return NewEnemyFromDef(d), nil
		}
	}
	return nil, fmt.Errorf("musuh dengan ID %s tidak ditemukan", id)
}

// CombatSession manages a turn-based battle between player and enemy
type CombatSession struct {
	Player          *character.Player
	Enemy           *Enemy
	PlayerDefending bool
	TurnCount       int
	IsOver          bool
	Won             bool
	Fled            bool
	Logs            []string
}

// NewCombatSession starts a new battle and logs initial engagement
func NewCombatSession(p *character.Player, enemy *Enemy) *CombatSession {
	s := &CombatSession{
		Player:    p,
		Enemy:     enemy,
		TurnCount: 1,
		Logs:      make([]string, 0),
	}

	playerInit := p.Stats.Agility + p.EquippedWeapon.Initiative
	s.addLog(fmt.Sprintf("[!] Bertemu dengan %s (HP %d/%d)", enemy.Name, enemy.HP, enemy.MaxHP))

	if playerInit >= enemy.Initiative {
		s.addLog("[i] Kecepatan inisiatif unggul! Giliran Anda melangkah terlebih dahulu")
	} else {
		s.addLog(fmt.Sprintf("[!] %s bergerak cepat bersiap menerkam", enemy.Name))
	}

	return s
}

// PlayerAttack executes a physical attack against the enemy
func (s *CombatSession) PlayerAttack() (int, bool, error) {
	if s.IsOver {
		return 0, false, fmt.Errorf("pertempuran sudah berakhir")
	}

	weaponMin := s.Player.EquippedWeapon.BaseDamage[0]
	weaponMax := s.Player.EquippedWeapon.BaseDamage[1]
	if weaponMax < weaponMin {
		weaponMax = weaponMin
	}

	rawDmg := rand.Intn(weaponMax-weaponMin+1) + weaponMin
	mightBonus := (s.Player.Stats.Might - 10) / 2
	if mightBonus < 0 {
		mightBonus = 0
	}

	critChance := s.Player.EquippedWeapon.CritRate + float64(s.Player.Stats.Agility)*0.005
	isCrit := rand.Float64() < critChance

	total := rawDmg + mightBonus
	if isCrit {
		total = int(float64(total) * 1.5)
	}

	netDmg := total - s.Enemy.Defense
	if netDmg < 1 {
		netDmg = 1
	}

	s.Enemy.HP -= netDmg
	if s.Enemy.HP < 0 {
		s.Enemy.HP = 0
	}

	if isCrit {
		s.addLog(fmt.Sprintf("[+] SERANGAN KRITIKAL! Tebasan Anda menembus pertahanan %s (-%d HP)", s.Enemy.Name, netDmg))
	} else {
		s.addLog(fmt.Sprintf("[+] Serangan %s mengenai %s (-%d HP)", s.Player.EquippedWeapon.Name, s.Enemy.Name, netDmg))
	}

	// Check if enemy defeated
	if s.Enemy.HP <= 0 {
		s.IsOver = true
		s.Won = true
		s.addLog(fmt.Sprintf("[+] %s roboh tak berdaya! Anda memenangkan pertempuran", s.Enemy.Name))
		s.addLog(fmt.Sprintf("[*] Memperoleh jarahan +%d Emas", s.Enemy.GoldReward))
		return netDmg, isCrit, nil
	}

	// Enemy counter-attacks
	s.enemyCounterAttack()
	s.TurnCount++
	return netDmg, isCrit, nil
}

// PlayerDefend enters defensive stance, reducing next enemy damage by half
func (s *CombatSession) PlayerDefend() error {
	if s.IsOver {
		return fmt.Errorf("pertempuran sudah berakhir")
	}

	s.PlayerDefending = true
	s.addLog("[i] Anda memasang kuda-kuda bertahan (+50% reduksi kerusakan musuh)")

	s.enemyCounterAttack()
	s.TurnCount++
	return nil
}

// PlayerHeal consumes a ration from expedition to recover player HP
func (s *CombatSession) PlayerHeal(healAmount int) error {
	if s.IsOver {
		return fmt.Errorf("pertempuran sudah berakhir")
	}

	oldHP := s.Player.HP
	s.Player.HP += healAmount
	if s.Player.HP > s.Player.MaxHP {
		s.Player.HP = s.Player.MaxHP
	}
	gained := s.Player.HP - oldHP

	s.addLog(fmt.Sprintf("[+] Memakan ransum darurat (+%d HP)", gained))

	s.enemyCounterAttack()
	s.TurnCount++
	return nil
}

// PlayerFlee attempts to retreat from the battlefield
func (s *CombatSession) PlayerFlee() (bool, error) {
	if s.IsOver {
		return false, fmt.Errorf("pertempuran sudah berakhir")
	}

	fleeChance := 0.40 + float64(s.Player.Stats.Agility)*0.02
	if fleeChance > 0.85 {
		fleeChance = 0.85
	}

	if rand.Float64() < fleeChance {
		s.IsOver = true
		s.Fled = true
		s.addLog("[!] Anda berhasil melompat mundur dan meloloskan diri dari kepungan")
		return true, nil
	}

	s.addLog("[!] Percobaan kabur gagal! Musuh memotong jalur mundur Anda")
	s.enemyCounterAttack()
	s.TurnCount++
	return false, nil
}

// enemyCounterAttack executes the enemy turn
func (s *CombatSession) enemyCounterAttack() {
	if s.Enemy.HP <= 0 {
		return
	}

	eMin := s.Enemy.MinDamage
	eMax := s.Enemy.MaxDamage
	if eMax < eMin {
		eMax = eMin
	}

	dmg := rand.Intn(eMax-eMin+1) + eMin
	if s.PlayerDefending {
		dmg = dmg / 2
		if dmg < 1 {
			dmg = 1
		}
		s.addLog(fmt.Sprintf("[-] Pertahanan Anda meredam serangan %s (-%d HP)", s.Enemy.Name, dmg))
		s.PlayerDefending = false
	} else {
		s.addLog(fmt.Sprintf("[-] %s menyerang Anda (-%d HP)", s.Enemy.Name, dmg))
	}

	s.Player.HP -= dmg
	if s.Player.HP <= 0 {
		s.Player.HP = 0
		s.IsOver = true
		s.Won = false
		s.addLog(fmt.Sprintf("[!] Tubuh Anda tumbang tak sadarkan diri akibat serangan %s", s.Enemy.Name))
	}
}

// addLog keeps the combat log to the most recent entries
func (s *CombatSession) addLog(entry string) {
	s.Logs = append(s.Logs, entry)
	if len(s.Logs) > 8 {
		s.Logs = s.Logs[len(s.Logs)-8:]
	}
}
