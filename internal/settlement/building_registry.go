package settlement

import (
	"fmt"

	"github.com/Fahmi-mi/terminal-odyssey/data"
)

var (
	cachedBuildingDefs map[string]data.BuildingDef
	buildingOrder      []string
)

func init() {
	defs, err := data.LoadBuildingDefs()
	if err == nil {
		cachedBuildingDefs = make(map[string]data.BuildingDef)
		for _, d := range defs {
			cachedBuildingDefs[d.Name] = d
			buildingOrder = append(buildingOrder, d.Name)
		}
	}
}

// GetBuildingInfo returns dynamic cost and status for a given building loaded from data
func (s *Settlement) GetBuildingInfo(name string) BuildingInfo {
	lvl := s.Buildings[name]
	nextLvl := lvl + 1

	def, exists := cachedBuildingDefs[name]
	if !exists {
		return BuildingInfo{Name: name, Level: lvl}
	}

	return BuildingInfo{
		Name:        def.Name,
		Description: def.Description,
		Level:       lvl,
		MaxLevel:    def.MaxLevel,
		WoodCost:    def.BaseWoodCost * nextLvl,
		StoneCost:   def.BaseStoneCost * nextLvl,
		GoldCost:    def.BaseGoldCost * nextLvl,
	}
}

// AllBuildingsList returns list of all building names in the configured order
func AllBuildingsList() []string {
	if len(buildingOrder) > 0 {
		return buildingOrder
	}
	return []string{
		BuildingTownHall,
		BuildingStorehouse,
		BuildingBlacksmith,
		BuildingTrainingGround,
		BuildingApothecary,
		BuildingTavern,
		BuildingCaravanPost,
	}
}

// UpgradeBuilding attempts to upgrade a building by spending required resources
func (s *Settlement) UpgradeBuilding(name string) error {
	info := s.GetBuildingInfo(name)
	if info.Level >= info.MaxLevel {
		return fmt.Errorf("%s sudah mencapai level maksimum (%d)", name, info.MaxLevel)
	}

	if s.Lumber < info.WoodCost {
		return fmt.Errorf("kayu tidak mencukupi (butuh %d, ada %d)", info.WoodCost, s.Lumber)
	}
	if s.Stone < info.StoneCost {
		return fmt.Errorf("batu tidak mencukupi (butuh %d, ada %d)", info.StoneCost, s.Stone)
	}
	if s.Treasury < info.GoldCost {
		return fmt.Errorf("emas tidak mencukupi (butuh %d, ada %d)", info.GoldCost, s.Treasury)
	}

	// Deduct resources
	s.Lumber -= info.WoodCost
	s.Stone -= info.StoneCost
	s.Treasury -= info.GoldCost

	// Upgrade
	s.Buildings[name]++
	s.RecalculateDefense()
	return nil
}
