package alchemy

import (
	"fmt"
	"math/rand"

	"github.com/Fahmi-mi/terminal-odyssey/data"
)

// Effect types for alchemy items
const (
	EffectHealHP        = "heal_hp"
	EffectRefillTorch   = "refill_torch"
	EffectRestoreSanity = "restore_sanity"
	EffectCureDebuff    = "cure_debuff"
	EffectBuffAtk       = "buff_atk"
)

// BrewResult represents the outcome of an alchemy brewing session
type BrewResult struct {
	Recipe      data.AlchemyRecipeDef
	YieldAmount int
	IsCritical  bool
	Log         string
}

// AlchemyManager coordinates recipe lookup and brewing mechanics
type AlchemyManager struct {
	Recipes []data.AlchemyRecipeDef
}

// NewAlchemyManager initializes alchemy recipes from data
func NewAlchemyManager() (*AlchemyManager, error) {
	defs, err := data.LoadAlchemyRecipeDefs()
	if err != nil {
		return nil, err
	}
	return &AlchemyManager{
		Recipes: defs,
	}, nil
}

// GetRecipe finds a recipe by ID
func (m *AlchemyManager) GetRecipe(id string) *data.AlchemyRecipeDef {
	for i := range m.Recipes {
		if m.Recipes[i].ID == id {
			return &m.Recipes[i]
		}
	}
	return nil
}

// CanBrew validates whether the village has the required building level and resources
func (m *AlchemyManager) CanBrew(recipeID string, apothecaryLevel int, treasury int, commodities map[string]int, lumber, stone, rations int) error {
	recipe := m.GetRecipe(recipeID)
	if recipe == nil {
		return fmt.Errorf("resep alkimia tidak ditemukan")
	}

	if apothecaryLevel < recipe.ApothecaryLevel {
		return fmt.Errorf("butuh Laboratorium Alkimia Level %d", recipe.ApothecaryLevel)
	}

	if treasury < recipe.GoldCost {
		return fmt.Errorf("kas emas tidak mencukupi (butuh %d Gold, ada %d)", recipe.GoldCost, treasury)
	}

	for ingredientID, requiredAmt := range recipe.Ingredients {
		available := 0
		switch ingredientID {
		case "lumber":
			available = lumber
		case "stone":
			available = stone
		case "rations":
			available = rations
		default:
			if commodities != nil {
				available = commodities[ingredientID]
			}
		}

		if available < requiredAmt {
			return fmt.Errorf("bahan %s tidak mencukupi (butuh %d, ada %d)", ingredientID, requiredAmt, available)
		}
	}

	return nil
}

// Brew executes the brewing process with potential Ingenuity critical bonus
func (m *AlchemyManager) Brew(recipeID string, apothecaryLevel int, treasury int, commodities map[string]int, lumber, stone, rations int, ingenuity int) (*BrewResult, error) {
	if err := m.CanBrew(recipeID, apothecaryLevel, treasury, commodities, lumber, stone, rations); err != nil {
		return nil, err
	}

	recipe := m.GetRecipe(recipeID)
	yield := 1
	isCrit := false

	// Ingenuity bonus: chance for extra yield (+1 potion)
	diff := ingenuity - 10
	if diff > 0 {
		critChance := float64(diff) * 0.04 // 4% per point above 10
		if critChance > 0.35 {
			critChance = 0.35
		}
		if rand.Float64() < critChance {
			yield = 2
			isCrit = true
		}
	}

	logMsg := fmt.Sprintf("Berhasil meracik 1 %s", recipe.Name)
	if isCrit {
		logMsg = fmt.Sprintf("Racikan Alkimia Sempurna! Kecerdikan memicu ekstraksi murni (+2 %s)", recipe.Name)
	}

	return &BrewResult{
		Recipe:      *recipe,
		YieldAmount: yield,
		IsCritical:  isCrit,
		Log:         logMsg,
	}, nil
}
