package tavern

import (
	"fmt"
	"math/rand"

	"github.com/Fahmi-mi/terminal-odyssey/data"
)

// Mercenary roles
const (
	RoleRogue    = "Rogue"
	RoleScholar  = "Scholar"
	RoleVanguard = "Vanguard"
	RoleAcolyte  = "Acolyte"
)

// Mercenary represents a recruitable companion in the tavern
type Mercenary struct {
	Def     data.CompanionDef
	IsHired bool
}

// TavernManager coordinates companions and ambient rumors
type TavernManager struct {
	Mercenaries []*Mercenary
	Rumors      []string
}

// NewTavernManager initializes the tavern from companion definitions
func NewTavernManager() (*TavernManager, error) {
	defs, err := data.LoadCompanionDefs()
	if err != nil {
		return nil, err
	}

	mercs := make([]*Mercenary, len(defs))
	for i, d := range defs {
		mercs[i] = &Mercenary{
			Def:     d,
			IsHired: false,
		}
	}

	rumors := []string{
		"Musafir utara berbisik bahwa Metropolis Dunhallow siap membayar mahal untuk pasokan rempah wangi",
		"Seorang penambang tua bersumpah melihat pintu suaka rahasia tersembunyi di balik altar katakombe",
		"Kafilah dagang mengabarkan bahwa musim dingin memperpanjang perjalanan gerobak hingga dua kali lipat",
		"Penduduk perbatasan memperingatkan keberadaan laba-laba raksasa yang menyarangkan racun di kedalaman lorong",
		"Beredar kabar bahwa sarjana kuno mampu menerjemahkan kutukan altar menjadi limpahan relik emas murni",
	}

	return &TavernManager{
		Mercenaries: mercs,
		Rumors:      rumors,
	}, nil
}

// MaxPartySize returns maximum companions allowed by Tavern level
func (tm *TavernManager) MaxPartySize(tavernLevel int) int {
	if tavernLevel < 1 {
		return 0
	}
	if tavernLevel > 3 {
		return 3
	}
	return tavernLevel
}

// GetMercenary returns mercenary by ID
func (tm *TavernManager) GetMercenary(id string) *Mercenary {
	for _, m := range tm.Mercenaries {
		if m.Def.ID == id {
			return m
		}
	}
	return nil
}

// CanHire validates whether a companion can be recruited
func (tm *TavernManager) CanHire(id string, tavernLevel int, treasury int, currentPartyCount int) error {
	if tavernLevel < 1 {
		return fmt.Errorf("kedai Minum belum didirikan di desa")
	}

	merc := tm.GetMercenary(id)
	if merc == nil {
		return fmt.Errorf("rekan petualang tidak ditemukan")
	}

	if merc.IsHired {
		return fmt.Errorf("%s sudah bergabung dalam rombongan", merc.Def.Name)
	}

	maxParty := tm.MaxPartySize(tavernLevel)
	if currentPartyCount >= maxParty {
		return fmt.Errorf("rombongan sudah penuh (%d/%d pendamping untuk Kedai Minum Level %d)", currentPartyCount, maxParty, tavernLevel)
	}

	if treasury < merc.Def.HireCost {
		return fmt.Errorf("kas emas tidak mencukupi biaya sewa (butuh %d Gold, ada %d)", merc.Def.HireCost, treasury)
	}

	return nil
}

// Hire recruits a companion into the party
func (tm *TavernManager) Hire(id string, tavernLevel int, treasury int, currentPartyCount int) (*Mercenary, error) {
	if err := tm.CanHire(id, tavernLevel, treasury, currentPartyCount); err != nil {
		return nil, err
	}

	merc := tm.GetMercenary(id)
	merc.IsHired = true
	return merc, nil
}

// Dismiss releases a companion from the active party
func (tm *TavernManager) Dismiss(id string) error {
	merc := tm.GetMercenary(id)
	if merc == nil {
		return fmt.Errorf("pendamping tidak ditemukan")
	}
	if !merc.IsHired {
		return fmt.Errorf("%s sedang tidak berada dalam rombongan", merc.Def.Name)
	}
	merc.IsHired = false
	return nil
}

// CanRest checks whether player can afford tavern meal and rest
func (tm *TavernManager) CanRest(treasury, rations int) error {
	if treasury < 10 {
		return fmt.Errorf("butuh 10 Gold untuk memesan hidangan kedai")
	}
	if rations < 1 {
		return fmt.Errorf("butuh 1 Ransum untuk santapan kedai")
	}
	return nil
}

// GetRandomRumor returns an atmospheric rumor
func (tm *TavernManager) GetRandomRumor() string {
	if len(tm.Rumors) == 0 {
		return "Suasana kedai minum tampak tenang dan bersahabat"
	}
	return tm.Rumors[rand.Intn(len(tm.Rumors))]
}
